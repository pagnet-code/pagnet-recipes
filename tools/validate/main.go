// Command validate checks every manifest.yaml in the recipes repository
// against the pagnet.dev/v1 manifest format (plan D13):
//
//   - apiVersion is pagnet.dev/v1
//   - kind is one of Recipe, Service, Agent, AgentTemplate
//   - metadata.name is present and kebab-case
//   - capability ids are dot-separated (e.g. echo.say)
//   - requested permissions are within the 8 membership permissions
//   - subscription modes are deliver or wake
//   - no fields outside the format ("no more fields")
//   - kind-appropriate fields (a Recipe has contains, a template has no
//     capabilities, ...)
//   - a Recipe's contains entries resolve to sub-manifest files
//
// Run it from tools/validate (go run .) or pass a repo root as the first
// argument. It is the local and CI gate for manifest correctness.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	kinds = map[string]bool{
		"Recipe": true, "Service": true, "Agent": true, "AgentTemplate": true,
	}
	permissions = map[string]bool{
		"discover": true, "communicate": true, "invoke": true,
		"event_publish": true, "event_subscribe": true,
		"task_read": true, "task_write": true, "operate": true,
	}
	kebab = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	capID = regexp.MustCompile(`^[a-z0-9]+(\.[a-z0-9]+)+$`)
)

// manifest is the complete pagnet.dev/v1 field set; the strict decoder
// rejects any field outside it.
type manifest struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
	} `yaml:"metadata"`
	Capabilities []struct {
		ID           string         `yaml:"id"`
		Version      int            `yaml:"version"`
		Description  string         `yaml:"description"`
		InputSchema  map[string]any `yaml:"inputSchema"`
		OutputSchema map[string]any `yaml:"outputSchema"`
		Tags         []string       `yaml:"tags"`
	} `yaml:"capabilities"`
	Subscriptions []struct {
		Event string `yaml:"event"`
		Mode  string `yaml:"mode"`
	} `yaml:"subscriptions"`
	Permissions struct {
		Requested []string `yaml:"requested"`
	} `yaml:"permissions"`
	Runtime     string   `yaml:"runtime"`
	Mission     string   `yaml:"mission"`
	Instruction string   `yaml:"instruction"`
	Contains    []string `yaml:"contains"`
}

func main() {
	root := repoRoot()
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		if !d.IsDir() && d.Name() == "manifest.yaml" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fatalf("walk %s: %v", root, err)
	}
	if len(files) == 0 {
		fatalf("no manifest.yaml files found under %s", root)
	}

	failed := false
	for _, f := range files {
		if err := check(f, root); err != nil {
			fmt.Fprintf(os.Stderr, "FAIL  %s: %v\n", rel(root, f), err)
			failed = true
			continue
		}
		fmt.Printf("ok    %s\n", rel(root, f))
	}
	if failed {
		os.Exit(1)
	}
	fmt.Printf("%d manifests valid\n", len(files))
}

func check(path, root string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	dec := yaml.NewDecoder(strings.NewReader(string(data)))
	dec.KnownFields(true)
	var m manifest
	if err := dec.Decode(&m); err != nil {
		return fmt.Errorf("yaml: %w", err)
	}
	if m.APIVersion != "pagnet.dev/v1" {
		return fmt.Errorf("apiVersion must be pagnet.dev/v1, got %q", m.APIVersion)
	}
	if !kinds[m.Kind] {
		return fmt.Errorf("kind must be Recipe, Service, Agent or AgentTemplate, got %q", m.Kind)
	}
	if !kebab.MatchString(m.Metadata.Name) {
		return fmt.Errorf("metadata.name must be kebab-case, got %q", m.Metadata.Name)
	}
	for _, c := range m.Capabilities {
		if !capID.MatchString(c.ID) {
			return fmt.Errorf("capability id %q must be dot-separated (e.g. echo.say)", c.ID)
		}
	}
	for _, s := range m.Subscriptions {
		if s.Mode != "deliver" && s.Mode != "wake" {
			return fmt.Errorf("subscription %q: mode must be deliver or wake, got %q", s.Event, s.Mode)
		}
	}
	for _, p := range m.Permissions.Requested {
		if !permissions[p] {
			return fmt.Errorf("permission %q is not one of the 8 membership permissions", p)
		}
	}
	switch m.Kind {
	case "Recipe":
		if len(m.Capabilities) > 0 || len(m.Subscriptions) > 0 || len(m.Permissions.Requested) > 0 {
			return fmt.Errorf("Recipe must not declare capabilities, subscriptions or permissions")
		}
		if len(m.Contains) == 0 {
			return fmt.Errorf("Recipe must list sub-manifests in contains")
		}
		for _, c := range m.Contains {
			sub := filepath.Join(filepath.Dir(path), c, "manifest.yaml")
			if _, err := os.Stat(sub); err != nil {
				return fmt.Errorf("contains %q: no sub-manifest at %s", c, rel(root, sub))
			}
		}
	case "AgentTemplate":
		if len(m.Capabilities) > 0 || len(m.Permissions.Requested) > 0 {
			return fmt.Errorf("AgentTemplate must not declare capabilities or permissions")
		}
	case "Service", "Agent":
		if m.Runtime != "" || m.Mission != "" || m.Instruction != "" || len(m.Contains) > 0 {
			return fmt.Errorf("%s must not declare template or recipe fields", m.Kind)
		}
	}
	return nil
}

func repoRoot() string {
	if len(os.Args) > 1 {
		abs, err := filepath.Abs(os.Args[1])
		if err != nil {
			fatalf("root: %v", err)
		}
		return abs
	}
	cwd, err := os.Getwd()
	if err != nil {
		fatalf("getwd: %v", err)
	}
	// The module lives at <repo>/tools/validate.
	return filepath.Clean(filepath.Join(cwd, "..", ".."))
}

func rel(root, p string) string {
	r, err := filepath.Rel(root, p)
	if err != nil {
		return p
	}
	return r
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "validate: "+format+"\n", args...)
	os.Exit(1)
}

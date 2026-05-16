package values

import (
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestFlattenNestedYAML(t *testing.T) {
	root := parseYAML(t, `
service:
  type: LoadBalancer
ports:
  web:
    nodePort: null
`)

	entries := Flatten(root)
	got := paths(entries)
	want := []string{"ports.web.nodePort", "service.type"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("paths = %#v, want %#v", got, want)
	}
}

func TestSearchMatchesKeywordsCaseInsensitiveOR(t *testing.T) {
	entries := []ValueEntry{
		{Path: "service.type", Value: "LoadBalancer"},
		{Path: "ports.web.nodePort", Value: nil},
		{Path: "service.externalTrafficPolicy", Value: "Cluster"},
	}

	matches := Search(entries, []string{"NODEPORT", "externalTrafficPolicy"})
	got := paths(matches)
	want := []string{"ports.web.nodePort", "service.externalTrafficPolicy"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("paths = %#v, want %#v", got, want)
	}
}

func TestBuildSkeletonYAML(t *testing.T) {
	out, err := BuildSkeletonYAML([]ValueEntry{
		{Path: "service.type", Value: "LoadBalancer"},
		{Path: "ports.web.nodePort", Value: nil},
		{Path: "ports.websecure.nodePort", Value: nil},
	})
	if err != nil {
		t.Fatal(err)
	}

	text := string(out)
	for _, want := range []string{
		"service:",
		"  type: LoadBalancer",
		"ports:",
		"    nodePort: null",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("output missing %q:\n%s", want, text)
		}
	}
}

func TestFlattenTreatsListAsLeaf(t *testing.T) {
	root := parseYAML(t, `
additionalArguments:
  - "--api.insecure=true"
`)

	entries := Flatten(root)
	if len(entries) != 1 {
		t.Fatalf("len(entries) = %d, want 1", len(entries))
	}
	if entries[0].Path != "additionalArguments" {
		t.Fatalf("path = %q", entries[0].Path)
	}
	list, ok := entries[0].Value.([]any)
	if !ok || len(list) != 1 || list[0] != "--api.insecure=true" {
		t.Fatalf("value = %#v", entries[0].Value)
	}
}

func TestNullValuePreserved(t *testing.T) {
	root := parseYAML(t, `
ports:
  web:
    nodePort: null
`)

	entries := Flatten(root)
	if len(entries) != 1 {
		t.Fatalf("len(entries) = %d, want 1", len(entries))
	}
	if entries[0].Value != nil {
		t.Fatalf("value = %#v, want nil", entries[0].Value)
	}

	out, err := BuildSkeletonYAML(entries)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "nodePort: null") {
		t.Fatalf("null was not preserved:\n%s", out)
	}
}

func parseYAML(t *testing.T, input string) any {
	t.Helper()
	var root any
	if err := yaml.Unmarshal([]byte(input), &root); err != nil {
		t.Fatal(err)
	}
	return root
}

func paths(entries []ValueEntry) []string {
	out := make([]string, len(entries))
	for i, entry := range entries {
		out[i] = entry.Path
	}
	return out
}

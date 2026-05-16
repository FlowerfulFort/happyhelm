package values

import (
	"fmt"
	"sort"
	"strings"
)

type ValueEntry struct {
	Path  string
	Value any
}

func Flatten(root any) []ValueEntry {
	var entries []ValueEntry
	flattenValue("", normalize(root), &entries)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})
	return entries
}

func flattenValue(prefix string, value any, entries *[]ValueEntry) {
	switch typed := value.(type) {
	case map[string]any:
		if len(typed) == 0 && prefix != "" {
			*entries = append(*entries, ValueEntry{Path: prefix, Value: typed})
			return
		}
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			path := key
			if prefix != "" {
				path = prefix + "." + key
			}
			flattenValue(path, typed[key], entries)
		}
	default:
		if prefix != "" {
			*entries = append(*entries, ValueEntry{Path: prefix, Value: typed})
		}
	}
}

func normalize(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for key, val := range typed {
			out[key] = normalize(val)
		}
		return out
	case map[any]any:
		out := make(map[string]any, len(typed))
		for key, val := range typed {
			out[fmt.Sprint(key)] = normalize(val)
		}
		return out
	case []any:
		out := make([]any, len(typed))
		for i, val := range typed {
			out[i] = normalize(val)
		}
		return out
	default:
		return typed
	}
}

func pathParts(path string) []string {
	raw := strings.Split(path, ".")
	parts := raw[:0]
	for _, part := range raw {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

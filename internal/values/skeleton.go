package values

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

func BuildSkeleton(entries []ValueEntry) (map[string]any, error) {
	root := make(map[string]any)
	for _, entry := range entries {
		parts := pathParts(entry.Path)
		if len(parts) == 0 {
			return nil, fmt.Errorf("invalid empty path")
		}

		current := root
		for i, part := range parts {
			if i == len(parts)-1 {
				current[part] = entry.Value
				continue
			}

			next, ok := current[part]
			if !ok {
				child := make(map[string]any)
				current[part] = child
				current = child
				continue
			}

			child, ok := next.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("path conflict at %q", part)
			}
			current = child
		}
	}
	return root, nil
}

func BuildSkeletonYAML(entries []ValueEntry) ([]byte, error) {
	skeleton, err := BuildSkeleton(entries)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	encoder := yaml.NewEncoder(&out)
	encoder.SetIndent(2)
	if err := encoder.Encode(skeleton); err != nil {
		return nil, fmt.Errorf("marshal selected values: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("close YAML encoder: %w", err)
	}
	return out.Bytes(), nil
}

package values

import "strings"

func Search(entries []ValueEntry, keywords []string) []ValueEntry {
	if len(keywords) == 0 {
		return entries
	}

	normalized := make([]string, 0, len(keywords))
	for _, keyword := range keywords {
		keyword = strings.TrimSpace(strings.ToLower(keyword))
		if keyword != "" {
			normalized = append(normalized, keyword)
		}
	}
	if len(normalized) == 0 {
		return entries
	}

	var matches []ValueEntry
	for _, entry := range entries {
		path := strings.ToLower(entry.Path)
		for _, keyword := range normalized {
			if path == keyword || strings.Contains(path, keyword) {
				matches = append(matches, entry)
				break
			}
		}
	}
	return matches
}

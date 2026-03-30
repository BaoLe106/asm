package util

import (
	"path/filepath"
	"runtime"
	"strings"
)

func ToSlash(rel string) string {
	return filepath.ToSlash(rel)
}

func NormalizeSkillList(raw string) []string {
	parts := strings.Split(raw, ",")
	seen := map[string]struct{}{}
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

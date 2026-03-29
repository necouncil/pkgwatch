package upstream

import (
	"strings"
)

func LatestVersion(sources []string) (string, error) {
	for _, src := range sources {
		src = stripURLPrefix(src)

		if strings.Contains(src, "github.com") {
			owner, repo, ok := parseGitHubURL(src)
			if !ok {
				continue
			}
			return GitHubLatest(owner, repo)
		}

		if strings.Contains(src, "gitlab.com") {
			project, ok := parseGitLabURL(src)
			if !ok {
				continue
			}
			return GitLabLatest(project)
		}
	}
	return "", nil
}

func IsNewer(latest, current string) bool {
	l := normalizeVersion(latest)
	c := normalizeVersion(current)
	return compareVersions(l, c) > 0
}

func stripURLPrefix(src string) string {
	if idx := strings.Index(src, "::"); idx != -1 {
		return src[idx+2:]
	}
	return src
}

func normalizeVersion(v string) string {
	v = strings.TrimPrefix(v, "v")
	v = strings.TrimPrefix(v, "V")
	return v
}

func compareVersions(a, b string) int {
	ap := splitVersion(a)
	bp := splitVersion(b)

	max := len(ap)
	if len(bp) > max {
		max = len(bp)
	}

	for i := 0; i < max; i++ {
		av := 0
		bv := 0
		if i < len(ap) {
			av = parseNum(ap[i])
		}
		if i < len(bp) {
			bv = parseNum(bp[i])
		}
		if av > bv {
			return 1
		}
		if av < bv {
			return -1
		}
	}
	return 0
}

func splitVersion(v string) []string {
	return strings.FieldsFunc(v, func(r rune) bool {
		return r == '.' || r == '-' || r == '_'
	})
}

func parseNum(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			break
		}
	}
	return n
}

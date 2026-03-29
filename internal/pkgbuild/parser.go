package pkgbuild

import (
	"bufio"
	"os"
	"strings"
)

type Package struct {
	Name    string
	Version string
	Source  []string
}

func Parse(path string) (*Package, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pkg := &Package{}
	scanner := bufio.NewScanner(f)
	var inSource bool
	var sourceBuf strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		if inSource {
			sourceBuf.WriteString(" " + trimmed)
			if strings.Contains(trimmed, ")") {
				inSource = false
				pkg.Source = parseArray(sourceBuf.String())
			}
			continue
		}

		if strings.HasPrefix(trimmed, "pkgname=") {
			pkg.Name = extractValue(trimmed, "pkgname=")
		} else if strings.HasPrefix(trimmed, "pkgver=") {
			pkg.Version = extractValue(trimmed, "pkgver=")
		} else if strings.HasPrefix(trimmed, "source=") {
			val := strings.TrimPrefix(trimmed, "source=")
			if strings.HasPrefix(val, "(") && !strings.Contains(val, ")") {
				inSource = true
				sourceBuf.Reset()
				sourceBuf.WriteString(val)
			} else {
				pkg.Source = parseArray(val)
			}
		}
	}

	return pkg, scanner.Err()
}

func extractValue(line, prefix string) string {
	val := strings.TrimPrefix(line, prefix)
	val = strings.Trim(val, `"'`)
	return val
}

func parseArray(s string) []string {
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(strings.TrimSpace(s), ")")
	var result []string
	for _, part := range strings.Fields(s) {
		part = strings.Trim(part, `"'`)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

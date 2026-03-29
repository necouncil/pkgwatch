package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func parseGitLabURL(src string) (project string, ok bool) {
	src = strings.TrimSuffix(src, ".git")
	parts := strings.Split(src, "gitlab.com/")
	if len(parts) < 2 {
		return "", false
	}
	project = strings.Trim(parts[1], "/")
	if project == "" {
		return "", false
	}
	return project, true
}

func GitLabLatest(project string) (string, error) {
	encoded := url.PathEscape(project)
	apiURL := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/releases", encoded)

	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gitlab api: %s", resp.Status)
	}

	var releases []struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", err
	}

	if len(releases) == 0 {
		return "", nil
	}

	return strings.TrimPrefix(releases[0].TagName, "v"), nil
}

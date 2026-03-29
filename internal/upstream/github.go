package upstream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func parseGitHubURL(src string) (owner, repo string, ok bool) {
	src = strings.TrimSuffix(src, ".git")
	parts := strings.Split(src, "github.com/")
	if len(parts) < 2 {
		return "", "", false
	}
	segments := strings.Split(strings.Trim(parts[1], "/"), "/")
	if len(segments) < 2 {
		return "", "", false
	}
	return segments[0], segments[1], true
}

func GitHubLatest(owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api: %s", resp.Status)
	}

	var result struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return strings.TrimPrefix(result.TagName, "v"), nil
}

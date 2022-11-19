package change

import "testing"

func TestParseGithubURL(t *testing.T) {
	urls := []struct {
		url   string
		owner string
		repo  string
	}{
		{"https://github.com/fengleng/mars.git", "go-mars", "mars"},
		{"https://github.com/fengleng/mars", "go-mars", "mars"},
		{"git@github.com:go-mars/mars.git", "go-mars", "mars"},
		{"https://github.com/go-mars/go-mars.dev.git", "go-mars", "go-mars.dev"},
	}
	for _, url := range urls {
		owner, repo := ParseGithubURL(url.url)
		if owner != url.owner {
			t.Fatalf("owner want: %s, got: %s", owner, url.owner)
		}
		if repo != url.repo {
			t.Fatalf("repo want: %s, got: %s", repo, url.repo)
		}
	}
}

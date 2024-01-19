package reporef

import (
	"errors"
	"strings"
)

type RepoRef struct {
	URL  string
	Path string
	Name string
	Tag  string
}

func ParseRepoRef(repoRef string) (*RepoRef, error) {
	parts := strings.Split(repoRef, "@")
	if len(parts) != 2 {
		return nil, errors.New("repoRef must be in the format <url>@<tag>")
	}

	url, tag := parts[0], parts[1]
	parts = strings.Split(url, "/")
	if len(parts) < 3 {
		return nil, errors.New("url must be in the format <host>/<owner>/<repo>")
	}
	name := parts[len(parts)-1]
	url = strings.Join(parts[:3], "/")
	path := strings.Join(parts[3:], "/")

	return &RepoRef{
		URL:  url,
		Path: path,
		Name: name,
		Tag:  tag,
	}, nil
}

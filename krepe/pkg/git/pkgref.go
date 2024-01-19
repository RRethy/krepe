package git

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrInvalidPkgRef = errors.New("package ref must be in the format github.com/<org>/<repo>[/<path>]@<tag>")

	urlPattern = regexp.MustCompile(`^github\.com/(?P<owner>[^/]+)/(?P<repo>[^/]+)(?P<path>(?:/[^/]+)*)@(?P<tag>.+)$`)
)

type PkgRef struct {
	Owner string
	Repo  string
	Path  []string
	Tag   string
	Name  string
}

func NewPkgRefFromString(pkgRef string) (*PkgRef, error) {
	matches := urlPattern.FindStringSubmatch(pkgRef)
	if matches == nil {
		return nil, fmt.Errorf("%w: does not match %s", ErrInvalidPkgRef, pkgRef)
	}

	owner := matches[1]
	repo := matches[2]
	path := strings.Split(matches[3], "/")[1:]
	if len(path) == 0 {
		path = nil
	}
	tag := matches[4]

	name := repo
	if len(path) > 0 {
		name = path[len(path)-1]
	}

	return &PkgRef{
		Owner: owner,
		Path:  path,
		Repo:  repo,
		Tag:   tag,
		Name:  name,
	}, nil
}

func (r *PkgRef) String() string {
	builder := []string{r.Owner, r.Repo}
	if len(r.Path) > 0 {
		builder = append(builder, r.Path...)
	}
	return "github.com/" + strings.Join(builder, "/") + "@" + r.Tag
}

func (r *PkgRef) URL() string {
	return "https://" + r.String()
}

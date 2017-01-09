package domain

import (
	"log"
	"os/exec"
	"strings"
)

type Source interface {
	Name() string
	Owner() string
}

// PathSource implements Source and retrieves the Git source under the current path
type PathSource struct {
	name  string
	owner string
	path  string
}

func NewPathSource(path string) (Source, error) {
	c := &PathSource{
		path: path,
	}
	c.resolve()
	return c, nil
}

func (c *PathSource) Name() string {
	return c.name
}

func (c *PathSource) Owner() string {
	return c.owner
}

func (c *PathSource) resolve() string {
	// out, err := exec.Command(fmt.Sprintf("cd %s", c.path)).Output()
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		log.Fatal(err)
	}

	url := strings.TrimSpace(string(out))
	urlItems := strings.Split(url, "/")

	c.owner = urlItems[len(urlItems)-2]
	c.name = urlItems[len(urlItems)-1]
	return string(out)
}

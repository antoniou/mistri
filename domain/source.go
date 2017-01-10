package domain

import (
	"fmt"
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
	if err := c.resolve(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *PathSource) Name() string {
	return c.name
}

func (c *PathSource) Owner() string {
	return c.owner
}

func (c *PathSource) resolve() error {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("set -e; cd %s;git remote get-url origin", c.path)).Output()
	if err != nil {
		return fmt.Errorf("Could not find repository under path %s: %s", c.path, err.Error())
	}

	url := strings.TrimSpace(string(out))
	urlItems := strings.Split(url, "/")

	c.owner = urlItems[len(urlItems)-2]
	c.name = urlItems[len(urlItems)-1]
	return nil
}

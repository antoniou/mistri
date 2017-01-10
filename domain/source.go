package domain

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Source interface {
	Name() string
	Owner() string
	Auth() string
}

// PathSource implements Source and retrieves the Git source under the current path
type PathSource struct {
	name  string
	owner string
	path  string
	auth  string
}

func NewPathSource(path string) (Source, error) {
	c := &PathSource{
		path: path,
	}
	if err := c.init(); err != nil {
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

func (c *PathSource) Auth() string {
	return c.auth
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

func (c *PathSource) authenticate() error {
	// !!!!!!!!FIXME!!!!!!!!
	buf := bytes.NewBuffer(nil)
	f, err := os.Open("../.githubaccess")
	if err != nil {
		log.Fatal(err)
		return err
	}

	io.Copy(buf, f)
	c.auth = string(buf.Bytes())
	return nil
}

func (c *PathSource) init() error {
	err := c.resolve()
	if err != nil {
		return err
	}

	err = c.authenticate()
	return err
}

package domain

import (
	"fmt"
	"os/exec"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SourceTestSuite struct {
	suite.Suite
}

func (suite *SourceTestSuite) SetupTest() {

}

func (suite *SourceTestSuite) TestCurrentPathSource() {
	repoPath := "/tmp/temprepo"
	origin := "git@github.com:antoniou/zero2Pipe.git"

	// Create temporary local repo
	exec.Command("bash", "-c", fmt.Sprintf("set -e; git init %s; cd %s; git remote add origin %s", repoPath, repoPath, origin)).Output()
	source, err := NewPathSource(repoPath)

	// Cleanup repo
	exec.Command("bash", "-c", fmt.Sprintf("rm -rf %s", repoPath))

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "zero2Pipe", source.Name())
	assert.Equal(suite.T(), "antoniou", source.Owner())
}

func (suite *SourceTestSuite) TestInvalidPathSource() {
	source, err := NewPathSource("/invalid/path")

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), source)
	assert.Contains(suite.T(), err.Error(), "Could not find repository")
}

// func (suite *SourceTestSuite) TestSourceNotSupported() {
// 	path := "."
// 	defer exec.Command("bash", "-c", fmt.Sprintf("cd %s; git remote remove deleteme", path)).Output()
// 	source, err := NewPathSource(".")
//
// 	assert.NotNil(suite.T(), err)
// 	assert.Nil(suite.T(), source)
// }
//
// func TestSourceTestSuite(t *testing.T) {
// 	suite.Run(t, new(SourceTestSuite))
// }

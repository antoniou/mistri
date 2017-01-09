package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SourceTestSuite struct {
	suite.Suite
}

func (suite *SourceTestSuite) SetupTest() {

}

func (suite *SourceTestSuite) TestCurrentPathSource() {
	source, err := NewPathSource(".")

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "zero2Pipe", source.Name())
	assert.Equal(suite.T(), "antoniou", source.Owner())
}

func TestSourceTestSuite(t *testing.T) {
	suite.Run(t, new(SourceTestSuite))
}

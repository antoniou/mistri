package domain

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Define a mock struct to be used in your unit tests of myFunc.
type mockLambdaClient struct {
	mock.Mock
}

type LambdaActorTestSuite struct {
	suite.Suite
	service mockLambdaClient
}

func (suite *LambdaActorTestSuite) SetupTest() {
	suite.service = mockLambdaClient{}
}

func (suite *LambdaActorTestSuite) TestCreateFunction() {

}

func TestLambdaActorTestSuite(t *testing.T) {
	suite.Run(t, new(LambdaActorTestSuite))
}

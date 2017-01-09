package domain

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Define a mock struct to be used in your unit tests of myFunc.
type mockCloudFormationClient struct {
	cloudformationiface.CloudFormationAPI
	mock.Mock
}

func (m *mockCloudFormationClient) CreateStack(i *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	args := m.Mock.Called()
	return nil, args.Error(1)
}

func (m *mockCloudFormationClient) WaitUntilStackCreateComplete(*cloudformation.DescribeStacksInput) error {
	args := m.Mock.Called()
	return args.Error(0)
}

type CloudformationActorTestSuite struct {
	suite.Suite
	service mockCloudFormationClient
}

func (suite *CloudformationActorTestSuite) SetupTest() {
	suite.service = mockCloudFormationClient{}
}

func (suite *CloudformationActorTestSuite) TestCreateStack() {
	actor := CloudFormationActor{
		Template:  "../templates/lambda-store.json",
		StackName: "s3-lambda-bucket",
	}

	suite.service.On("CreateStack").Return(nil, nil)
	suite.service.On("WaitUntilStackCreateComplete").Return(nil)

	err := actor.createStack(&suite.service)
	assert.Nil(suite.T(), err)
	suite.service.AssertCalled(suite.T(), "CreateStack")
	suite.service.AssertCalled(suite.T(), "WaitUntilStackCreateComplete")
	suite.service.AssertExpectations(suite.T())
}

func (suite *CloudformationActorTestSuite) TestCreateStackInvalidTemplateFile() {
	actor := CloudFormationActor{
		Template:  "../templates/invalid-file.json",
		StackName: "s3-lambda-bucket",
	}
	err := actor.createStack(&suite.service)

	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "no such file")
	suite.service.AssertNotCalled(suite.T(), "CreateStack")
	suite.service.AssertExpectations(suite.T())
}

func (suite *CloudformationActorTestSuite) TestCreateStackReturnsWithErrors() {
	actor := CloudFormationActor{
		Template:  "../templates/lambda-store.json",
		StackName: "s3-lambda-bucket",
	}

	suite.service.On("CreateStack").Return(nil, fmt.Errorf("500 Server failure"))
	err := actor.createStack(&suite.service)

	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "500")
	suite.service.AssertCalled(suite.T(), "CreateStack")
	suite.service.AssertNotCalled(suite.T(), "WaitUntilStackCreateComplete")
	suite.service.AssertExpectations(suite.T())
}

func (suite *CloudformationActorTestSuite) TestWaitUntilStackCreateCompleteReturnsWithErrors() {
	actor := CloudFormationActor{
		Template:  "../templates/lambda-store.json",
		StackName: "s3-lambda-bucket",
	}

	suite.service.On("CreateStack").Return(nil, nil)
	suite.service.On("WaitUntilStackCreateComplete").Return(fmt.Errorf("500 Server failure"))
	err := actor.createStack(&suite.service)

	assert.NotNil(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "500")
	suite.service.AssertCalled(suite.T(), "CreateStack")
	suite.service.AssertCalled(suite.T(), "WaitUntilStackCreateComplete")
	suite.service.AssertExpectations(suite.T())
}

func TestCloudformationActorTestSuite(t *testing.T) {
	suite.Run(t, new(CloudformationActorTestSuite))
}

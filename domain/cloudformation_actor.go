package domain

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

// CloudFormationActor implements the Actor interface and provisions a
// cloudformation stack
type CloudFormationActor struct {
	Template   string
	StackName  string
	Parameters map[string]string
}

func (c *CloudFormationActor) Run(interface{}) error {
	session, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return err
	}

	svc := cloudformation.New(session)
	err = c.createStack(svc)
	if err != nil {
		log.Printf("Failed to create stack: %s", err)
	}

	return err
}

func (c *CloudFormationActor) createStackInput() (*cloudformation.CreateStackInput, error) {
	template, err := c.templateContents()
	if err != nil {
		return nil, err
	}

	parameters := make([]*cloudformation.Parameter, 0)
	for pkey, pvalue := range c.Parameters {
		parameters = append(parameters, &cloudformation.Parameter{
			ParameterKey:     aws.String(pkey),
			ParameterValue:   aws.String(pvalue),
			UsePreviousValue: aws.Bool(true),
		})
	}

	stackInput := &cloudformation.CreateStackInput{
		StackName: aws.String(c.StackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_NAMED_IAM"),
		},
		TemplateBody: aws.String(template),
		Parameters:   parameters,
	}

	return stackInput, nil

}

func (c *CloudFormationActor) createStack(service cloudformationiface.CloudFormationAPI) error {
	params, err := c.createStackInput()
	if err != nil {
		fmt.Printf("Error creating Stack Input Parameters: %s", err.Error())
		return err
	}

	resp, err := service.CreateStack(params)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return err
	}

	dparams := &cloudformation.DescribeStacksInput{
		StackName: aws.String(c.StackName),
	}
	err = service.WaitUntilStackCreateComplete(dparams)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// Pretty-print the response data.
	fmt.Println(resp)
	return nil
}

func (c *CloudFormationActor) templateContents() (string, error) {
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(c.Template)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	io.Copy(buf, f)
	return string(buf.Bytes()), nil
}

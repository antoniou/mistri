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
	Template  string
	StackName string
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

func (c *CloudFormationActor) createStack(service cloudformationiface.CloudFormationAPI) error {
	template, err := c.templateContents()
	if err != nil {
		return err
	}
	params := &cloudformation.CreateStackInput{
		StackName: aws.String(c.StackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_NAMED_IAM"),
		},
		TemplateBody: aws.String(template),
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

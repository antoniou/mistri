package domain

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type AWSCodePipelineContext struct {
	Props   map[string]string
	Session *session.Session
}

func (c *AWSCodePipelineContext) init() error {
	var err error
	c.Session, err = session.NewSession()
	if err != nil {
		fmt.Println("Failed to create AWS session,", err)
		return err
	}

	c.Props["region"] = *c.Session.Config.Region
	account, ok := c.account()
	c.Props["account"] = account
	return ok
}

func (c *AWSCodePipelineContext) account() (string, error) {
	svc := iam.New(c.Session)
	resp, err := svc.GetUser(&iam.GetUserInput{})
	if err != nil {
		log.Fatal("Failed to retrieve user,", err)
		return "", err
	}

	return strings.Split(*resp.User.Arn, ":")[4], nil

}

func NewAWSCodePipelineContext() *AWSCodePipelineContext {
	context := AWSCodePipelineContext{
		Props: make(map[string]string),
	}
	context.init()
	return &context
}

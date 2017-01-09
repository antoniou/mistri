package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	cli := New()

	assert.NotNil(t, cli)
	assert.IsType(t, new(Client), cli)
}

func TestRegisterCommands(t *testing.T) {
	cli := &Client{}

	err := cli.RegisterCommands()

	assert.Nil(t, err)
	assert.IsType(t, new(CreateCommand), cli.commands["create"])

	// Verify that commands are not registered twice
	err = cli.RegisterCommands()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

// type AppClientMock struct {
// 	mock.Mock
// }

// func TestRunClient(t *testing.T) {
// 	cli, _ := New()
// 	appMock := new(AppClientMock)
// 	arguments := []string{"create", "help"}
// 	appMock.On("Run", arguments).Return(nil)
//
// 	cli.Run(arguments)
//
// 	appMock.AssertNumberOfCalls(t, "Run", 2)
// }

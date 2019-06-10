package cmd_test

import (
	"bytes"
	"testing"

	"github.com/reaandrew/maul/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err

}

func TestHelloWorld(t *testing.T) {
	output, err := executeCommand(cmd.RootCmd)

	assert.Nil(t, err, "Unexpected error: %v", err)
	assert.Equal(t, "Hello, World!\n", output)
}

package cmd_test

import (
	"bytes"
	"os"
	"regexp"
	"testing"

	"github.com/reaandrew/schmokin/cmd"
	"github.com/reaandrew/schmokin/utils"
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

func TestOutput(t *testing.T) {
	file := utils.CreateTestFile([]string{
		"http://localhost:8080/1",
	})
	defer os.Remove(file.Name())

	output, err := executeCommand(cmd.RootCmd, "-u", file.Name(), "-n", "1", "-c", "1")
	assert.Nil(t, err)

	patterns := []string{
		`Random[^\s]+\s[\w]+`,
		`Worker Count[^\s]+\s[\d]+`,
		`Successful Transactions[^\s]+\s[\d]+`,
		`Failed Transactions[^\s]+\s[\d]+`,
		`Concurrency[^\s]+\s[\d\.]+`,
		`Shortest Transaction[^\s]+\s[\d]+s`,
		`Longest Transaction[^\s]+\s[\d]+s`,
		`Elapsed Time \(ms\)[^\s]+\s[\d]+s`,
		`Availability \(%\)[^\s]+\s[\d]+`,
		`Transactions[^\s]+\s[\d]+`,
		`Data Receive Rate \(bytes/sec\)[^\s]+\s[\d]+ B`,
		`Data Send Rate \(bytes/sec\)[^\s]+\s[\d]+ B`,
		`Total Bytes Sent[^\s]+\s[\d]+ B`,
		`Total Bytes Received[^\s]+\s[\d]+ B`,
		`Average Transaction Rate \(requests/sec\)[^\s]+\s[^0][\d\.]+`,
		`Average Response Time \(ms\)[^\s]+\s[\d\.]+`,
	}

	for _, pattern := range patterns {
		matched, err := regexp.Match(pattern, []byte(output))
		assert.Nil(t, err)
		assert.True(t, matched, pattern)
	}
}

package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_jsonFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./panalysis -json", " ")

	status := cli.Run(args)
	_ = status
}

func TestRun_configFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./panalysis -config", " ")

	status := cli.Run(args)
	_ = status
}

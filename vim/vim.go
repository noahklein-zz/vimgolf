package vim

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func AttemptChallenge(ctx context.Context, challengeText string, command string) (string, error) {
	tf, err := tempFileWithContents("", "challenge", []byte(challengeText))
	if err != nil {
		return "", err
	}
	fileName := tf.Name()
	defer os.Remove(fileName)

	if err := runVimCommand(ctx, command, fileName); err != nil {
		return "", fmt.Errorf("Vim command failed to run %v:%v", command, err)
	}

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("Failed to read temp file %s: %v", fileName, err)
	}
	return string(b), nil
}

func tempFileWithContents(dir, prefix string, b []byte) (*os.File, error) {
	tf, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return nil, fmt.Errorf("Failed to create temp file: %v", err)
	}
	defer tf.Close()

	if _, err = tf.Write(b); err != nil {
		return nil, fmt.Errorf("Failed to write to temp file %v: %v", tf, err)
	}

	return tf, nil
}

func runVimCommand(ctx context.Context, command string, file string) error {
	cs := []string{"-nZ", "-u", "NONE"}
	for _, command := range lexCommand(command) {
		if strings.Trim(command, " ") == "" {
			continue
		}
		cs = append(cs, parseCommand(command))
	}
	cs = append(cs, "+wq", file)

	cmd := exec.CommandContext(ctx, "ex", cs...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%v: %s", err, stderr.String())
	}
	return nil
}

func prepareCommand(cmd string) []string {
	lexemes := lexCommand(cmd)
	var out []string
	for _, c := range lexemes {
		out = append(out, parseCommand(c))
	}
	return out
}

func parseCommand(cmd string) string {
	cmd = strings.TrimLeft(cmd, "+")
	if strings.HasPrefix(cmd, ":") {
		return fmt.Sprintf("+%s", cmd)
	}
	return fmt.Sprintf("+normal %s", cmd)
}

func lexCommand(cmd string) []string {
	return strings.Split(cmd, "<Esc>")
}

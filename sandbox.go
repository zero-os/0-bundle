package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
)

const (
	//EnvFile default location of env file
	EnvFile = "/etc/env"

	//BufferSize of data captured from the sandbox stdout stderr (tail)
	BufferSize = 32 * 1024 // 32KB
)

func parseEnv(r io.Reader) ([]string, error) {
	reader := bufio.NewReader(r)
	var env []string
	var err error
	for err == nil {
		var line string
		line, err = reader.ReadString('\n')
		log.Debugf("line: %s", line)
		if err != nil && err != io.EOF {
			return env, err
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		//TODO: some env line validation goes here
		env = append(env, line)
	}

	return env, nil
}

func environ(root string) ([]string, error) {

	name := path.Join(root, EnvFile)

	log.Debugf("opening file: %s", name)
	file, err := os.Open(name)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	var env []string

	for err == nil {
		var line string
		line, err = reader.ReadString('\n')
		log.Debugf("line: %s", line)
		if err != nil && err != io.EOF {
			return env, err
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		//TODO: some env line validation goes here
		env = append(env, line)
	}

	return env, nil
}

//sandbox, runs the sandbox and return the captured stdout, stderr, and exit error
func sandbox(root string, userenv []string) ([]byte, []byte, error) {
	//read env
	log.Debugf("reading the env")
	flistenv, err := environ(root)
	if err != nil {
		return nil, nil, err
	}

	stdout := NewTailBuffer(BufferSize)
	stderr := NewTailBuffer(BufferSize)

	//start
	cmd := exec.Cmd{
		Path: "/etc/start",
		Dir:  "/",
		Env:  append(flistenv, userenv...),
		SysProcAttr: &syscall.SysProcAttr{
			Chroot: root,
		},
		Stdout: io.MultiWriter(os.Stdout, stdout),
		Stderr: io.MultiWriter(os.Stderr, stderr),
	}

	err = cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

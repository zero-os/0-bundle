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
	EnvFile = "/etc/env"
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

func sandbox(root string, userenv []string) error {

	//read env
	log.Debugf("reading the env")
	flistenv, err := environ(root)
	if err != nil {
		return err
	}

	//start
	cmd := exec.Cmd{
		Path: "/etc/start",
		Dir:  "/",
		Env:  append(flistenv, userenv...),
		SysProcAttr: &syscall.SysProcAttr{
			Chroot: root,
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	return cmd.Run()
}

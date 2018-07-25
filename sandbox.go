package main

import (
	"bufio"
	"fmt"
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

type Sandbox struct {
	Root       string
	UserEnv    []string
	EntryPoint string
	Args       []string

	cmd *exec.Cmd
}

//Run runs the sandbox, blocks until the sandbox
//entry point exits
func (s *Sandbox) Run() ([]byte, []byte, error) {
	log.Debugf("reading the env")
	flistenv, err := environ(s.Root)
	if err != nil {
		return nil, nil, err
	}

	stdout := NewTailBuffer(BufferSize)
	stderr := NewTailBuffer(BufferSize)

	args := []string{s.EntryPoint}
	args = append(args, s.Args...)
	//start
	cmd := exec.Cmd{
		Path: s.EntryPoint,
		Dir:  "/",
		Args: args,
		Env:  append(flistenv, s.UserEnv...),
		SysProcAttr: &syscall.SysProcAttr{
			Chroot: s.Root,
		},
		Stdout: io.MultiWriter(os.Stdout, stdout),
		Stderr: io.MultiWriter(os.Stderr, stderr),
	}
	s.cmd = &cmd
	err = cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

//Signal the sandbox process
func (s *Sandbox) Signal(signal os.Signal) error {
	if s.cmd == nil || s.cmd.Process == nil {
		return fmt.Errorf("sandbox is not started")
	}

	return s.cmd.Process.Signal(signal)
}

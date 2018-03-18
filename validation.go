package main

import (
	"os"
)

func isRoot() bool {
	return os.Getuid() == 0
}

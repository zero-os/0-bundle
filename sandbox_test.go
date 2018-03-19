package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEnv(t *testing.T) {
	r := strings.NewReader(
		`
# some comments
ENV=1
NAME=AZMY

AGE=100
`)

	out, err := parseEnv(r)

	if ok := assert.Nil(t, err); !ok {
		t.Fatal()
	}

	if ok := assert.Len(t, out, 3); !ok {
		t.Error()
	}
}

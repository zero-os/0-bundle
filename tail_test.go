package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTailBufferNotFull(t *testing.T) {
	buf := NewTailBuffer(10)
	x := "012345" // exactly 6
	c, err := buf.Write([]byte(x))

	if ok := assert.Nil(t, err); !ok {
		t.Fatal()
	}

	if ok := assert.Equal(t, c, len(x)); !ok {
		t.Error()
	}

	res := buf.Bytes()

	if ok := assert.Equal(t, x, string(res)); !ok {
		t.Error()
	}
}

func TestTailBufferOverflow(t *testing.T) {
	buf := NewTailBuffer(10)
	x := "0123456789ABC" // exactly 13
	c, err := buf.Write([]byte(x))

	if ok := assert.Nil(t, err); !ok {
		t.Fatal()
	}

	if ok := assert.Equal(t, c, len(x)); !ok {
		t.Error()
	}

	res := buf.Bytes()

	if ok := assert.Equal(t, "3456789ABC", string(res)); !ok {
		t.Error()
	}
}

func TestTailBufferZeroSize(t *testing.T) {
	buf := NewTailBuffer(0)
	x := "0123456789ABC" // exactly 13
	_, err := buf.Write([]byte(x))

	if ok := assert.NotNil(t, err); !ok {
		t.Fatal()
	}

}

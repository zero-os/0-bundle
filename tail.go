package main

import (
	"syscall"
)

//TailBuffer keeps a fixed sized buffer of the last
//writen bytes, older bytes are discarded
type TailBuffer struct {
	buf  []byte
	cur  int
	full bool
}

//Write writes bytes to buffer
func (t *TailBuffer) Write(p []byte) (int, error) {
	dst := t.buf[t.cur:]
	//fmt.Println(cap(dst))
	if cap(dst) > len(p) {
		copy(dst, p)
		t.cur += len(p)
		return len(p), nil
	}

	n := cap(dst)
	if n == 0 {
		return 0, syscall.ENOSPC
	}
	copy(dst, p[:n])
	t.cur = 0
	t.full = true
	m, err := t.Write(p[n:])
	n += m
	return n, err
}

//Bytes returns the last writen bytes up to size
func (t *TailBuffer) Bytes() []byte {
	if t.full {
		//we need to read from cur to end, then rotate
		result := make([]byte, len(t.buf))
		copy(result, t.buf[t.cur:])
		copy(result[len(t.buf)-t.cur:], t.buf)

		return result
	}
	result := make([]byte, t.cur)
	copy(result, t.buf)
	return result
}

//NewTailBuffer creates a new tail buffer
func NewTailBuffer(s int) *TailBuffer {
	return &TailBuffer{buf: make([]byte, s)}
}

package lockfile

import (
	"testing"
)

var lockfileName = "pid/test.pid"

func TestLockfile(t *testing.T) {
	l, err := NewLockfile(lockfileName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = l.Get()
	if err == nil {
		t.Fatal(l)
	}
	err = l.Lock()
	if err != nil {
		t.Fatal(err)
	}
	err = l.Lock()
	if err != nil {
		t.Fatal(err)
	}
	_, err = l.Get()
	if err != nil {
		t.Fatal(err)
	}
	err = l.Unlock()
	if err != nil {
		t.Fatal(err)
	}
}

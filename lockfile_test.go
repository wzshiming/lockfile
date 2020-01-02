package lockfile

import (
	"testing"
)

var lockfileName = "test.pid"

func TestLockfile(t *testing.T) {
	l, err := NewLockfile(lockfileName)
	if err != nil {
		t.Fatal(l)
	}
	_, err = l.Get()
	if err == nil {
		t.Fatal(l)
	}
	err = l.Lock()
	if err != nil {
		t.Fatal(l)
	}
	err = l.Lock()
	if err == nil {
		t.Fatal(l)
	}
	_, err = l.Get()
	if err != nil {
		t.Fatal(l)
	}
	err = l.Unlock()
	if err != nil {
		t.Fatal(l)
	}
}

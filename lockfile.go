package lockfile

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// lockfile is a lock file
type lockfile string

var (
	errNotExist  = errors.New("error lockfile not exist")
	errIsRunning = errors.New("error is running")
)

// NewLockfile returns a new lockfile.
func NewLockfile(path string) (lockfile, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return lockfile(path), nil
}

// put the pid to lockfile.
func (p lockfile) put() error {
	lockfile := string(p)
	if p == "" {
		return errNotExist
	}
	err := os.MkdirAll(filepath.Dir(lockfile), os.FileMode(0755))
	if err != nil {
		return err
	}

	pid := os.Getpid()

	err = ioutil.WriteFile(lockfile, []byte(strconv.Itoa(pid)), os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("error write lockfile %s: %s", lockfile, err.Error())
	}

	return nil
}

// Get the pid from the lockfile.
func (p lockfile) Get() (int, error) {
	lockfile := string(p)
	if lockfile == "" {
		return 0, errNotExist
	}

	d, err := ioutil.ReadFile(lockfile)
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(string(bytes.TrimSpace(d)))
	if err != nil {
		return 0, fmt.Errorf("error read lockfile %s: %s", lockfile, err.Error())
	}

	return pid, nil
}

// Lock the lockfile.
func (p lockfile) Lock() error {
	pid, err := p.Get()
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		ok, err := isRunning(pid)
		if err != nil {
			return err
		}
		if ok {
			return errIsRunning
		}
	}
	return p.put()
}

// Unlock the lockfile.
func (p lockfile) Unlock() error {
	pid, err := p.Get()
	if err != nil {
		return nil
	}
	ok, err := isRunning(pid)
	if err != nil {
		return nil
	}
	if !ok {
		return nil
	}

	mpid := os.Getpid()
	if mpid != pid {
		return errIsRunning
	}
	return os.Remove(string(p))
}

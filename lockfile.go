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

// Lockfile is a lock file
type Lockfile string

var (
	errNotExist  = errors.New("error Lockfile not exist")
	errIsRunning = errors.New("error is running")
)

// NewLockfile returns a new Lockfile.
func NewLockfile(path string) (Lockfile, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return Lockfile(path), nil
}

// put the pid to Lockfile.
func (p Lockfile) put() error {
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
		return fmt.Errorf("error write Lockfile %s: %s", lockfile, err.Error())
	}

	return nil
}

// Get the pid from the Lockfile.
func (p Lockfile) Get() (int, error) {
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
		return 0, fmt.Errorf("error read Lockfile %s: %s", lockfile, err.Error())
	}

	return pid, nil
}

// Lock the Lockfile.
func (p Lockfile) Lock() error {
	pid, err := p.Get()
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		mpid := os.Getpid()
		if mpid == pid {
			return nil
		}
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

// Unlock the Lockfile.
func (p Lockfile) Unlock() error {
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

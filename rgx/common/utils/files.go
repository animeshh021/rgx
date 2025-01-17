package utils

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"rgx/common/log"
)

func RemoveDir(dir string) error {
	e1 := removeContents(dir)
	if e1 != nil {
		return e1
	}
	e2 := os.Remove(dir)
	if e2 != nil {
		return e2
	}
	return nil
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer func(d *os.File) {
		err := d.Close()
		if err != nil {
			log.Error("removeContents: could not close file")
		}
	}(d)
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func Hash(filename, algorithm string) (Checksum, error) {
	f, e := os.Open(filename)
	if e != nil {
		return Checksum{"", ""}, e
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Error("could not close downloaded file after hashing: %s", err.Error())
		}
	}(f)

	var hasher hash.Hash
	switch algorithm {
	case "sha1":
		hasher = sha1.New()
	default:
		hasher = sha256.New()
	}
	if _, e := io.Copy(hasher, f); e != nil {
		return Checksum{"", ""}, e
	}
	return Checksum{
			Algorithm: algorithm,
			Hash:      hex.EncodeToString(hasher.Sum(nil)),
		},
		nil
}

func Copy(src, dest string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func(source *os.File) {
		err := source.Close()
		if err != nil {
			log.Warn("could not close: %s: %s", src, err.Error())
		}
	}(source)

	destination, err := os.Create(dest)
	if err != nil {
		return 0, err
	}
	defer func(destination *os.File) {
		err := destination.Close()
		if err != nil {
			log.Warn("could not close: %s: %s", dest, err.Error())
		}
	}(destination)
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func ExtractToTemp(source, artifactType string) (string, error) {
	if !Exists(source) {
		return "", errors.New("file not found: " + source)
	}
	tempdir, err := os.MkdirTemp(TempDir(), artifactType+string('-'))
	if err != nil {
		return "", err
	}
	Extract(source, tempdir)
	return tempdir, nil
}

func CleanDirs(fn func() []string) {
	dirs := fn()
	for _, dir := range dirs {
		err := RemoveDir(dir)
		if err != nil {
			log.Warn("could not clean up %s", dir)
		}
	}
}

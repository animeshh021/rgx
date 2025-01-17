package utils

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"rgx/common/log"
)

func Extract(archiveName, targetDir string) bool {

	log.Trace("starting to decompress %s ...", archiveName)

	if strings.HasSuffix(archiveName, ".tar.gz") {
		// tar zxf the downloaded xyz.tar.gz to targetDir
		tgz, e := os.Open(archiveName)
		ErrCheck(e)
		untar(tgz, targetDir)
		e = tgz.Close()
		ErrCheck(e)
	} else if strings.HasSuffix(archiveName, ".zip") {
		// unzip the archive, typically on windows
		e := unzip(archiveName, targetDir)
		ErrCheck(e)
	}
	if !Exists(targetDir) {
		log.Trace("tried to extract arhive to %s, but it doesnt exist ", targetDir)
		return false
	} else {
		log.Trace("extracted archive to %s", targetDir)
		return true
	}
}

func untar(gzipStream io.Reader, targetDir string) {
	cwd := getwd()
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Fatal("extract: failed to open gzip stream")
	}

	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("extract: failed to open tar header: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filepath.Join(targetDir, header.Name), 0755); err != nil {
				log.Fatal("extract: failed to create directory: %s", err.Error())
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Join(targetDir, header.Name), 0755); err != nil {
				log.Fatal("extract: failed to create parent directories: %s", err.Error())
			}
			outFile, err := os.Create(filepath.Join(targetDir, header.Name))
			if err != nil {
				log.Fatal("extract: failed to create file: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatal("extract: failed to copy data: %s", err.Error())
			}
			closeError := outFile.Close()
			if closeError != nil {
				log.Debug("error closing file")
			}
			if header.Mode&0111 != 0 { // has executable bit set
				e := os.Chmod(filepath.Join(targetDir, header.Name), os.FileMode(header.Mode))
				if e != nil {
					log.Warn("could not set permissions: %v", err)
				}
			}
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(filepath.Join(targetDir, header.Name)), 0755); err != nil {
				log.Fatal("extract: failed to create parent directories: %s", err.Error())
			}
			err := os.Symlink(header.Linkname, filepath.Join(targetDir, header.Name))
			if err != nil {
				log.Error("extract: error creating symlink: %v", err)
			}
		default:
			log.Fatal("extract: unknown type: %v in %s", header.Typeflag, header.Name)
		}
	}
	err = os.Chdir(cwd)
	ErrCheck(err)
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	err = os.MkdirAll(dest, 0755)
	ErrCheck(err)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(path, f.Mode())
			ErrCheck(err)
		} else {
			err := os.MkdirAll(filepath.Dir(path), f.Mode())
			ErrCheck(err)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func getwd() string {
	cwd, e := os.Getwd()
	ErrCheck(e)
	return cwd
}

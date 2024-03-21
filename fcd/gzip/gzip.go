package gzip

import (
	gz "compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"fcd/path"
)

func Compress(source, target string) error {
	source = path.ExpandDirAndEnv(source)
	target = path.ExpandDirAndEnv(target)

	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	archiver := gz.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	return err
}

func Uncompress(source, target string) error {
	source = path.ExpandDirAndEnv(source)
	target = path.ExpandDirAndEnv(target)

	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gz.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

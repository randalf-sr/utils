package archive

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fcd/path"
)

func UnTar(tarball, target string) error {
	tarball = path.ExpandDirAndEnv(tarball)
	target = path.ExpandDirAndEnv(target)

	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()

	stripped := filepath.Base(tarball)
	_ = os.MkdirAll(filepath.Join(target, strings.TrimSuffix(stripped, filepath.Ext(stripped))), 0755)

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		// sometimes tar files have a directory entry before the file entry and it has not come across the dir, so creating each dir regardless
		_ = os.MkdirAll(filepath.Dir(path), 0755)

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(file, tarReader)
		_ = file.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func Tar(source, target string) error {
	source = path.ExpandDirAndEnv(source)
	target = path.ExpandDirAndEnv(target)

	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarFile, err := os.Create(target)
	if err != nil {
		return err
	}

	defer tarFile.Close()

	archive := tar.NewWriter(tarFile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := archive.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}

			defer file.Close()
			_, err = io.Copy(archive, file)

			return err
		})
}

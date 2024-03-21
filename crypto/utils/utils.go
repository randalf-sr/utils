package utils

import (
	"archive/tar"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tyler-smith/go-bip39"
)

func EncryptInto(fromPlainText io.Reader, toCipherText io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, block.BlockSize())
	_, err = rand.Read(iv)
	if err != nil {
		return err
	}

	_, err = toCipherText.Write(iv)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	cipherWriter := &cipher.StreamWriter{S: stream, W: toCipherText}

	archiver := gzip.NewWriter(cipherWriter)
	defer archiver.Close()

	_, err = io.Copy(archiver, fromPlainText)

	return err
}

func EncryptFile(filename string, key []byte) error {
	filename = ExpandHomeDir(filename)
	plainfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer plainfile.Close()

	cipherfile, err := os.Create(filename + ".encrypted")
	if err != nil {
		return err
	}
	defer cipherfile.Close()

	return EncryptInto(plainfile, cipherfile, key)
}

func DecryptInto(fromCipherText io.Reader, toPlainText io.Writer, key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := make([]byte, block.BlockSize())
	_, err = fromCipherText.Read(iv)
	if err != nil {
		return err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	cipherReader := &cipher.StreamReader{S: stream, R: fromCipherText}

	archive, err := gzip.NewReader(cipherReader)
	if err != nil {
		return err
	}
	defer archive.Close()

	_, err = io.Copy(toPlainText, archive)

	return err
}

func DecryptFile(filename string, key []byte) error {
	filename = ExpandHomeDir(filename)

	cipherfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer cipherfile.Close()

	plainfile, err := os.Create(filename[:len(filename)-10])
	if err != nil {
		return err
	}
	defer plainfile.Close()

	return DecryptInto(cipherfile, plainfile, key)
}

func UnTar(tarball, target string) error {
	tarball = ExpandHomeDir(tarball)
	target = ExpandHomeDir(target)

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
	source = ExpandHomeDir(source)
	target = ExpandHomeDir(target)

	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	tarball := tar.NewWriter(tarFile)
	defer tarball.Close()

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

			if err := tarball.WriteHeader(header); err != nil {
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
			_, err = io.Copy(tarball, file)
			return err
		})
}

func Gzip(source, target string) error {
	source = ExpandHomeDir(source)
	target = ExpandHomeDir(target)

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

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	return err
}

func UnGzip(source, target string) error {
	source = ExpandHomeDir(source)
	target = ExpandHomeDir(target)

	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
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

func ExpandHomeDirAndEnv(path string) string {
	path = Ignore(path, TryExpandHomeDir)
	return os.ExpandEnv(path)
}

func ExpandHomeDir(path string) string {
	return Ignore(path, TryExpandHomeDir)
}

func TryExpandHomeDir(path string) (expandedPath string, expanded bool) {
	if !strings.HasPrefix(path, "~") {
		return path, false
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return home + path[1:], true
}

func Ignore[T any, R1 any, R2 any](arg T, fn func(arg T) (R1, R2)) R1 {
	r, _ := fn(arg)
	return r
}

func GetKeyForMnemonic(mnemonicText, secret string) []byte {
	mnemonicText = strings.ReplaceAll(mnemonicText, "\r", "")
	mnemonicText = strings.ReplaceAll(mnemonicText, "\n", "")
	mnemonicText = strings.ReplaceAll(mnemonicText, "\t", " ")

	parts := strings.Split(mnemonicText, " ")
	mnemonicText = ""
	for _, p := range parts {
		if len(p) > 0 {
			mnemonicText += p + " "
		}
	}

	seed := bip39.NewSeed(strings.TrimSpace(mnemonicText), strings.TrimSpace(secret))
	hash := sha256.Sum256(seed)

	return hash[:]
}

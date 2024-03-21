package crypto

import (
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

	"fcd/path"

	"github.com/tyler-smith/go-bip39"
)

func EncryptFile(sourceFile, targetFile, mnemonic, secret string) error {
	key := GetBytesFromMnemonicFile(mnemonic, secret)
	sourceFile = path.ExpandDirAndEnv(sourceFile)
	targetFile = path.ExpandDirAndEnv(targetFile)
	plainfile, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer plainfile.Close()

	if path.Exists(targetFile) {
		return fmt.Errorf("target file '%s' already exists", targetFile)
	}

	cipherfile, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer cipherfile.Close()

	return encryptInto(plainfile, cipherfile, key)
}

func DecryptFile(sourceFile string, targetFile, mnemonic string, secret string) error {
	key := GetBytesFromMnemonicFile(mnemonic, secret)
	sourceFile = path.ExpandDirAndEnv(sourceFile)

	cipherfile, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer cipherfile.Close()

	if path.Exists(targetFile) {
		return fmt.Errorf("target file '%s' already exists", targetFile)
	}

	plainfile, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer plainfile.Close()

	return decryptInto(cipherfile, plainfile, key)
}

func GetBytesFromMnemonicFile(mnemonic, secret string) []byte {
	path := getPath(mnemonic)
	if path == nil {
		fmt.Printf("Mnemonic file '%s' not found\n", mnemonic)
		os.Exit(-1)
	}

	data, err := os.ReadFile(*path)
	if err != nil {
		fmt.Printf("Error reading mnemonic file '%s': %v\n", *path, err)
		os.Exit(-1)
	}
	return GetBytesForMnemonicText(string(data), secret)
}

func GetBytesForMnemonicText(mnemonicText, secret string) []byte {
	seed := bip39.NewSeed(sanitizeMnemonic(mnemonicText), strings.TrimSpace(secret))
	hash := sha256.Sum256(seed)

	return hash[:]
}

func encryptInto(fromPlainText io.Reader, toCipherText io.Writer, key []byte) error {
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

func decryptInto(fromCipherText io.Reader, toPlainText io.Writer, key []byte) error {
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

func getPath(mnemonic string) *string {
	mnemonicPath := strings.TrimSpace(os.ExpandEnv(mnemonic))
	if path.Exists(mnemonicPath) {
		return &mnemonicPath
	}

	mnemonicPath = path.ExpandDir(filepath.Join("~/.mnemonic", mnemonicPath))
	if path.Exists(mnemonicPath) {
		return &mnemonicPath
	}

	return nil
}

func sanitizeMnemonic(mnemonic string) string {
	mnemonic = strings.ReplaceAll(mnemonic, "\r", "")
	mnemonic = strings.ReplaceAll(mnemonic, "\n", "")
	mnemonic = strings.ReplaceAll(mnemonic, "\t", " ")

	parts := strings.Split(mnemonic, " ")
	mnemonic = ""
	for _, p := range parts {
		if len(p) > 0 {
			mnemonic += p + " "
		}
	}

	return strings.TrimSpace(mnemonic)
}

package main

import (
	"os"
	"testing"
)

const (
	test_CLEAR_TEXT = "x882gc4lgpIP9NLaNru2ox3yTN7d8cSj2vhcJz3CxY6xZsSb9m"
	test_SECRET     = "f8aSIzuyKgcYx1oevgJCj5Ct3lIeLM2asC2a3pM2B4TkH5UM"
	test_MNEMONIC   = "xray sierra golf romeo topaz alpha papa quebec delta echo victor yankee"
)

func TestEncryptionDecryption(t *testing.T) {
	var files []string
	defer func() {
		for _, file := range files {
			os.Remove(file)
		}
	}()

	mnemonic = writeToTmpFile(t, "mnemonic_*.txt", test_MNEMONIC)
	files = append(files, mnemonic)

	sourcePath = writeToTmpFile(t, "data_*.txt", test_CLEAR_TEXT)
	files = append(files, sourcePath)

	targetPath = sourcePath + ".enc"
	files = append(files, targetPath)

	encrypt = true
	secret = test_SECRET

	err := run()
	if err != nil {
		t.Errorf("Error encrypting file: %v", err)
	}

	encrypt = false
	decrypt = true
	sourcePath = targetPath
	targetPath = sourcePath + ".dec"
	files = append(files, targetPath)

	err = run()
	if err != nil {
		t.Errorf("Error decrypting file: %v", err)
	}

	decryptedString := readStringFromFile(t, targetPath)
	if decryptedString != test_CLEAR_TEXT {
		t.Errorf("Decrypted string does not match original string: %v", decryptedString)
	}
}

func getSize(t *testing.T, path string) int64 {
	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Errorf("Error getting file info: %v", err)
	}

	return fileInfo.Size()
}

func writeToTmpFile(t *testing.T, pattern string, value string) string {
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Errorf("Error creating temp file: %v", err)
	}

	_, err = file.WriteString(value)
	if err != nil {
		t.Errorf("Error writing to temp file: %v", err)
	}

	return file.Name()
}

func readStringFromFile(t *testing.T, path string) string {
	value, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	return string(value)
}

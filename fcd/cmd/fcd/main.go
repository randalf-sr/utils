package main

import (
	"flag"
	"fmt"
	"os"

	"fcd/crypto"
)

var (
	sourcePath string
	targetPath string
	mnemonic   string
	secret     string
	encrypt    bool
	decrypt    bool
	help       bool
)

func main() {
	parseArgs()

	err := run()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(-1)
	}
}

func run() error {
	if encrypt {
		return crypto.EncryptFile(sourcePath, targetPath, mnemonic, secret)
	} else {
		return crypto.DecryptFile(sourcePath, targetPath, mnemonic, secret)
	}
}

func parseArgs() {
	flag.StringVar(&sourcePath, "s", "", "source file path.")
	flag.StringVar(&targetPath, "t", "", "target file path.")
	flag.StringVar(&mnemonic, "m", "", "required path to the mnemonic file. Can also be an OS env variable. If not found, an attempt will be made to locate it under ~/.mnemonic.")
	flag.StringVar(&secret, "p", "", "password/secret to use with the mnemonic for encrypting/decrypting the file.")
	flag.BoolVar(&encrypt, "e", false, "encrypt the file.")
	flag.BoolVar(&decrypt, "d", false, "decrypt the file.")

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if (encrypt && decrypt) || (!encrypt && !decrypt) {
		flag.Usage()
		fmt.Println("either encrypt or decrypt must be specified.")
		os.Exit(-1)
	}

	if sourcePath == "" {
		flag.Usage()
		fmt.Println("source file path must be specified.")
		os.Exit(-1)
	}

	if mnemonic == "" {
		flag.Usage()
		fmt.Println("mnemonic file path must be specified.")
		os.Exit(-1)
	}

	if secret == "" {
		flag.Usage()
		fmt.Println("secret/password must be specified.")
		os.Exit(-1)
	}

	if targetPath != "" {
		return
	}

	if decrypt {
		flag.Usage()
		fmt.Println("target file must be specified when decrypting.")
		os.Exit(-1)
	}

	fmt.Printf("No target file path specified. Will use '%s.encrypted'.", sourcePath)
	targetPath = sourcePath + ".encrypted"
}

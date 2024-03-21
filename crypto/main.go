package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"rutils/utils"
	"strings"
)

func main() {
	for _, v := range os.Args {
		if v == "-h" || v == "--help" {
			fmt.Println("Example usage for tar, gzip, encrypt, decrypt:")
			fmt.Println("  ## Unzip ##")
			fmt.Println("    > rutils -a -gz -s [your_file].tar.gz -d [your_file].tar")
			fmt.Println("  ## Untar ##")
			fmt.Println("    > rutils -a -tar -s [your_file].tar -d [your_file]")
			fmt.Println("  ## Encrypt ##")
			fmt.Println("    > rutils.exe -a enc -s [your_file] -m [your_mnemonic-key-file] -p [your_secret_key]")
			fmt.Println("  ## Decrypt ##")
			fmt.Println("    > rutils -a -enc -s [your_file].encrypted -m [your_mnemonic-key-file] -p [your_secret_key]")
			fmt.Println()
			return
		}
	}

	s, t, m, p, a := loadArgs()

	var err error
	switch {
	case a == 1:
		err = utils.Gzip(s, t)

	case a == 2:
		err = utils.UnGzip(s, t)

	case a == 3:
		err = utils.Tar(s, t)

	case a == 4:
		err = utils.UnTar(s, t)

	case a == 5:
		err = utils.EncryptFile(s, getBytesFromMnemonic(m, p))

	case a == 6:
		err = utils.DecryptFile(s, getBytesFromMnemonic(m, p))
	}

	if err == nil {
		return
	}

	fmt.Printf("Error: %v", err)
	os.Exit(-1)
}

const nothing = "__SmIH9nWNWp__"

func loadArgs() (source string, target string, mnemonic string, secret string, action int) {
	s := flag.String("s", "", "source file/dir.")
	t := flag.String("d", nothing, "destination file/dir.")
	m := flag.String("m", nothing, "required mnemonic when encrypting/decrypting. Value should be a file name found under ~/.mnemonic that contains the mnemonic.")
	p := flag.String("p", nothing, "secret to use with the mnemonic.")
	flag.Func("a", "gz | -gz | tar | -tar | enc | -enc", func(s string) error {
		a := strings.ToLower(s)

		if action != 0 {
			return errors.New("action already set")
		}

		switch a {
		case "gz":
			action = 1
		case "-gz":
			action = 2
		case "tar":
			action = 3
		case "-tar":
			action = 4
		case "enc":
			action = 5
		case "-enc":
			action = 6
		default:
			return errors.New("invalid action")
		}

		return nil
	})

	flag.Parse()

	if action == 0 {
		*t = ""
		flag.Usage()
		fmt.Println("No action specified.")
		os.Exit(-1)
	}

	if *s == "" {
		flag.Usage()
		fmt.Println("No source file specified.")
		os.Exit(-1)
	}

	if action >= 5 && *m == nothing {
		flag.Usage()
		fmt.Println("Mnemonic file not specified.")
		os.Exit(-1)
	}

	if action <= 4 && *t == nothing {
		flag.Usage()
		fmt.Println("No destination file specified.")
		os.Exit(-1)
	}

	if action > 4 && *p == nothing {
		*p = ""
	}

	return *s, *t, *m, *p, action
}

func getBytesFromMnemonic(mnemonic, secret string) []byte {
	path := utils.ExpandHomeDir(filepath.Join("~/.mnemonic", strings.TrimSpace(mnemonic)))
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error retrieving mnemonic file '%s': %v\n", path, err)
		os.Exit(-1)
	}

	return utils.GetKeyForMnemonic(string(data), secret)
}

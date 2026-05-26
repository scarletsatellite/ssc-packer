package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"ssc-tool/archive"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func main() {
	if len(os.Args) < 3 {
		showHelp()
		return
	}

	args := os.Args[1:]
	withPassword := false
	safeMode := false
	var cleanArgs []string

	for _, arg := range args {
		switch arg {
		case "--passwd":
			withPassword = true
		case "--safe":
			safeMode = true
		default:
			cleanArgs = append(cleanArgs, arg)
		}
	}

	if len(cleanArgs) < 2 {
		showHelp()
		return
	}

	command := cleanArgs[0]

	switch command {
	case "--pack":
		targets := cleanArgs[1:]
		var password []byte
		var err error

		if withPassword {
			password, err = askPassword("Arşiv için şifre belirleyin: ")
			if err != nil {
				fmt.Printf("Şifre alınamadı: %v\n", err)
				return
			}
		} else {
			password = []byte{}
		}

		baseName := strings.TrimSuffix(filepath.Base(targets[0]), filepath.Ext(targets[0]))
		archiveName := baseName + ".ssc"

		if err := archive.Pack(targets, archiveName, password, safeMode); err != nil {
			fmt.Printf("Arşivleme Hatası: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Arşivleme tamamlandı.")

	case "--extract":
		archivePath := cleanArgs[1]

		err := archive.Extract(archivePath, []byte{})

		if errors.Is(err, archive.ErrPasswordRequired) {
			password, askErr := askPassword("Bu arşiv şifreli. Lütfen şifreyi girin: ")
			if askErr != nil {
				fmt.Printf("Şifre alınamadı: %v\n", askErr)
				return
			}

			if err = archive.Extract(archivePath, password); err != nil {
				fmt.Printf("\nÇıkarma Hatası: %v\n", err)
				os.Exit(1)
			}
		} else if err != nil {
			fmt.Printf("\nÇıkarma Hatası: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Çıkarma işlemi tamamlandı.")

	case "--help", "--yardim":
		showHelp()

	default:
		fmt.Printf("Bilinmeyen komut: %s\n", command)
		showHelp()
	}
}

func askPassword(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	pass, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return pass, err
}

func showHelp() {
	fmt.Println("Kullanım:")
	fmt.Println("  Şifresiz Arşivleme:           sscpacker.exe --pack [dosya/klasör]")
	fmt.Println("  Şifreli Arşivleme:            sscpacker.exe --passwd --pack [dosya/klasör]")
	fmt.Println("  Güvenli Arşivleme:            sscpacker.exe --safe --pack [dosya/klasör]")
	fmt.Println("  Güvenli + Şifreli Arşivleme:  sscpacker.exe --safe --passwd --pack [dosya/klasör]")
	fmt.Println("  Arşivden Çıkarma:             sscpacker.exe --extract [arsiv_adi.ssc]")
}

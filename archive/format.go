package archive

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"ssc-tool/compression"
	"ssc-tool/crypto"
	"strings"
)

const (
	MagicString = "SSC!"
	SaltSize    = 16
	NonceSize   = 12
)

var ErrPasswordRequired = errors.New("Sifre Gerekli")

type FileMeta struct {
	Name         string `json:"name"`
	OriginalSize int64  `json:"orig_size"`
	PackedSize   int64  `json:"pack_size"`
	Offset       int64  `json:"offset"`
}

// safe mod açıkken çalışan benzersiz arşiv ismi üretici
func getUniqueArchiveName(basePath string) string {
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return basePath
	}

	ext := filepath.Ext(basePath)
	sub := strings.TrimSuffix(basePath, ext)

	counter := 1
	for {
		candidate := fmt.Sprintf("%s-%02d%s", sub, counter, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
		counter++
	}
}

// --extract için otomatik çalışan klasör ismi üretici
func getUniqueDirName(baseDir string) string {
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		return baseDir
	}
	counter := 1
	for {
		candidate := fmt.Sprintf("%s-%02d", baseDir, counter)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
		counter++
	}
}

// --pack fonksiyonuna safemode parametresi eklendi
func Pack(targets []string, archiveName string, password []byte, safeMode bool) error {
	var metadata []FileMeta
	var bodyBuffer bytes.Buffer
	var currentOffset int64 = 0

	for _, target := range targets {
		info, err := os.Stat(target)
		if err != nil {
			return fmt.Errorf("hedef bulunamadı: %s, hata: %v", target, err)
		}

		if info.IsDir() {
			baseDir := filepath.Dir(target)
			err := filepath.Walk(target, func(path string, fileInfo os.FileInfo, err error) error {
				if err != nil || fileInfo.IsDir() {
					return err
				}
				relPath, _ := filepath.Rel(baseDir, path)
				return processFile(path, relPath, fileInfo.Size(), &metadata, &bodyBuffer, &currentOffset)
			})
			if err != nil {
				return err
			}
		} else {
			internalPath := filepath.Clean(target)
			if filepath.IsAbs(internalPath) {
				internalPath = filepath.Base(target)
			}

			if err := processFile(target, internalPath, info.Size(), &metadata, &bodyBuffer, &currentOffset); err != nil {
				return err
			}
		}
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	salt := make([]byte, SaltSize)
	nonce := make([]byte, NonceSize)

	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("salt uretilemedi: %v", err)
	}
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("nonce uretilemedi: %v", err)
	}

	var plainPayload bytes.Buffer
	headerSize := int32(len(metadataJSON))
	binary.Write(&plainPayload, binary.LittleEndian, headerSize)
	plainPayload.Write(metadataJSON)
	plainPayload.Write(bodyBuffer.Bytes())

	encryptedData, err := crypto.EncryptGCM(plainPayload.Bytes(), password, salt, nonce)
	if err != nil {
		return err
	}

	// hibrit mimari, safe mod aktifse benzersiz isim türet, değilse direkt üstüne yaz
	finalArchiveName := archiveName
	if safeMode {
		finalArchiveName = getUniqueArchiveName(archiveName)
	}

	outFile, err := os.Create(finalArchiveName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	outFile.WriteString(MagicString)
	outFile.Write(salt)
	outFile.Write(nonce)
	outFile.Write(encryptedData)

	fmt.Printf("\nBaşarıyla arşivlendi -> %s\n", finalArchiveName)
	return nil
}

func processFile(path, internalPath string, size int64, metadata *[]FileMeta, bodyBuf *bytes.Buffer, offset *int64) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	compressed, err := compression.CompressXZ(data, filepath.Base(path), size)
	if err != nil {
		return err
	}

	packedSize := int64(len(compressed))
	*metadata = append(*metadata, FileMeta{
		Name:         filepath.ToSlash(internalPath),
		OriginalSize: size,
		PackedSize:   packedSize,
		Offset:       *offset,
	})

	bodyBuf.Write(compressed)
	*offset += packedSize
	return nil
}

func Extract(archiveFile string, password []byte) error {
	arch, err := os.Open(archiveFile)
	if err != nil {
		return err
	}
	defer arch.Close()

	magicBuf := make([]byte, 4)
	if _, err := io.ReadFull(arch, magicBuf); err != nil {
		return err
	}
	if string(magicBuf) != MagicString {
		return fmt.Errorf("geçersiz dosya formatı")
	}

	salt := make([]byte, SaltSize)
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(arch, salt); err != nil {
		return err
	}
	if _, err := io.ReadFull(arch, nonce); err != nil {
		return err
	}

	encryptedData, err := io.ReadAll(arch)
	if err != nil {
		return err
	}

	plainData, err := crypto.DecryptGCM(encryptedData, password, salt, nonce)
	if err != nil {
		if len(password) == 0 {
			return ErrPasswordRequired
		}
		return fmt.Errorf("şifre çözme başarısız! Yanlış şifre")
	}

	plainReader := bytes.NewReader(plainData)

	var headerSize int32
	binary.Read(plainReader, binary.LittleEndian, &headerSize)

	headerBuf := make([]byte, headerSize)
	if _, err := io.ReadFull(plainReader, headerBuf); err != nil {
		return err
	}

	var metadata []FileMeta
	if err := json.Unmarshal(headerBuf, &metadata); err != nil {
		return err
	}

	dataStartOffset, _ := plainReader.Seek(0, io.SeekCurrent)

	baseTargetDir := strings.TrimSuffix(filepath.Base(archiveFile), filepath.Ext(archiveFile)) + "_extracted"
	finalTargetDir := getUniqueDirName(baseTargetDir)

	fmt.Printf("Dosyalar şuraya çıkarılıyor: %s/\n\n", finalTargetDir)

	for _, meta := range metadata {
		_, err := plainReader.Seek(dataStartOffset+meta.Offset, io.SeekStart)
		if err != nil {
			return err
		}

		compressedBuf := make([]byte, meta.PackedSize)
		if _, err := io.ReadFull(plainReader, compressedBuf); err != nil {
			return err
		}

		outPath := filepath.Join(finalTargetDir, filepath.FromSlash(meta.Name))
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return err
		}

		err = compression.DecompressXZ(compressedBuf, filepath.Base(meta.Name), meta.OriginalSize, outFile)
		outFile.Close()
		if err != nil {
			return err
		}
	}

	fmt.Println("\nİşlem başarıyla tamamlandı.")
	return nil
}

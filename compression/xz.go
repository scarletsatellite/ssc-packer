package compression

import (
	"bytes"
	"io"
	"ssc-tool/ui"

	"github.com/ulikunitz/xz"
)

func CompressXZ(data []byte, fileName string, origSize int64) ([]byte, error) {
	var buf bytes.Buffer

	xzWriter, err := xz.NewWriter(&buf)
	if err != nil {
		return nil, err
	}

	pw := &ui.ProgressWriter{
		W:        xzWriter,
		Total:    origSize,
		FileName: fileName,
	}

	if origSize == 0 {
		pw.Written = 0
		pw.Total = 1
		if _, err := pw.Write([]byte{}); err != nil {
			return nil, err
		}
	} else {
		if _, err := pw.Write(data); err != nil {
			return nil, err
		}
	}

	if err := xzWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecompressXZ(compressedData []byte, fileName string, origSize int64, out io.Writer) error {
	xzReader, err := xz.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return err
	}

	pw := &ui.ProgressWriter{
		W:        out,
		Total:    origSize,
		FileName: fileName,
	}

	if origSize > 0 {
		if _, err := io.Copy(pw, xzReader); err != nil {
			return err
		}
	}

	return nil
}

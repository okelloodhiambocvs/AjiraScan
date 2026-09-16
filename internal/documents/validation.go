package documents

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const MaxExtractedTextBytes = 2 << 20

type FileType string

const (
	PDF  FileType = "application/pdf"
	DOCX FileType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	TXT  FileType = "text/plain"
)

type Upload struct {
	Filename  string
	MediaType string
	Size      int64
	Head      []byte
}

func ValidateUpload(upload Upload, maxBytes int64) (FileType, error) {
	if upload.Size <= 0 || upload.Size > maxBytes {
		return "", errors.New("invalid upload size")
	}
	extension := strings.ToLower(filepath.Ext(upload.Filename))
	mediaType, _, err := mime.ParseMediaType(upload.MediaType)
	if err != nil {
		return "", errors.New("invalid content type")
	}
	switch extension {
	case ".pdf":
		if mediaType != string(PDF) || !bytes.HasPrefix(upload.Head, []byte("%PDF-")) {
			return "", errors.New("invalid PDF signature or content type")
		}
		return PDF, nil
	case ".docx":
		if mediaType != string(DOCX) || !bytes.HasPrefix(upload.Head, []byte("PK")) {
			return "", errors.New("invalid DOCX signature or content type")
		}
		return DOCX, nil
	case ".txt":
		if mediaType != string(TXT) {
			return "", errors.New("invalid text content type")
		}
		if !utf8Like(upload.Head) {
			return "", errors.New("text upload must be UTF-8")
		}
		return TXT, nil
	default:
		return "", errors.New("unsupported file type")
	}
}

func NewObjectKey(ownerID string, fileType FileType) string {
	extension := map[FileType]string{PDF: ".pdf", DOCX: ".docx", TXT: ".txt"}[fileType]
	return "cv/" + ownerID + "/" + uuid.NewString() + extension
}

func ReadBounded(reader io.Reader, max int64) ([]byte, error) {
	limited := io.LimitReader(reader, max+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > max {
		return nil, errors.New("content exceeds configured limit")
	}
	return data, nil
}

func utf8Like(data []byte) bool {
	for _, b := range data {
		if b == 0 {
			return false
		}
	}
	return true
}

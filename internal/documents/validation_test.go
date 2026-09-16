package documents

import "testing"

func TestValidateUploadRejectsSpoofedPDF(t *testing.T) {
	_, err := ValidateUpload(Upload{Filename: "cv.pdf", MediaType: "application/pdf", Size: 10, Head: []byte("not a pdf")}, 100)
	if err == nil {
		t.Fatal("expected spoofed PDF rejection")
	}
}

func TestValidateUploadAcceptsPDF(t *testing.T) {
	got, err := ValidateUpload(Upload{Filename: "cv.pdf", MediaType: "application/pdf", Size: 10, Head: []byte("%PDF-1.7")}, 100)
	if err != nil || got != PDF {
		t.Fatalf("expected PDF, got %s, %v", got, err)
	}
}

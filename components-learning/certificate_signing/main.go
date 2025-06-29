package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/skip2/go-qrcode"
)

func main() {
	// Generate QR code as bytes
	qrBytes, err := qrcode.Encode(
		"https://verify.example.com/certificate/12345",
		qrcode.Medium,
		256,
	)
	if err != nil {
		log.Fatal("Failed to generate QR code:", err)
	}

	// Convert to base64
	qrBase64 := base64.StdEncoding.EncodeToString(qrBytes)
	qrDataURL := "data:image/png;base64," + qrBase64

	// Certificate data
	data := struct {
		Name        string
		Course      string
		Date        string
		QRImagePath string
	}{
		Name:        "John Doe",
		Course:      "Advanced Go Programming",
		Date:        "June 29, 2025",
		QRImagePath: qrDataURL, // Use data URL
	}

	// 2. Read and render HTML template
	tmpl, err := template.ParseFiles("./cert.html")
	if err != nil {
		log.Fatal("Failed to parse template:", err)
	}

	var htmlBuffer bytes.Buffer
	if err := tmpl.Execute(&htmlBuffer, data); err != nil {
		log.Fatal("Failed to render template:", err)
	}

	// 3. Generate PDF
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal("Failed to init PDF generator:", err)
	}

	page := wkhtmltopdf.NewPageReader(bytes.NewReader(htmlBuffer.Bytes()))
	pdfg.AddPage(page)

	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Dpi.Set(300)

	if err := pdfg.Create(); err != nil {
		log.Fatal("Failed to generate PDF:", err)
	}

	if err := pdfg.WriteFile("certificate.pdf"); err != nil {
		log.Fatal("Failed to save PDF:", err)
	}

	fmt.Println("Certificate successfully generated: certificate.pdf")
}

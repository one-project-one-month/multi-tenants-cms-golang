package service

import (
	"encoding/json"
	"fmt"
	"github.com/SwanHtetAungPhyo/certificate_signing/types"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
	"os"
)

// Main function to generate the PDF
func GeneratePDF(cert *types.CompletionCertificate, filename string) error {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	drawMainBorder(pdf)
	drawHeader(pdf)
	drawTitle(pdf)
	drawRecipientSection(pdf, cert)
	drawCourseSection(pdf, cert)
	drawFooter(pdf, cert)
	embedCertificateData(pdf, cert)

	return pdf.OutputFileAndClose(filename)
}

// ---------- UI Drawing Helpers ----------

func drawMainBorder(pdf *gofpdf.Fpdf) {
	pdf.SetDrawColor(66, 103, 178)
	pdf.SetLineWidth(3)
	pdf.Rect(15, 15, 267, 180, "D")

	pdf.SetDrawColor(220, 230, 240)
	pdf.SetLineWidth(1)
	pdf.Rect(25, 25, 247, 160, "D")

	drawCornerFlourishes(pdf)
}

func drawCornerFlourishes(pdf *gofpdf.Fpdf) {
	pdf.SetDrawColor(66, 103, 178)
	pdf.SetLineWidth(2)

	// top-left
	pdf.Line(25, 35, 45, 35)
	pdf.Line(35, 25, 35, 45)

	// top-right
	pdf.Line(252, 35, 272, 35)
	pdf.Line(262, 25, 262, 45)

	// bottom-left
	pdf.Line(25, 175, 45, 175)
	pdf.Line(35, 165, 35, 185)

	// bottom-right
	pdf.Line(252, 175, 272, 175)
	pdf.Line(262, 165, 262, 185)
}

func drawHeader(pdf *gofpdf.Fpdf) {
	pdf.SetFillColor(66, 103, 178)
	pdf.Rect(40, 35, 60, 20, "F")

	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(45, 42)
	pdf.Cell(50, 8, "TECH ACADEMY")

	pdf.SetDrawColor(66, 103, 178)
	pdf.SetLineWidth(2)
	pdf.Line(120, 45, 270, 45)
}

func drawTitle(pdf *gofpdf.Fpdf) {
	pdf.SetFont("Arial", "B", 32)
	pdf.SetTextColor(45, 45, 45)
	pdf.SetY(65)
	pdf.CellFormat(0, 15, "CERTIFICATE", "", 1, "C", false, 0, "")
}

// ---------- Certificate Sections ----------

func drawRecipientSection(pdf *gofpdf.Fpdf, cert *types.CompletionCertificate) {
	pdf.SetFont("Arial", "", 14)
	pdf.SetY(90)
	pdf.CellFormat(0, 10, "This certifies that", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 22)
	pdf.CellFormat(0, 14, cert.RecipientName, "", 1, "C", false, 0, "")
}

func drawCourseSection(pdf *gofpdf.Fpdf, cert *types.CompletionCertificate) {
	pdf.SetFont("Arial", "", 14)
	pdf.SetY(115)
	pdf.CellFormat(0, 10, fmt.Sprintf("has successfully completed the course"), "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(0, 12, cert.CourseName, "", 1, "C", false, 0, "")

	if cert.Instructor != "" {
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(0, 10, fmt.Sprintf("Instructor: %s", cert.Instructor), "", 1, "C", false, 0, "")
	}
}

// ---------- Footer and QR ----------

func drawFooter(pdf *gofpdf.Fpdf, cert *types.CompletionCertificate) {
	pdf.SetFont("Arial", "", 12)
	pdf.SetY(-50)
	footerText := fmt.Sprintf("Issued by %s on %s", cert.Institution, cert.IssueDate.Format("Jan 2, 2006"))
	pdf.CellFormat(0, 10, footerText, "", 1, "C", false, 0, "")

	if cert.Grade != "" {
		pdf.CellFormat(0, 8, "Grade: "+cert.Grade, "", 1, "C", false, 0, "")
	}

	pdf.CellFormat(0, 8, fmt.Sprintf("Credits: %.2f", cert.Credits), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 8, fmt.Sprintf("Certificate ID: %s", cert.CertificateID), "", 1, "C", false, 0, "")
}

// ---------- Embed Metadata & QR ----------

func embedCertificateData(pdf *gofpdf.Fpdf, cert *types.CompletionCertificate) {
	jsonData, err := json.Marshal(cert)
	if err != nil {
		pdf.SetY(-20)
		pdf.CellFormat(0, 10, "Error embedding QR", "", 1, "C", false, 0, "")
		return
	}

	qrBytes, err := qrcode.Encode(string(jsonData), qrcode.Medium, 128)
	if err != nil {
		pdf.SetY(-20)
		pdf.CellFormat(0, 10, "Error generating QR", "", 1, "C", false, 0, "")
		return
	}

	qrPath := "tmp_qr.png"
	err = os.WriteFile(qrPath, qrBytes, 0644)
	if err != nil {
		pdf.SetY(-20)
		pdf.CellFormat(0, 10, "Error writing QR file", "", 1, "C", false, 0, "")
		return
	}
	defer os.Remove(qrPath)

	pdf.Image(qrPath, 235, 35, 40, 0, false, "", 0, "")
}

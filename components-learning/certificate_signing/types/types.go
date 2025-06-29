package types

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
	"strings"
	"time"
)

type CompletionCertificate struct {
	// Certificate identification
	CertificateID string `json:"certificate_id"`
	SerialNumber  string `json:"serial_number"`

	// Recipient information
	RecipientName  string `json:"recipient_name"`
	RecipientEmail string `json:"recipient_email,omitempty"`
	RecipientID    string `json:"recipient_id,omitempty"`

	// Course/Program details
	CourseName  string `json:"course_name"`
	CourseCode  string `json:"course_code,omitempty"`
	Institution string `json:"institution"`
	Instructor  string `json:"instructor,omitempty"`

	// Completion details
	CompletionDate time.Time  `json:"completion_date"`
	IssueDate      time.Time  `json:"issue_date"`
	ExpiryDate     *time.Time `json:"expiry_date,omitempty"`
	Grade          string     `json:"grade,omitempty"`
	Credits        float64    `json:"credits,omitempty"`

	// Verification data
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
	Hash      string `json:"hash"`

	// Additional metadata
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *CompletionCertificate) GetSignableData() []byte {
	signable := struct {
		CertificateID  string     `json:"certificate_id"`
		SerialNumber   string     `json:"serial_number"`
		RecipientName  string     `json:"recipient_name"`
		RecipientEmail string     `json:"recipient_email,omitempty"`
		RecipientID    string     `json:"recipient_id,omitempty"`
		CourseName     string     `json:"course_name"`
		CourseCode     string     `json:"course_code,omitempty"`
		Institution    string     `json:"institution"`
		Instructor     string     `json:"instructor,omitempty"`
		CompletionDate time.Time  `json:"completion_date"`
		IssueDate      time.Time  `json:"issue_date"`
		ExpiryDate     *time.Time `json:"expiry_date,omitempty"`
		Grade          string     `json:"grade,omitempty"`
		Credits        float64    `json:"credits,omitempty"`
		PublicKey      string     `json:"public_key"`
		Version        string     `json:"version"`
		CreatedAt      time.Time  `json:"created_at"`
	}{
		CertificateID:  c.CertificateID,
		SerialNumber:   c.SerialNumber,
		RecipientName:  c.RecipientName,
		RecipientEmail: c.RecipientEmail,
		RecipientID:    c.RecipientID,
		CourseName:     c.CourseName,
		CourseCode:     c.CourseCode,
		Institution:    c.Institution,
		Instructor:     c.Instructor,
		CompletionDate: c.CompletionDate,
		IssueDate:      c.IssueDate,
		ExpiryDate:     c.ExpiryDate,
		Grade:          c.Grade,
		Credits:        c.Credits,
		PublicKey:      c.PublicKey,
		Version:        c.Version,
		CreatedAt:      c.CreatedAt,
	}

	data, _ := json.Marshal(signable)
	return data
}

func (c *CompletionCertificate) Sgin(privateKey *ecdsa.PrivateKey) error {
	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	c.PublicKey = fmt.Sprintf("0x%x", publicKeyBytes)

	data := c.GetSignableData()
	hash := crypto.Keccak256Hash(data)
	c.Hash = hash.Hex()

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		color.Red(err.Error())
		return err
	}
	c.Signature = fmt.Sprintf("0x%x", signature)
	return nil

}
func (c *CompletionCertificate) Verify() (bool, error) {
	publicKeyBytes, err := hex.DecodeString(strings.TrimPrefix(c.PublicKey, "0x"))
	if err != nil {
		color.Red(err.Error())
		return false, err
	}
	_, err = crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		color.Red(err.Error())
		return false, err
	}

	data := c.GetSignableData()
	hash := crypto.Keccak256Hash(data)
	if c.Hash != hash.Hex() {
		return false, fmt.Errorf("hash mismatch")
	}
	signatureByte, err := hex.DecodeString(strings.TrimPrefix(c.Signature, "0x"))
	if err != nil {
		color.Red(err.Error())
		return false, err
	}
	return crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureByte[:64]), nil
}

func (c *CompletionCertificate) GetSignature() string {
	return c.Signature
}

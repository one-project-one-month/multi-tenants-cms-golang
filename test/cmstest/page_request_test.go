package cmstest

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const pageRequestBaseURL = "http://localhost:8081/cms/page-request"

var createdPageRequestID uuid.UUID

type pageRequestResponse struct {
	Data struct {
		ID uuid.UUID `json:"id"`
	} `json:"data"`
}

// 1. Create a page request with logo file
func TestCreatePageRequest(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("ownerId", "11111111-1111-1111-1111-111111111111")
	_ = writer.WriteField("requestType", "LANDING_PAGE")
	_ = writer.WriteField("title", "Test Landing Page")
	_ = writer.WriteField("description", "Landing page for test")
	logoPath := "/Users/swanhtet/Desktop/multi-tenants-cms-golang/img_1.png"
	file, err := os.Open(logoPath)
	assert.NoError(t, err)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
	}(file)

	part, err := writer.CreateFormFile("logo", filepath.Base(logoPath))
	assert.NoError(t, err)

	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	err = writer.Close()
	if err != nil {
		t.Fatal(err.Error())
		return
	}

	// Send request
	req, err := http.NewRequest("POST", pageRequestBaseURL+"/", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
	}(resp.Body)

	jsonBody, _ := io.ReadAll(resp.Body)
	t.Log(string(jsonBody))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var pr pageRequestResponse
	err = json.NewDecoder(resp.Body).Decode(&pr)
	assert.NoError(t, err)
	createdPageRequestID = pr.Data.ID

	t.Logf("✅ Created Page Request ID: %s", createdPageRequestID)
}

// 2. Get all page requests with pagination
func TestGetAllPageRequests(t *testing.T) {
	resp, err := http.Get(pageRequestBaseURL + "?page=1&limit=5")
	assert.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err.Error())
		}
	}(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(body), "page_request"), "Unexpected body: %s", body)

	t.Logf("✅ Fetched Page Requests")
}

// 3. Change status of created page request
func TestChangePageRequestStatus(t *testing.T) {
	if createdPageRequestID == uuid.Nil {
		t.Fatal("❌ No page request created to change status")
	}

	type ChangeStatusPageRequest struct {
		RequestID string `json:"requestId"`
		Status    string `json:"status"`
	}
	payload := ChangeStatusPageRequest{
		RequestID: createdPageRequestID.String(),
		Status:    "APPROVED",
	}

	body, _ := json.Marshal(payload)
	resp, err := http.Post(pageRequestBaseURL+"/status", "application/json", bytes.NewReader(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	t.Logf("✅ Changed status of Page Request %s to APPROVED", createdPageRequestID)
}

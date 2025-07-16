package cmstest

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

const (
	authBaseURL        = "http://localhost:8081/cms/auth"
	ownerBaseURL       = "http://localhost:8081/cms/owners"
	pageRequestBaseURL = "http://localhost:8081/cms/page-request"
)

var (
	testEmail    = "swanhtet102002@gmail.com"
	testPassword = "StrongP@ssword123"

	accessToken          string
	refreshToken         string
	userID               uuid.UUID
	createdOwnerID       string
	createdPageRequestID uuid.UUID
)

func prettyPrintJSON(data []byte) string {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return string(data)
	}
	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return string(data)
	}
	return string(pretty)
}

func TestCompleteFlow(t *testing.T) {
	t.Run("1_Register", RegisterTest)
	t.Run("2_Login", LoginTest)
	t.Run("3_CreateOwner", CreateOwnerTest)
	t.Run("4_GetAllOwners", GetAllOwnersTest)
	t.Run("5_GetOwnerByID", GetOwnerByIDTest)
	t.Run("6_UpdateOwner", UpdateOwnerTest)
	t.Run("7_CreatePageRequest", CreatePageRequestTest)
	t.Run("8_GetAllPageRequests", GetAllPageRequestsTest)
	//t.Run("9_ChangePageRequestStatus", ChangePageRequestStatusTest)
	t.Run("10_DeleteOwner", DeleteOwnerTest)
	t.Run("11_Logout", LogoutTest)
}

func RegisterTest(t *testing.T) {
	body := RegisterRequest{
		Name:     "API Tester",
		Email:    testEmail,
		Password: testPassword,
	}

	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post(authBaseURL+"/register", "application/json", bytes.NewReader(jsonBody))
	assert.NoError(t, err)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(resp.Body)

	bodyBytes, _ := io.ReadAll(resp.Body)
	t.Log("Register Response:", prettyPrintJSON(bodyBytes))

	if resp.StatusCode == http.StatusCreated {
		t.Log("✅ Registration successful")

		var regResp struct {
			Data struct {
				User UserResponse `json:"user"`
			} `json:"data"`
		}

		err = json.Unmarshal(bodyBytes, &regResp)
		if err == nil {
			userID = regResp.Data.User.ID
			t.Logf("✅ User ID extracted: %s", userID)
		}
	} else if resp.StatusCode == http.StatusConflict {
		t.Log("⚠️ User already exists - attempting to get user ID from login")
	} else {
		t.Errorf("Unexpected status: %d", resp.StatusCode)
	}
}

func LoginTest(t *testing.T) {
	body := LoginRequest{
		Email:    testEmail,
		Password: testPassword,
	}

	jsonBody, _ := json.Marshal(body)
	resp, err := http.Post(authBaseURL+"/login", "application/json", bytes.NewReader(jsonBody))
	assert.NoError(t, err)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	bodyBytes, _ := io.ReadAll(resp.Body)
	t.Log("Login Response:", prettyPrintJSON(bodyBytes))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp struct {
		Data AuthResponse `json:"data"`
	}

	err = json.Unmarshal(bodyBytes, &authResp)
	assert.NoError(t, err)

	accessToken = authResp.Data.AccessToken
	refreshToken = authResp.Data.RefreshToken

	if userID == uuid.Nil {
		userID = authResp.Data.User.ID
	}

	t.Logf("✅ Logged in. User ID: %s, AccessToken: %.20s..., RefreshToken: %.20s...",
		userID, accessToken, refreshToken)

	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	assert.NotEqual(t, uuid.Nil, userID)
}

func CreateOwnerTest(t *testing.T) {
	req := OwnerCreateRequest{
		Name:      "Test Owner",
		Email:     "testowner@example.com",
		NameSpace: "test-space-123",
		Password:  "SecurePass123",
	}

	jsonBody, _ := json.Marshal(req)
	resp, err := http.Post(ownerBaseURL+"/", "application/json", bytes.NewReader(jsonBody))

	assert.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err.Error())
		}
	}(resp.Body)

	bodyBytes, _ := io.ReadAll(resp.Body)
	t.Log("Create Owner Response:", prettyPrintJSON(bodyBytes))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result struct {
		Data OwnerResponse `json:"data"`
	}
	_ = json.Unmarshal(bodyBytes, &result)
	createdOwnerID = result.Data.ID.String()

	t.Logf("✅ Created Owner: ID = %s", createdOwnerID)
}

func GetAllOwnersTest(t *testing.T) {
	resp, err := http.Get(ownerBaseURL + "/")
	assert.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(resp.Body)

	bodyBytes, _ := io.ReadAll(resp.Body)
	t.Log("Get All Owners Response:", prettyPrintJSON(bodyBytes))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result struct {
		Data []OwnerResponse `json:"data"`
	}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.True(t, len(result.Data) > 0)

	t.Logf("✅ Fetched %d owner(s)", len(result.Data))
}

func GetOwnerByIDTest(t *testing.T) {
	if createdOwnerID == "" {
		t.Fatal("Owner not created")
	}

	resp, err := http.Get(ownerBaseURL + "/" + createdOwnerID)
	assert.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(resp.Body)

	bodyBytes, _ := io.ReadAll(resp.Body)
	t.Log("Get Owner By ID Response:", prettyPrintJSON(bodyBytes))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result struct {
		Data OwnerResponse `json:"data"`
	}
	err = json.Unmarshal(bodyBytes, &result)
	assert.NoError(t, err)
	assert.Equal(t, createdOwnerID, result.Data.ID.String())

	t.Logf("✅ Fetched Owner by ID: %s", createdOwnerID)
}

func UpdateOwnerTest(t *testing.T) {
	if createdOwnerID == "" {
		t.Fatal("Owner not created")
	}

	req := OwnerUpdateRequest{
		Name:      "Updated Owner",
		NameSpace: "updatedspace",
	}
	jsonBody, _ := json.Marshal(req)

	request, err := http.NewRequest(http.MethodPut, ownerBaseURL+"/"+createdOwnerID, bytes.NewReader(jsonBody))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	assert.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(resp.Body)

	bodyBytes, _ := io.ReadAll(resp.Body)
	t.Log("Update Owner Response:", prettyPrintJSON(bodyBytes))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result struct {
		Data OwnerResponse `json:"data"`
	}
	_ = json.Unmarshal(bodyBytes, &result)

	assert.Equal(t, "Updated Owner", result.Data.Name)
	t.Logf("✅ Updated Owner: %s", result.Data.ID)
}

func CreatePageRequestTest(t *testing.T) {
	if createdOwnerID == "" {
		t.Fatal("Owner not created - cannot create page request")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("ownerId", createdOwnerID)
	_ = writer.WriteField("requestType", "LANDING_PAGE")
	_ = writer.WriteField("title", "Test Landing Page")
	_ = writer.WriteField("description", "Landing page for test")

	logoPath := "/Users/swanhtet/Desktop/multi-tenants-cms-golang/img_1.png"
	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		tmpFile, err := os.CreateTemp("", "test_logo*.png")
		if err != nil {
			t.Fatal("Failed to create temp file:", err)
		}
		defer os.Remove(tmpFile.Name())

		tmpFile.Write([]byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
		tmpFile.Close()
		logoPath = tmpFile.Name()
	}

	file, err := os.Open(logoPath)
	assert.NoError(t, err)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			t.Error(err.Error())
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

	req, err := http.NewRequest("POST", pageRequestBaseURL+"/", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err.Error())
		}
	}(resp.Body)

	bodyBytes, _ := io.ReadAll(resp.Body)
	t.Log("Create Page Request Response:", prettyPrintJSON(bodyBytes))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var pr struct {
		Data struct {
			ID uuid.UUID `json:"id"`
		} `json:"data"`
	}

	err = json.Unmarshal(bodyBytes, &pr)
	assert.NoError(t, err)
	createdPageRequestID = pr.Data.ID

	t.Logf("✅ Created Page Request ID: %s for Owner: %s", createdPageRequestID, createdOwnerID)
}

func GetAllPageRequestsTest(t *testing.T) {
	req, err := http.NewRequest("GET", pageRequestBaseURL+"?page=1&limit=5", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err.Error())
		}
	}(resp.Body)

	bodyBytes, _ := io.ReadAll(resp.Body)
	t.Log("Get All Page Requests Response:", prettyPrintJSON(bodyBytes))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	t.Logf("✅ Fetched Page Requests")
}

func ChangePageRequestStatusTest(t *testing.T) {
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
	req, err := http.NewRequest("POST", pageRequestBaseURL+"/status", bytes.NewReader(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	t.Log("Change Page Request Status Response:", prettyPrintJSON(bodyBytes))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	t.Logf("✅ Changed status of Page Request %s to APPROVED", createdPageRequestID)
}

func DeleteOwnerTest(t *testing.T) {
	if createdOwnerID == "" {
		t.Fatal("Owner not created")
	}

	req := OwnerDeleteRequest{
		IDs:         []string{createdOwnerID},
		ForceDelete: true,
	}
	jsonBody, _ := json.Marshal(req)

	request, err := http.NewRequest(http.MethodDelete, ownerBaseURL+"/", bytes.NewReader(jsonBody))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(request)
	assert.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Error(err)
		}
	}(resp.Body)

	bodyBytes, _ := io.ReadAll(resp.Body)
	t.Log("Delete Owner Response:", prettyPrintJSON(bodyBytes))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	t.Logf("✅ Deleted Owner: %s", createdOwnerID)
}

func LogoutTest(t *testing.T) {
	if accessToken == "" && refreshToken == "" {
		t.Fatal("Tokens not set. Run LoginTest first.")
	}

	payload := map[string]string{
		"refresh_token": refreshToken,
	}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", authBaseURL+"/logout", bytes.NewReader(jsonPayload))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	bodyBytes, _ := io.ReadAll(resp.Body)
	t.Log("Logout Response:", prettyPrintJSON(bodyBytes))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	accessToken = ""
	refreshToken = ""

	t.Logf("✅ Logout successful for user: %s", userID)
}

//func TestAuthFlow(t *testing.T) {
//	t.Run("Register", RegisterTest)
//	t.Run("Login", LoginTest)
//	t.Run("Logout", LogoutTest)
//}

func TestOwnerFlow(t *testing.T) {
	if accessToken == "" {
		LoginTest(t)
	}

	t.Run("CreateOwner", CreateOwnerTest)
	t.Run("GetAllOwners", GetAllOwnersTest)
	t.Run("GetOwnerByID", GetOwnerByIDTest)
	t.Run("UpdateOwner", UpdateOwnerTest)
	t.Run("DeleteOwner", DeleteOwnerTest)
}

func TestPageRequestFlow(t *testing.T) {
	if accessToken == "" {
		LoginTest(t)
	}
	if createdOwnerID == "" {
		CreateOwnerTest(t)
	}

	//t.Run("CreatePageRequest", CreatePageRequestTest)
	t.Run("GetAllPageRequests", GetAllPageRequestsTest)
	//t.Run("ChangePageRequestStatus", ChangePageRequestStatusTest)
}

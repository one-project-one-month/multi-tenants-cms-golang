package cmstest

import (
	"bytes"

	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

const ownerBaseURL = "http://localhost:8081/cms/owners"

var createdOwnerID string

func TestOwnerEndpoints(t *testing.T) {
	t.Run("CreateOwner", CreateOwnerTest)
	t.Run("GetAllOwners", GetAllOwnersTest)
	t.Run("GetOwnerByID", GetOwnerByIDTest)
	t.Run("UpdateOwner", UpdateOwnerTest)
	t.Run("DeleteOwner", DeleteOwnerTest)
}

func CreateOwnerTest(t *testing.T) {
	req := OwnerCreateRequest{
		Name:      "Test Owner",
		Email:     "testowner@example.com",
		NameSpace: "testspace",
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
	assert.Equal(t, http.StatusCreated, resp.StatusCode, string(bodyBytes))

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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result struct {
		Data []OwnerResponse `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.True(t, len(result.Data) > 0)

	t.Logf(" Fetched %d owner(s)", len(result.Data))
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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result struct {
		Data OwnerResponse `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, createdOwnerID, result.Data.ID.String())

	t.Logf("✅ Fetched Owner by ID")
}

func UpdateOwnerTest(t *testing.T) {
	if createdOwnerID == "" {
		t.Fatal(" Owner not created")
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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result struct {
		Data OwnerResponse `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&result)

	assert.Equal(t, "Updated Owner", result.Data.Name)
	t.Logf("Updated Owner: %s", result.Data.ID)
}

func DeleteOwnerTest(t *testing.T) {
	if createdOwnerID == "" {
		t.Fatal(" Owner not created")
	}

	req := OwnerDeleteRequest{
		IDs:         []string{createdOwnerID},
		ForceDelete: true,
	}
	jsonBody, _ := json.Marshal(req)

	request, err := http.NewRequest(http.MethodDelete, ownerBaseURL+"/", bytes.NewReader(jsonBody))
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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	t.Logf("✅ Deleted Owner: %s", createdOwnerID)
}

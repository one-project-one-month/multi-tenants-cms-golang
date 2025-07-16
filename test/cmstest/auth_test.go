package cmstest

import (
	"testing"
)

const baseURL = "http://localhost:8081/cms/auth"

//var testEmail = "swanhtet102002@gmail.com"
//var testPassword = "StrongP@ssword123"
////var accessToken string
//var refreshToken string

func TestRegisterLoginLogout(t *testing.T) {
	t.Run("Register", RegisterTest)
	//t.Run("Login", LoginTest)
	//t.Run("Logout", LogoutTest)
}

//func RegisterTest(t *testing.T) {
//	body := RegisterRequest{
//		Name:     "API Tester",
//		Email:    testEmail,
//		Password: testPassword,
//	}
//
//	jsonBody, _ := json.Marshal(body)
//	resp, err := http.Post(baseURL+"/register", "application/json", bytes.NewReader(jsonBody))
//	assert.NoError(t, err)
//
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Fatal(err.Error())
//		}
//	}(resp.Body)
//	bodyBytes, _ := io.ReadAll(resp.Body)
//
//	if resp.StatusCode == http.StatusCreated {
//		t.Log("✅ Registration successful")
//	} else if resp.StatusCode == http.StatusConflict {
//		t.Log("⚠️ User already exists")
//	} else {
//		t.Errorf("Unexpected status: %d - %s", resp.StatusCode, string(bodyBytes))
//	}
//}
//
//func LoginTest(t *testing.T) {
//	body := LoginRequest{
//		Email:    testEmail,
//		Password: testPassword,
//	}
//	jsonBody, _ := json.Marshal(body)
//	resp, err := http.Post(baseURL+"/login", "application/json", bytes.NewReader(jsonBody))
//	assert.NoError(t, err)
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//	}(resp.Body)
//
//	t.Log(resp.Body)
//	assert.Equal(t, http.StatusOK, resp.StatusCode)
//
//	var authResp struct {
//		Data AuthResponse `json:"data"`
//	}
//
//	err = json.NewDecoder(resp.Body).Decode(&authResp)
//	assert.NoError(t, err)
//
//	accessToken = authResp.Data.AccessToken
//	refreshToken = authResp.Data.RefreshToken
//
//	t.Logf("✅ Logged in. AccessToken: %.20s..., RefreshToken: %.20s...", accessToken, refreshToken)
//	assert.NotEmpty(t, accessToken)
//	assert.NotEmpty(t, refreshToken)
//}
//
//func LogoutTest(t *testing.T) {
//	if accessToken == "" && refreshToken == "" {
//		t.Fatal(" Tokens not set. Run LoginTest first.")
//	}
//
//	payload := map[string]string{
//		"refresh_token": refreshToken,
//	}
//	jsonPayload, _ := json.Marshal(payload)
//
//	req, err := http.NewRequest("POST", baseURL+"/logout", bytes.NewReader(jsonPayload))
//	assert.NoError(t, err)
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("Authorization", "Bearer "+accessToken)
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	assert.NoError(t, err)
//	defer func(Body io.ReadCloser) {
//		err := Body.Close()
//		if err != nil {
//			log.Fatal(err)
//		}
//	}(resp.Body)
//
//	assert.Equal(t, http.StatusOK, resp.StatusCode)
//	t.Log("✅ Logout successful")
//}

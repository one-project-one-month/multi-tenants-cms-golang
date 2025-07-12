# \MFAAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AuthMfaSetupPost**](MFAAPI.md#AuthMfaSetupPost) | **Post** /auth/mfa/setup | Setup MFA
[**AuthMfaVerifyPost**](MFAAPI.md#AuthMfaVerifyPost) | **Post** /auth/mfa/verify | Verify MFA setup



## AuthMfaSetupPost

> MFASetupResponse AuthMfaSetupPost(ctx).Execute()

Setup MFA



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.MFAAPI.AuthMfaSetupPost(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MFAAPI.AuthMfaSetupPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AuthMfaSetupPost`: MFASetupResponse
	fmt.Fprintf(os.Stdout, "Response from `MFAAPI.AuthMfaSetupPost`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiAuthMfaSetupPostRequest struct via the builder pattern


### Return type

[**MFASetupResponse**](MFASetupResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AuthMfaVerifyPost

> AuthMfaVerifyPost(ctx).MFAVerificationRequest(mFAVerificationRequest).Execute()

Verify MFA setup



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	mFAVerificationRequest := *openapiclient.NewMFAVerificationRequest(int32(123), "Code_example") // MFAVerificationRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.MFAAPI.AuthMfaVerifyPost(context.Background()).MFAVerificationRequest(mFAVerificationRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MFAAPI.AuthMfaVerifyPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthMfaVerifyPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **mFAVerificationRequest** | [**MFAVerificationRequest**](MFAVerificationRequest.md) |  | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


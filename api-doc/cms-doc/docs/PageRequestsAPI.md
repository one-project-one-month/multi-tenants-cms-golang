# \PageRequestsAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PageRequestGet**](PageRequestsAPI.md#PageRequestGet) | **Get** /page-request | Get all page requests
[**PageRequestPost**](PageRequestsAPI.md#PageRequestPost) | **Post** /page-request | Create page request



## PageRequestGet

> PageRequestGet200Response PageRequestGet(ctx).Page(page).Limit(limit).Execute()

Get all page requests



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
	page := int32(56) // int32 |  (optional) (default to 1)
	limit := int32(56) // int32 |  (optional) (default to 10)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PageRequestsAPI.PageRequestGet(context.Background()).Page(page).Limit(limit).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PageRequestsAPI.PageRequestGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PageRequestGet`: PageRequestGet200Response
	fmt.Fprintf(os.Stdout, "Response from `PageRequestsAPI.PageRequestGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPageRequestGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **page** | **int32** |  | [default to 1]
 **limit** | **int32** |  | [default to 10]

### Return type

[**PageRequestGet200Response**](PageRequestGet200Response.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PageRequestPost

> PageRequestResponse PageRequestPost(ctx).OwnerId(ownerId).RequestType(requestType).Title(title).Description(description).PageUrl(pageUrl).Logo(logo).Execute()

Create page request



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
	ownerId := "38400000-8cf0-11bd-b23e-10b96e4ef00d" // string | 
	requestType := "requestType_example" // string | 
	title := "title_example" // string | 
	description := "description_example" // string | 
	pageUrl := "pageUrl_example" // string |  (optional)
	logo := os.NewFile(1234, "some_file") // *os.File |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.PageRequestsAPI.PageRequestPost(context.Background()).OwnerId(ownerId).RequestType(requestType).Title(title).Description(description).PageUrl(pageUrl).Logo(logo).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `PageRequestsAPI.PageRequestPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `PageRequestPost`: PageRequestResponse
	fmt.Fprintf(os.Stdout, "Response from `PageRequestsAPI.PageRequestPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPageRequestPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ownerId** | **string** |  | 
 **requestType** | **string** |  | 
 **title** | **string** |  | 
 **description** | **string** |  | 
 **pageUrl** | **string** |  | 
 **logo** | ***os.File** |  | 

### Return type

[**PageRequestResponse**](PageRequestResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: multipart/form-data
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


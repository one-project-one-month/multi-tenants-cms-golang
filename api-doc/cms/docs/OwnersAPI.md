# \OwnersAPI

All URIs are relative to *http://localhost:8080*

Method | HTTP request | Description
------------- | ------------- | -------------
[**OwnersDelete**](OwnersAPI.md#OwnersDelete) | **Delete** /owners | Delete owners
[**OwnersGet**](OwnersAPI.md#OwnersGet) | **Get** /owners | Get all owners
[**OwnersIdGet**](OwnersAPI.md#OwnersIdGet) | **Get** /owners/{id} | Get owner by ID
[**OwnersIdPut**](OwnersAPI.md#OwnersIdPut) | **Put** /owners/{id} | Update owner
[**OwnersPost**](OwnersAPI.md#OwnersPost) | **Post** /owners | Create owner



## OwnersDelete

> OwnersDelete(ctx).OwnerDeleteRequest(ownerDeleteRequest).Execute()

Delete owners



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
	ownerDeleteRequest := *openapiclient.NewOwnerDeleteRequest([]string{"Ids_example"}) // OwnerDeleteRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.OwnersAPI.OwnersDelete(context.Background()).OwnerDeleteRequest(ownerDeleteRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OwnersAPI.OwnersDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiOwnersDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ownerDeleteRequest** | [**OwnerDeleteRequest**](OwnerDeleteRequest.md) |  | 

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


## OwnersGet

> []OwnerResponse OwnersGet(ctx).Execute()

Get all owners



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
	resp, r, err := apiClient.OwnersAPI.OwnersGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OwnersAPI.OwnersGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OwnersGet`: []OwnerResponse
	fmt.Fprintf(os.Stdout, "Response from `OwnersAPI.OwnersGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiOwnersGetRequest struct via the builder pattern


### Return type

[**[]OwnerResponse**](OwnerResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OwnersIdGet

> OwnerResponse OwnersIdGet(ctx, id).Execute()

Get owner by ID



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
	id := "id_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OwnersAPI.OwnersIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OwnersAPI.OwnersIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OwnersIdGet`: OwnerResponse
	fmt.Fprintf(os.Stdout, "Response from `OwnersAPI.OwnersIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiOwnersIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**OwnerResponse**](OwnerResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OwnersIdPut

> OwnerResponse OwnersIdPut(ctx, id).OwnerUpdateRequest(ownerUpdateRequest).Execute()

Update owner



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
	id := "id_example" // string | 
	ownerUpdateRequest := *openapiclient.NewOwnerUpdateRequest("Name_example", "Namespace_example") // OwnerUpdateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OwnersAPI.OwnersIdPut(context.Background(), id).OwnerUpdateRequest(ownerUpdateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OwnersAPI.OwnersIdPut``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OwnersIdPut`: OwnerResponse
	fmt.Fprintf(os.Stdout, "Response from `OwnersAPI.OwnersIdPut`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiOwnersIdPutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **ownerUpdateRequest** | [**OwnerUpdateRequest**](OwnerUpdateRequest.md) |  | 

### Return type

[**OwnerResponse**](OwnerResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## OwnersPost

> OwnerResponse OwnersPost(ctx).OwnerCreateRequest(ownerCreateRequest).Execute()

Create owner



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
	ownerCreateRequest := *openapiclient.NewOwnerCreateRequest("Name_example", "Email_example", "Namespace_example", "Password_example") // OwnerCreateRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.OwnersAPI.OwnersPost(context.Background()).OwnerCreateRequest(ownerCreateRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `OwnersAPI.OwnersPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `OwnersPost`: OwnerResponse
	fmt.Fprintf(os.Stdout, "Response from `OwnersAPI.OwnersPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiOwnersPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ownerCreateRequest** | [**OwnerCreateRequest**](OwnerCreateRequest.md) |  | 

### Return type

[**OwnerResponse**](OwnerResponse.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


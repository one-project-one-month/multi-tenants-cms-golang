# \TenantServiceAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**TenantServiceCreateTenant**](TenantServiceAPI.md#TenantServiceCreateTenant) | **Post** /lms/v1/tenants | Create Tenant



## TenantServiceCreateTenant

> TenantsCreateTenantResponse TenantServiceCreateTenant(ctx).Body(body).Execute()

Create Tenant



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
	body := *openapiclient.NewTenantsCreateTenantRequest() // TenantsCreateTenantRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.TenantServiceAPI.TenantServiceCreateTenant(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `TenantServiceAPI.TenantServiceCreateTenant``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `TenantServiceCreateTenant`: TenantsCreateTenantResponse
	fmt.Fprintf(os.Stdout, "Response from `TenantServiceAPI.TenantServiceCreateTenant`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTenantServiceCreateTenantRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**TenantsCreateTenantRequest**](TenantsCreateTenantRequest.md) |  | 

### Return type

[**TenantsCreateTenantResponse**](TenantsCreateTenantResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


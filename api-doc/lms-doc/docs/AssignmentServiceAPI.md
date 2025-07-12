# \AssignmentServiceAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AssignmentServiceCreateAssignment**](AssignmentServiceAPI.md#AssignmentServiceCreateAssignment) | **Post** /lms/v1/assignment | Create Assignment



## AssignmentServiceCreateAssignment

> AssignmentCreateAssignmentResponse AssignmentServiceCreateAssignment(ctx).Body(body).Execute()

Create Assignment



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
	body := *openapiclient.NewAssignmentCreateAssignmentRequest() // AssignmentCreateAssignmentRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.AssignmentServiceAPI.AssignmentServiceCreateAssignment(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `AssignmentServiceAPI.AssignmentServiceCreateAssignment``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `AssignmentServiceCreateAssignment`: AssignmentCreateAssignmentResponse
	fmt.Fprintf(os.Stdout, "Response from `AssignmentServiceAPI.AssignmentServiceCreateAssignment`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAssignmentServiceCreateAssignmentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AssignmentCreateAssignmentRequest**](AssignmentCreateAssignmentRequest.md) |  | 

### Return type

[**AssignmentCreateAssignmentResponse**](AssignmentCreateAssignmentResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


# \CourseServiceAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CourseServiceCreateCourse**](CourseServiceAPI.md#CourseServiceCreateCourse) | **Post** /lms/v1/course | Create Course
[**CourseServiceGetCourse**](CourseServiceAPI.md#CourseServiceGetCourse) | **Get** /lms/v1/course | Create Course



## CourseServiceCreateCourse

> CourseCreateCourseResponse CourseServiceCreateCourse(ctx).Body(body).Execute()

Create Course



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
	body := *openapiclient.NewCourseCreateCourseRequest() // CourseCreateCourseRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CourseServiceAPI.CourseServiceCreateCourse(context.Background()).Body(body).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CourseServiceAPI.CourseServiceCreateCourse``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CourseServiceCreateCourse`: CourseCreateCourseResponse
	fmt.Fprintf(os.Stdout, "Response from `CourseServiceAPI.CourseServiceCreateCourse`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCourseServiceCreateCourseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CourseCreateCourseRequest**](CourseCreateCourseRequest.md) |  | 

### Return type

[**CourseCreateCourseResponse**](CourseCreateCourseResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CourseServiceGetCourse

> CourseGetCourseResponse CourseServiceGetCourse(ctx).CourseTitle(courseTitle).Execute()

Create Course



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
	courseTitle := "courseTitle_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.CourseServiceAPI.CourseServiceGetCourse(context.Background()).CourseTitle(courseTitle).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `CourseServiceAPI.CourseServiceGetCourse``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CourseServiceGetCourse`: CourseGetCourseResponse
	fmt.Fprintf(os.Stdout, "Response from `CourseServiceAPI.CourseServiceGetCourse`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCourseServiceGetCourseRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **courseTitle** | **string** |  | 

### Return type

[**CourseGetCourseResponse**](CourseGetCourseResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


# PaginateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Page** | Pointer to **int32** |  | [optional] 
**Limit** | Pointer to **int32** |  | [optional] 

## Methods

### NewPaginateRequest

`func NewPaginateRequest() *PaginateRequest`

NewPaginateRequest instantiates a new PaginateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPaginateRequestWithDefaults

`func NewPaginateRequestWithDefaults() *PaginateRequest`

NewPaginateRequestWithDefaults instantiates a new PaginateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPage

`func (o *PaginateRequest) GetPage() int32`

GetPage returns the Page field if non-nil, zero value otherwise.

### GetPageOk

`func (o *PaginateRequest) GetPageOk() (*int32, bool)`

GetPageOk returns a tuple with the Page field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPage

`func (o *PaginateRequest) SetPage(v int32)`

SetPage sets Page field to given value.

### HasPage

`func (o *PaginateRequest) HasPage() bool`

HasPage returns a boolean if a field has been set.

### GetLimit

`func (o *PaginateRequest) GetLimit() int32`

GetLimit returns the Limit field if non-nil, zero value otherwise.

### GetLimitOk

`func (o *PaginateRequest) GetLimitOk() (*int32, bool)`

GetLimitOk returns a tuple with the Limit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLimit

`func (o *PaginateRequest) SetLimit(v int32)`

SetLimit sets Limit field to given value.

### HasLimit

`func (o *PaginateRequest) HasLimit() bool`

HasLimit returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



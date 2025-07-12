# OwnerCreateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Email** | **string** |  | 
**Namespace** | **string** |  | 
**Password** | **string** |  | 

## Methods

### NewOwnerCreateRequest

`func NewOwnerCreateRequest(name string, email string, namespace string, password string, ) *OwnerCreateRequest`

NewOwnerCreateRequest instantiates a new OwnerCreateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOwnerCreateRequestWithDefaults

`func NewOwnerCreateRequestWithDefaults() *OwnerCreateRequest`

NewOwnerCreateRequestWithDefaults instantiates a new OwnerCreateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *OwnerCreateRequest) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *OwnerCreateRequest) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *OwnerCreateRequest) SetName(v string)`

SetName sets Name field to given value.


### GetEmail

`func (o *OwnerCreateRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *OwnerCreateRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *OwnerCreateRequest) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetNamespace

`func (o *OwnerCreateRequest) GetNamespace() string`

GetNamespace returns the Namespace field if non-nil, zero value otherwise.

### GetNamespaceOk

`func (o *OwnerCreateRequest) GetNamespaceOk() (*string, bool)`

GetNamespaceOk returns a tuple with the Namespace field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNamespace

`func (o *OwnerCreateRequest) SetNamespace(v string)`

SetNamespace sets Namespace field to given value.


### GetPassword

`func (o *OwnerCreateRequest) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *OwnerCreateRequest) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *OwnerCreateRequest) SetPassword(v string)`

SetPassword sets Password field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



# OwnerDeleteRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Ids** | **[]string** |  | 
**ForceDelete** | Pointer to **bool** |  | [optional] 

## Methods

### NewOwnerDeleteRequest

`func NewOwnerDeleteRequest(ids []string, ) *OwnerDeleteRequest`

NewOwnerDeleteRequest instantiates a new OwnerDeleteRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOwnerDeleteRequestWithDefaults

`func NewOwnerDeleteRequestWithDefaults() *OwnerDeleteRequest`

NewOwnerDeleteRequestWithDefaults instantiates a new OwnerDeleteRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIds

`func (o *OwnerDeleteRequest) GetIds() []string`

GetIds returns the Ids field if non-nil, zero value otherwise.

### GetIdsOk

`func (o *OwnerDeleteRequest) GetIdsOk() (*[]string, bool)`

GetIdsOk returns a tuple with the Ids field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIds

`func (o *OwnerDeleteRequest) SetIds(v []string)`

SetIds sets Ids field to given value.


### GetForceDelete

`func (o *OwnerDeleteRequest) GetForceDelete() bool`

GetForceDelete returns the ForceDelete field if non-nil, zero value otherwise.

### GetForceDeleteOk

`func (o *OwnerDeleteRequest) GetForceDeleteOk() (*bool, bool)`

GetForceDeleteOk returns a tuple with the ForceDelete field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetForceDelete

`func (o *OwnerDeleteRequest) SetForceDelete(v bool)`

SetForceDelete sets ForceDelete field to given value.

### HasForceDelete

`func (o *OwnerDeleteRequest) HasForceDelete() bool`

HasForceDelete returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



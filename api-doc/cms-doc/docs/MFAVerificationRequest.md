# MFAVerificationRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TokenId** | **int32** | Temporary token ID from setup | 
**Code** | **string** | 6-digit verification code from authenticator app | 

## Methods

### NewMFAVerificationRequest

`func NewMFAVerificationRequest(tokenId int32, code string, ) *MFAVerificationRequest`

NewMFAVerificationRequest instantiates a new MFAVerificationRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMFAVerificationRequestWithDefaults

`func NewMFAVerificationRequestWithDefaults() *MFAVerificationRequest`

NewMFAVerificationRequestWithDefaults instantiates a new MFAVerificationRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTokenId

`func (o *MFAVerificationRequest) GetTokenId() int32`

GetTokenId returns the TokenId field if non-nil, zero value otherwise.

### GetTokenIdOk

`func (o *MFAVerificationRequest) GetTokenIdOk() (*int32, bool)`

GetTokenIdOk returns a tuple with the TokenId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenId

`func (o *MFAVerificationRequest) SetTokenId(v int32)`

SetTokenId sets TokenId field to given value.


### GetCode

`func (o *MFAVerificationRequest) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *MFAVerificationRequest) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *MFAVerificationRequest) SetCode(v string)`

SetCode sets Code field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



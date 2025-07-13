# MFALoginRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Email** | **string** |  | 
**Password** | **string** |  | 
**MfaCode** | **string** | 6-digit MFA code from authenticator app | 

## Methods

### NewMFALoginRequest

`func NewMFALoginRequest(email string, password string, mfaCode string, ) *MFALoginRequest`

NewMFALoginRequest instantiates a new MFALoginRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMFALoginRequestWithDefaults

`func NewMFALoginRequestWithDefaults() *MFALoginRequest`

NewMFALoginRequestWithDefaults instantiates a new MFALoginRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEmail

`func (o *MFALoginRequest) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *MFALoginRequest) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *MFALoginRequest) SetEmail(v string)`

SetEmail sets Email field to given value.


### GetPassword

`func (o *MFALoginRequest) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *MFALoginRequest) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *MFALoginRequest) SetPassword(v string)`

SetPassword sets Password field to given value.


### GetMfaCode

`func (o *MFALoginRequest) GetMfaCode() string`

GetMfaCode returns the MfaCode field if non-nil, zero value otherwise.

### GetMfaCodeOk

`func (o *MFALoginRequest) GetMfaCodeOk() (*string, bool)`

GetMfaCodeOk returns a tuple with the MfaCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMfaCode

`func (o *MFALoginRequest) SetMfaCode(v string)`

SetMfaCode sets MfaCode field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



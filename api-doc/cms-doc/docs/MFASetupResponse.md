# MFASetupResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Secret** | Pointer to **string** | The TOTP secret key | [optional] 
**QrCodeUrl** | Pointer to **string** | URL to generate QR code | [optional] 
**QrCodeImage** | Pointer to **string** | Base64 encoded PNG of QR code | [optional] 
**TokenId** | Pointer to **int32** | Temporary token ID for verification | [optional] 
**ManualEntry** | Pointer to **string** | Secret key for manual entry | [optional] 

## Methods

### NewMFASetupResponse

`func NewMFASetupResponse() *MFASetupResponse`

NewMFASetupResponse instantiates a new MFASetupResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMFASetupResponseWithDefaults

`func NewMFASetupResponseWithDefaults() *MFASetupResponse`

NewMFASetupResponseWithDefaults instantiates a new MFASetupResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSecret

`func (o *MFASetupResponse) GetSecret() string`

GetSecret returns the Secret field if non-nil, zero value otherwise.

### GetSecretOk

`func (o *MFASetupResponse) GetSecretOk() (*string, bool)`

GetSecretOk returns a tuple with the Secret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecret

`func (o *MFASetupResponse) SetSecret(v string)`

SetSecret sets Secret field to given value.

### HasSecret

`func (o *MFASetupResponse) HasSecret() bool`

HasSecret returns a boolean if a field has been set.

### GetQrCodeUrl

`func (o *MFASetupResponse) GetQrCodeUrl() string`

GetQrCodeUrl returns the QrCodeUrl field if non-nil, zero value otherwise.

### GetQrCodeUrlOk

`func (o *MFASetupResponse) GetQrCodeUrlOk() (*string, bool)`

GetQrCodeUrlOk returns a tuple with the QrCodeUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQrCodeUrl

`func (o *MFASetupResponse) SetQrCodeUrl(v string)`

SetQrCodeUrl sets QrCodeUrl field to given value.

### HasQrCodeUrl

`func (o *MFASetupResponse) HasQrCodeUrl() bool`

HasQrCodeUrl returns a boolean if a field has been set.

### GetQrCodeImage

`func (o *MFASetupResponse) GetQrCodeImage() string`

GetQrCodeImage returns the QrCodeImage field if non-nil, zero value otherwise.

### GetQrCodeImageOk

`func (o *MFASetupResponse) GetQrCodeImageOk() (*string, bool)`

GetQrCodeImageOk returns a tuple with the QrCodeImage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQrCodeImage

`func (o *MFASetupResponse) SetQrCodeImage(v string)`

SetQrCodeImage sets QrCodeImage field to given value.

### HasQrCodeImage

`func (o *MFASetupResponse) HasQrCodeImage() bool`

HasQrCodeImage returns a boolean if a field has been set.

### GetTokenId

`func (o *MFASetupResponse) GetTokenId() int32`

GetTokenId returns the TokenId field if non-nil, zero value otherwise.

### GetTokenIdOk

`func (o *MFASetupResponse) GetTokenIdOk() (*int32, bool)`

GetTokenIdOk returns a tuple with the TokenId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTokenId

`func (o *MFASetupResponse) SetTokenId(v int32)`

SetTokenId sets TokenId field to given value.

### HasTokenId

`func (o *MFASetupResponse) HasTokenId() bool`

HasTokenId returns a boolean if a field has been set.

### GetManualEntry

`func (o *MFASetupResponse) GetManualEntry() string`

GetManualEntry returns the ManualEntry field if non-nil, zero value otherwise.

### GetManualEntryOk

`func (o *MFASetupResponse) GetManualEntryOk() (*string, bool)`

GetManualEntryOk returns a tuple with the ManualEntry field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetManualEntry

`func (o *MFASetupResponse) SetManualEntry(v string)`

SetManualEntry sets ManualEntry field to given value.

### HasManualEntry

`func (o *MFASetupResponse) HasManualEntry() bool`

HasManualEntry returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



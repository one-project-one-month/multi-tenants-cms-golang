# PageRequestResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**OwnerId** | Pointer to **string** |  | [optional] 
**RequestType** | Pointer to **string** |  | [optional] 
**Title** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**PageUrl** | Pointer to **NullableString** |  | [optional] 
**LogoUrl** | Pointer to **NullableString** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewPageRequestResponse

`func NewPageRequestResponse() *PageRequestResponse`

NewPageRequestResponse instantiates a new PageRequestResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPageRequestResponseWithDefaults

`func NewPageRequestResponseWithDefaults() *PageRequestResponse`

NewPageRequestResponseWithDefaults instantiates a new PageRequestResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *PageRequestResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *PageRequestResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *PageRequestResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *PageRequestResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetOwnerId

`func (o *PageRequestResponse) GetOwnerId() string`

GetOwnerId returns the OwnerId field if non-nil, zero value otherwise.

### GetOwnerIdOk

`func (o *PageRequestResponse) GetOwnerIdOk() (*string, bool)`

GetOwnerIdOk returns a tuple with the OwnerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOwnerId

`func (o *PageRequestResponse) SetOwnerId(v string)`

SetOwnerId sets OwnerId field to given value.

### HasOwnerId

`func (o *PageRequestResponse) HasOwnerId() bool`

HasOwnerId returns a boolean if a field has been set.

### GetRequestType

`func (o *PageRequestResponse) GetRequestType() string`

GetRequestType returns the RequestType field if non-nil, zero value otherwise.

### GetRequestTypeOk

`func (o *PageRequestResponse) GetRequestTypeOk() (*string, bool)`

GetRequestTypeOk returns a tuple with the RequestType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestType

`func (o *PageRequestResponse) SetRequestType(v string)`

SetRequestType sets RequestType field to given value.

### HasRequestType

`func (o *PageRequestResponse) HasRequestType() bool`

HasRequestType returns a boolean if a field has been set.

### GetTitle

`func (o *PageRequestResponse) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *PageRequestResponse) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *PageRequestResponse) SetTitle(v string)`

SetTitle sets Title field to given value.

### HasTitle

`func (o *PageRequestResponse) HasTitle() bool`

HasTitle returns a boolean if a field has been set.

### GetDescription

`func (o *PageRequestResponse) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *PageRequestResponse) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *PageRequestResponse) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *PageRequestResponse) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetPageUrl

`func (o *PageRequestResponse) GetPageUrl() string`

GetPageUrl returns the PageUrl field if non-nil, zero value otherwise.

### GetPageUrlOk

`func (o *PageRequestResponse) GetPageUrlOk() (*string, bool)`

GetPageUrlOk returns a tuple with the PageUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPageUrl

`func (o *PageRequestResponse) SetPageUrl(v string)`

SetPageUrl sets PageUrl field to given value.

### HasPageUrl

`func (o *PageRequestResponse) HasPageUrl() bool`

HasPageUrl returns a boolean if a field has been set.

### SetPageUrlNil

`func (o *PageRequestResponse) SetPageUrlNil(b bool)`

 SetPageUrlNil sets the value for PageUrl to be an explicit nil

### UnsetPageUrl
`func (o *PageRequestResponse) UnsetPageUrl()`

UnsetPageUrl ensures that no value is present for PageUrl, not even an explicit nil
### GetLogoUrl

`func (o *PageRequestResponse) GetLogoUrl() string`

GetLogoUrl returns the LogoUrl field if non-nil, zero value otherwise.

### GetLogoUrlOk

`func (o *PageRequestResponse) GetLogoUrlOk() (*string, bool)`

GetLogoUrlOk returns a tuple with the LogoUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLogoUrl

`func (o *PageRequestResponse) SetLogoUrl(v string)`

SetLogoUrl sets LogoUrl field to given value.

### HasLogoUrl

`func (o *PageRequestResponse) HasLogoUrl() bool`

HasLogoUrl returns a boolean if a field has been set.

### SetLogoUrlNil

`func (o *PageRequestResponse) SetLogoUrlNil(b bool)`

 SetLogoUrlNil sets the value for LogoUrl to be an explicit nil

### UnsetLogoUrl
`func (o *PageRequestResponse) UnsetLogoUrl()`

UnsetLogoUrl ensures that no value is present for LogoUrl, not even an explicit nil
### GetCreatedAt

`func (o *PageRequestResponse) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *PageRequestResponse) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *PageRequestResponse) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *PageRequestResponse) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



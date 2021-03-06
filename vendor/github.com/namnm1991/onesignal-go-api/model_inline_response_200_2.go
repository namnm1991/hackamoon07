/*
OneSignal

A powerful way to send personalized messages at scale and build effective customer engagement strategies. Learn more at onesignal.com

API version: 1.0.2
Contact: devrel@onesignal.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package onesignal

import (
	"encoding/json"
)

// InlineResponse2002 struct for InlineResponse2002
type InlineResponse2002 struct {
	Success *string `json:"success,omitempty"`
	DestinationUrl *string `json:"destination_url,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _InlineResponse2002 InlineResponse2002

// NewInlineResponse2002 instantiates a new InlineResponse2002 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInlineResponse2002() *InlineResponse2002 {
	this := InlineResponse2002{}
	return &this
}

// NewInlineResponse2002WithDefaults instantiates a new InlineResponse2002 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInlineResponse2002WithDefaults() *InlineResponse2002 {
	this := InlineResponse2002{}
	return &this
}

// GetSuccess returns the Success field value if set, zero value otherwise.
func (o *InlineResponse2002) GetSuccess() string {
	if o == nil || o.Success == nil {
		var ret string
		return ret
	}
	return *o.Success
}

// GetSuccessOk returns a tuple with the Success field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse2002) GetSuccessOk() (*string, bool) {
	if o == nil || o.Success == nil {
		return nil, false
	}
	return o.Success, true
}

// HasSuccess returns a boolean if a field has been set.
func (o *InlineResponse2002) HasSuccess() bool {
	if o != nil && o.Success != nil {
		return true
	}

	return false
}

// SetSuccess gets a reference to the given string and assigns it to the Success field.
func (o *InlineResponse2002) SetSuccess(v string) {
	o.Success = &v
}

// GetDestinationUrl returns the DestinationUrl field value if set, zero value otherwise.
func (o *InlineResponse2002) GetDestinationUrl() string {
	if o == nil || o.DestinationUrl == nil {
		var ret string
		return ret
	}
	return *o.DestinationUrl
}

// GetDestinationUrlOk returns a tuple with the DestinationUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InlineResponse2002) GetDestinationUrlOk() (*string, bool) {
	if o == nil || o.DestinationUrl == nil {
		return nil, false
	}
	return o.DestinationUrl, true
}

// HasDestinationUrl returns a boolean if a field has been set.
func (o *InlineResponse2002) HasDestinationUrl() bool {
	if o != nil && o.DestinationUrl != nil {
		return true
	}

	return false
}

// SetDestinationUrl gets a reference to the given string and assigns it to the DestinationUrl field.
func (o *InlineResponse2002) SetDestinationUrl(v string) {
	o.DestinationUrl = &v
}

func (o InlineResponse2002) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Success != nil {
		toSerialize["success"] = o.Success
	}
	if o.DestinationUrl != nil {
		toSerialize["destination_url"] = o.DestinationUrl
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return json.Marshal(toSerialize)
}

func (o *InlineResponse2002) UnmarshalJSON(bytes []byte) (err error) {
	varInlineResponse2002 := _InlineResponse2002{}

	if err = json.Unmarshal(bytes, &varInlineResponse2002); err == nil {
		*o = InlineResponse2002(varInlineResponse2002)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "success")
		delete(additionalProperties, "destination_url")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableInlineResponse2002 struct {
	value *InlineResponse2002
	isSet bool
}

func (v NullableInlineResponse2002) Get() *InlineResponse2002 {
	return v.value
}

func (v *NullableInlineResponse2002) Set(val *InlineResponse2002) {
	v.value = val
	v.isSet = true
}

func (v NullableInlineResponse2002) IsSet() bool {
	return v.isSet
}

func (v *NullableInlineResponse2002) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInlineResponse2002(val *InlineResponse2002) *NullableInlineResponse2002 {
	return &NullableInlineResponse2002{value: val, isSet: true}
}

func (v NullableInlineResponse2002) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInlineResponse2002) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



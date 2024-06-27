package model

import "time"

type DeviceState = string

type DeviceMethod = string

type Device struct {
	DeviceBase
	Created  time.Time      `json:"created"`
	Updated  time.Time      `json:"updated"`
	UserData DeviceUserData `json:"user_data,omitempty"`
}

type DeviceBase struct {
	ID  string `json:"id"`
	Ref string `json:"ref"`
	DeviceData
}

type DeviceData struct {
	Name       string            `json:"name"`
	State      DeviceState       `json:"state"`
	Type       string            `json:"type"`
	Attributes []DeviceAttribute `json:"attributes"`
}

type DeviceUserData struct {
	Name       string            `json:"name"`
	Updated    time.Time         `json:"updated"`
	Attributes []DeviceAttribute `json:"attributes"`
}

type DeviceAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DeviceMessage struct {
	Method   DeviceMethod `json:"method"`
	DeviceID string       `json:"device_id"`
	Data     *DeviceData  `json:"data"`
}

type DevicesFilter struct {
	IDs   []string
	State string
	Type  string
	Ref   string
}

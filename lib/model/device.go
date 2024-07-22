package model

import "time"

type DeviceState = string

type DeviceMethod = string

type DeviceBase struct {
	DeviceData
	UserData DeviceUserData `json:"user_data,omitempty"`
}

type DeviceData struct {
	DeviceDataBase
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type DeviceDataBase struct {
	ID         string            `json:"id"`
	Ref        string            `json:"ref"`
	Name       string            `json:"name"`
	State      DeviceState       `json:"state"`
	Type       string            `json:"type"`
	Attributes []DeviceAttribute `json:"attributes"`
}

type DeviceUserDataBase struct {
	Name       string            `json:"name"`
	Attributes []DeviceAttribute `json:"attributes"`
}

type DeviceUserData struct {
	DeviceUserDataBase
	Updated time.Time `json:"updated"`
}

type DeviceAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DevicesFilter struct {
	IDs   []string
	State string
	Type  string
	Ref   string
}

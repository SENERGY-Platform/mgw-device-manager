package model

type DeviceMessage struct {
	Method   DeviceMethod       `json:"method"`
	DeviceID string             `json:"device_id"`
	Data     *DeviceMessageData `json:"data"`
}

type DeviceMessageData struct {
	Name       string            `json:"name"`
	State      DeviceState       `json:"state"`
	Type       string            `json:"device_type"`
	Attributes []DeviceAttribute `json:"attributes"`
}

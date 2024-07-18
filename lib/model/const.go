package model

const (
	Online  DeviceState = "online"
	Offline DeviceState = "offline"
)

const (
	Set    DeviceMethod = "set"
	Delete DeviceMethod = "delete"
)

const (
	DevicesPath = "devices"
	SrvInfoPath = "info"
)

const (
	HeaderRequestID = "X-Request-ID"
	HeaderApiVer    = "X-Api-Version"
	HeaderSrvName   = "X-Service"
)

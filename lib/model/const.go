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
	DevicesPath    = "devices"
	RestrictedPath = "restricted"
	SrvInfoPath    = "info"
)

const (
	HeaderRequestID = "X-Request-ID"
	HeaderApiVer    = "X-Api-Version"
	HeaderSrvName   = "X-Service"
)

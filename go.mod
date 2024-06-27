module github.com/SENERGY-Platform/mgw-device-manager

go 1.22

require (
	github.com/SENERGY-Platform/go-service-base/sql-db-hdl v0.0.1
	github.com/SENERGY-Platform/go-service-base/srv-info-hdl v0.0.3
	github.com/SENERGY-Platform/go-service-base/util v1.0.0
	github.com/SENERGY-Platform/go-service-base/watchdog v0.4.2
	github.com/SENERGY-Platform/mgw-device-manager/lib v0.0.0-00010101000000-000000000000
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/y-du/go-log-level v1.0.0
)

require (
	github.com/SENERGY-Platform/go-service-base/srv-info-hdl/lib v0.0.2 // indirect
	github.com/y-du/go-env-loader v0.5.2 // indirect
)

replace github.com/SENERGY-Platform/mgw-device-manager/lib => ./lib

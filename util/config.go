package util

import (
	"github.com/SENERGY-Platform/go-service-base/config-hdl"
	sb_logger "github.com/SENERGY-Platform/go-service-base/logger"
	envldr "github.com/y-du/go-env-loader"
	"github.com/y-du/go-log-level/level"
	"reflect"
)

type DatabaseConfig struct {
	Timeout    int64  `json:"timeout" env_var:"DB_TIMEOUT"`
	Path       string `json:"path" env_var:"DB_PATH"`
	SchemaPath string `json:"schema_path" env_var:"DB_SCHEMA_PATH"`
}

type MqttClientConfig struct {
	Server            string `json:"server" env_var:"MQTT_SERVER"`
	KeepAlive         int64  `json:"keep_alive" env_var:"MQTT_KEEP_ALIVE"`
	PingTimeout       int64  `json:"ping_timeout" env_var:"MQTT_PING_TIMEOUT"`
	ConnectTimeout    int64  `json:"connect_timeout" env_var:"MQTT_CONNECT_TIMEOUT"`
	ConnectRetryDelay int64  `json:"connect_retry_delay" env_var:"MQTT_CONNECT_RETRY_DELAY"`
	MaxReconnectDelay int64  `json:"max_reconnect_delay" env_var:"MQTT_MAX_RECONNECT_DELAY"`
	WaitTimeout       int64  `json:"wait_timeout" env_var:"MQTT_WAIT_TIMEOUT"`
	QOSLevel          byte   `json:"qos_level" env_var:"MQTT_QOS_LEVEL"`
}

type LoggerConfig struct {
	Level        level.Level `json:"level" env_var:"LOGGER_LEVEL"`
	Utc          bool        `json:"utc" env_var:"LOGGER_UTC"`
	Path         string      `json:"path" env_var:"LOGGER_PATH"`
	FileName     string      `json:"file_name" env_var:"LOGGER_FILE_NAME"`
	Terminal     bool        `json:"terminal" env_var:"LOGGER_TERMINAL"`
	Microseconds bool        `json:"microseconds" env_var:"LOGGER_MICROSECONDS"`
	Prefix       string      `json:"prefix" env_var:"LOGGER_PREFIX"`
}

type Config struct {
	Logger          LoggerConfig     `json:"logger" env_var:"LOGGER_CONFIG"`
	Database        DatabaseConfig   `json:"database" env_var:"DATABASE_CONFIG"`
	MqttClient      MqttClientConfig `json:"mqtt_client" env_var:"MQTT_CLIENT_CONFIG"`
	MGWDeploymentID string           `json:"mgw_deployment_id" env_var:"MGW_DID"`
	MQTTLog         bool             `json:"mqtt_log" env_var:"MQTT_LOG"`
	MQTTDebugLog    bool             `json:"mqtt_debug_log" env_var:"MQTT_DEBUG_LOG"`
	ServerPort      uint             `json:"server_port" env_var:"SERVER_PORT"`
	MessageBuffer   int              `json:"message_buffer" env_var:"MESSAGE_BUFFER"`
}

var defaultMqttClientConfig = MqttClientConfig{
	KeepAlive:         30000000000, // 30s
	PingTimeout:       10000000000, // 10s
	ConnectTimeout:    30000000000, // 30s
	ConnectRetryDelay: 30000000000, // 30s
	MaxReconnectDelay: 30000000000, // 30s
	WaitTimeout:       5000000000,  // 5s
	QOSLevel:          2,
}

func NewConfig(path string) (*Config, error) {
	cfg := Config{
		Logger: LoggerConfig{
			Level:        level.Warning,
			Utc:          true,
			Microseconds: true,
			Terminal:     true,
		},
		Database: DatabaseConfig{
			Timeout:    5000000000,
			Path:       "/opt/device-manager/data",
			SchemaPath: "include/storage_schema.sql",
		},
		MqttClient:    defaultMqttClientConfig,
		ServerPort:    80,
		MessageBuffer: 50000,
	}
	err := config_hdl.Load(&cfg, nil, map[reflect.Type]envldr.Parser{reflect.TypeOf(level.Off): sb_logger.LevelParser}, nil, path)
	return &cfg, err
}

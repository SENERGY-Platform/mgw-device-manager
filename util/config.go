package util

import (
	sb_util "github.com/SENERGY-Platform/go-service-base/util"
	"github.com/y-du/go-log-level/level"
)

type DatabaseConfig struct {
	Name       string `json:"name" env_var:"DB_NAME"`
	Timeout    int64  `json:"timeout" env_var:"DB_TIMEOUT"`
	Path       string `json:"path" env_var:"PATH"`
	SchemaPath string `json:"schema_path" env_var:"DB_SCHEMA_PATH"`
}

type MqttClientConfig struct {
	Server            string `json:"server" env_var:"LOCAL_MQTT_SERVER"`
	KeepAlive         int64  `json:"keep_alive" env_var:"LOCAL_MQTT_KEEP_ALIVE"`
	PingTimeout       int64  `json:"ping_timeout" env_var:"LOCAL_MQTT_PING_TIMEOUT"`
	ConnectTimeout    int64  `json:"connect_timeout" env_var:"LOCAL_MQTT_CONNECT_TIMEOUT"`
	ConnectRetryDelay int64  `json:"connect_retry_delay" env_var:"LOCAL_MQTT_CONNECT_RETRY_DELAY"`
	MaxReconnectDelay int64  `json:"max_reconnect_delay" env_var:"LOCAL_MQTT_MAX_RECONNECT_DELAY"`
	WaitTimeout       int64  `json:"wait_timeout" env_var:"LOCAL_MQTT_WAIT_TIMEOUT"`
	QOSLevel          byte   `json:"qos_level" env_var:"LOCAL_MQTT_QOS_LEVEL"`
}

type Config struct {
	Logger          sb_util.LoggerConfig `json:"logger" env_var:"LOGGER_CONFIG"`
	Database        DatabaseConfig       `json:"database" env_var:"DATABASE_CONFIG"`
	MqttClient      MqttClientConfig     `json:"mqtt_client" env_var:"MQTT_CLIENT_CONFIG"`
	MGWDeploymentID string               `json:"mgw_deployment_id" env_var:"MGW_DID"`
	MQTTLog         bool                 `json:"mqtt_log" env_var:"MQTT_LOG"`
	MQTTDebugLog    bool                 `json:"mqtt_debug_log" env_var:"MQTT_DEBUG_LOG"`
	ServerPort      uint                 `json:"server_port" env_var:"SERVER_PORT"`
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
		Logger: sb_util.LoggerConfig{
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
		MqttClient: defaultMqttClientConfig,
		ServerPort: 80,
	}
	err := sb_util.LoadConfig(path, &cfg, nil, nil, nil)
	return &cfg, err
}

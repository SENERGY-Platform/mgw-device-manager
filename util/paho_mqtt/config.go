package paho_mqtt

import (
	"github.com/SENERGY-Platform/mgw-device-manager/util"
	"github.com/eclipse/paho.mqtt.golang"
	"time"
)

func SetClientOptions(co *mqtt.ClientOptions, clientID string, mqttConf util.MqttClientConfig) {
	co.AddBroker(mqttConf.Server)
	co.SetClientID(clientID)
	co.SetKeepAlive(time.Duration(mqttConf.KeepAlive))
	co.SetPingTimeout(time.Duration(mqttConf.PingTimeout))
	co.SetConnectTimeout(time.Duration(mqttConf.ConnectTimeout))
	co.SetConnectRetryInterval(time.Duration(mqttConf.ConnectRetryDelay))
	co.SetMaxReconnectInterval(time.Duration(mqttConf.MaxReconnectDelay))
	co.SetWriteTimeout(time.Second * 5)
	co.ConnectRetry = true
	co.AutoReconnect = true
}

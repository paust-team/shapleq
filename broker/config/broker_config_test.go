package config

import (
	"fmt"
	"github.com/paust-team/shapleq/common"
	"testing"
)

func TestBrokerConfigLoad(t *testing.T) {
	brokerConfig := NewBrokerConfig()

	defaultZKTimeout := brokerConfig.ZKTimeout()

	brokerConfig.Load("./config.yml")
	if brokerConfig.ZKTimeout() == defaultZKTimeout {
		t.Errorf("value in file is not wrapped")
	}
}

func TestBrokerConfigSet(t *testing.T) {
	brokerConfig := NewBrokerConfig().Load("./config.yml")

	if brokerConfig.Port() != common.DefaultBrokerPort {
		t.Errorf("wrong default value")
	}

	var expectedPort uint = 11010
	brokerConfig.SetPort(expectedPort)

	if brokerConfig.Port() != expectedPort {
		t.Errorf("value is not set")
	}
}

func TestBrokerConfigStructured(t *testing.T) {
	brokerConfig := NewBrokerConfig().Load("./config.yml")

	expectedHost := "172.0.0.1"
	var expectedPort uint = 10000

	brokerConfig.SetZKHost(expectedHost)
	brokerConfig.SetZKPort(expectedPort)

	if brokerConfig.ZKAddr() != fmt.Sprintf("%s:%d", expectedHost, expectedPort) {
		t.Errorf("value is not set")
	}
}

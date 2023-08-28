package apis_test

import (
	"testing"

	"github.com/jaksonlin/go-kefctl/apis"
)

func TestPower(t *testing.T) {
	apis.SetPower("10.0.0.2", 50001, false)
}

func TestVolume(t *testing.T) {
	apis.GetVolume("10.0.0.2", 50001)
	apis.SetVolume("10.0.0.2", 50001, 50)
	apis.GetVolume("10.0.0.2", 50001)
}

func TestSwitchSource(t *testing.T) {
	apis.SwitchInput("10.0.0.2", 50001, "off", "60", "wifi")
}

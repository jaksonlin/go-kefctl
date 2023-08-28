package apis

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

var GET_SOURCE = []byte{0x47, 0x30, 0x80}

// sendCommand send control command to the speaker
func sendCommand(host string, port int, command []byte) []byte {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("Error connecting:", err)
		return nil
	}
	defer conn.Close()
	tcpConn := conn.(*net.TCPConn)
	tcpConn.SetNoDelay(true)

	_, err = conn.Write(command)
	if err != nil {
		fmt.Println("Error sending command:", err)
		return nil
	}

	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		fmt.Println("Error receiving response:", err)
		return nil
	}

	return response
}

// SetPower turn the speaker on or off
func SetPower(host string, port int, powerOff bool) {
	response := sendCommand(host, port, GET_SOURCE)
	sourceBits := fmt.Sprintf("%08b", response[3]) // force a 8 bits string, for 6 bits, the control failed

	if powerOff && sourceBits[2:4] == "00" {
		sourceBits = sourceBits[:2] + "01" + sourceBits[4:]
		source, _ := strconv.ParseUint(sourceBits, 2, 8)
		sendCommand(host, port, []byte{0x53, 0x30, 0x81, byte(source)})
	}

	if powerOff {
		sourceBits = "1" + sourceBits[1:]
	} else {
		sourceBits = "0" + sourceBits[1:]
	}

	source, _ := strconv.ParseUint(sourceBits, 2, 8)
	sendCommand(host, port, []byte{0x53, 0x30, 0x81, byte(source)})

}

// GetVolume get the current volume of the speaker
func GetVolume(host string, port int) int {
	getVolume := []byte{0x47, 0x25, 0x80}
	response := sendCommand(host, port, getVolume)
	volume := int(response[3])
	return volume
}

// SetVolume set the volume of the speaker, volume range from 0 to 100
func SetVolume(host string, port int, volume int) {
	response := sendCommand(host, port, []byte{0x53, 0x25, 0x81, byte(volume)})
	fmt.Println(response)
}

// SwitchInput contains a set of control for different components, as they share the same control command, so we put them together
func SwitchInput(host string, port int, inverse string, standby string, input string) {

	response := sendCommand(host, port, GET_SOURCE)
	sourceBits := fmt.Sprintf("%06b", response[3]) // here we need a 6 bits string instead of 8 bits

	// Always power on
	sourceBits = "0" + sourceBits[1:]

	// 2 bits for inverse

	if inverse != "" {
		val := strings.ToLower(inverse)
		if val == "on" {
			sourceBits = sourceBits[:1] + "1" + sourceBits[2:]
		} else if val == "off" {
			sourceBits = sourceBits[:1] + "0" + sourceBits[2:]
		} else {
			log.Fatal("Unknown value:", val)
		}
	}

	// 2 bits for standby
	if standby != "" {
		if standby == "60" {
			sourceBits = sourceBits[:2] + "01" + sourceBits[4:]
		} else if standby == "0" {
			sourceBits = sourceBits[:2] + "10" + sourceBits[4:]
		} else if standby == "20" {
			sourceBits = sourceBits[:2] + "00" + sourceBits[4:]
		} else {
			log.Fatal("Unknown value:", standby)
		}
	}

	// 4 bits for input
	if input != "" {
		val := strings.ToLower(input)
		if val == "wifi" {
			sourceBits = sourceBits[:4] + "0010"
		} else if val == "usb" {
			sourceBits = sourceBits[:4] + "1100"
		} else if val == "bluetooth" {
			sourceBits = sourceBits[:4] + "1001"
		} else if val == "aux" {
			sourceBits = sourceBits[:4] + "1010"
		} else if val == "optical" {
			sourceBits = sourceBits[:4] + "1011"
		} else {
			log.Fatal("Unknown value:", val)
		}
	}

	source, _ := strconv.ParseUint(sourceBits, 2, 8)
	sendCommand(host, port, []byte{0x53, 0x30, 0x81, byte(source)})

}

// ShowStatus show the current status of the speaker, CLI helper function
func ShowStatus(host string, port int) {
	volume := GetVolume(host, port)
	muted := "No"
	if volume >= 128 {
		muted = "Yes"
		volume -= 128
	}
	fmt.Println("Volume: ", volume, "%")
	fmt.Println("Muted: ", muted)

	response := sendCommand(host, port, GET_SOURCE)
	sourceBits := fmt.Sprintf("%08b", response[3])

	power := "On"
	if sourceBits[0] == '1' {
		power = "Off"
	}

	inverse := "Off"
	if sourceBits[1] == '1' {
		inverse = "On"
	}

	standby := "Unknown"
	if sourceBits[2:4] == "00" {
		standby = "20 Minutes"
	} else if sourceBits[2:4] == "01" {
		standby = "60 Minutes"
	} else if sourceBits[2:4] == "10" {
		standby = "Never"
	}

	input := "Unknown"
	if sourceBits[4:8] == "0010" {
		input = "Wifi"
	} else if sourceBits[4:8] == "1100" {
		input = "USB"
	} else if sourceBits[4:8] == "1001" {
		input = "Bluetooth (paired)"
	} else if sourceBits[4:8] == "1111" {
		input = "Bluetooth (unpaired)"
	} else if sourceBits[4:8] == "1010" {
		input = "Aux"
	} else if sourceBits[4:8] == "1011" {
		input = "Optical"
	}

	fmt.Println("Source: ", input)
	fmt.Println("Standby: ", standby)
	fmt.Println("Inverse: ", inverse)
	fmt.Println("Power: ", power)

}

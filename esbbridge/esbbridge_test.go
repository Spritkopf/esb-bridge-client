package esbbridge

import (
	"fmt"
	"testing"
)

var testPipelineAddress = [5]byte{111, 111, 111, 111, 1}
var testDevice string = "/dev/ttyACM0"

//TestOpenSuccess tests that the virtual COM port can be opened
func TestOpenSuccess(t *testing.T) {
	err := Open("/dev/ttyACM0")

	if err != nil {
		t.Fatalf(err.Error())
	}

	Close()
}

// TestGetFwVersionNotOpen tests error handling when not connected
func TestGetFwVersionNotOpen(t *testing.T) {

	_, err := GetFwVersion()

	if err == nil {
		t.Fatalf("GetFwVersion should return an error when not connected (i.e. Open() was not called beforehand)")
	}

	Close()

}

// TestGetFwVersion tests correct read of firmware version
func TestGetFwVersion(t *testing.T) {

	err := Open(testDevice)

	if err != nil {
		t.Fatalf(err.Error())
	}

	version, err := GetFwVersion()

	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("Version: %v\n", version)
	Close()

}

// TestTransferNotOpen tests error handling when not connected
func TestTransferNotOpen(t *testing.T) {
	_, err := Transfer(testPipelineAddress, nil)

	if err == nil {
		t.Fatalf("Transfer should return an error when not connected (i.e. Open() was not called beforehand)")
	}

	Close()
}

// TestTransferPayloadSize tests error handling for too large and too short payloads
func TestTransferPayloadSize(t *testing.T) {
	var veryLongPayload [64]byte
	var veryShortPayload [2]byte

	Open(testDevice)

	_, err := Transfer(testPipelineAddress, veryLongPayload[:])

	if err == nil {
		t.Fatalf("Transfer should return an error when Payload is longer than 32 bytes")
	}

	_, err = Transfer(testPipelineAddress, veryShortPayload[:])
	if err == nil {
		t.Fatalf("Transfer should return an error when Payload is shorter than 1 bytes")
	}

	_, err = Transfer(testPipelineAddress, nil)
	if err == nil {
		t.Fatalf("Transfer should return an error when Payload is nil")
	}

	Close()
}

// TestTransfer tests the transfer of ESB packages by requesting the firware version of a supported device
// Note: the ESB command ID ESB_CMD_VERSION (0x10) should be common to all the custom esb compatible devices
func TestTransfer(t *testing.T) {

	errOpen := Open(testDevice)

	if errOpen != nil {
		t.Fatalf("Open() failed with error %v", errOpen)
	}

	payload := []byte{0x10}
	ansPayload, err := Transfer(testPipelineAddress, payload)

	if err != nil {
		t.Fatalf("Transfer() failed with error %v", err)
	}

	if len(ansPayload) != 5 {
		t.Fatalf("Answer payload has unexpected size, Got %v", ansPayload)
	}

	Close()
}

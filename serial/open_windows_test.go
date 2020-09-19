package serial

import (
	"testing"

	"github.com/windhooked/serial"
)

func TestXonXoff(t *testing.T) {
	// open
	f, err := serial.Open(OpenOptions{
		PortName:              "CNA0",
		BaudRate:              115200,
		DataBits:              8,
		StopBits:              1,
		ParityMode:            0,
		RTSCTSFlowControl:     false,
		InterCharacterTimeout: 100,
		MinimumReadSize:       1,

		XFlowControl: &XFlowControl{
			TXContinueOnXOFF: true,
			InX:              true,
			OutX:             true,
		},
	})
	if err != nil {
		t.Fatalf("Error opening serial port: ", err)
	}
	defer f.Close()

}

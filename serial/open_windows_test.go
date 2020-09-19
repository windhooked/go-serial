package serial

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

// requires com0com virtual port or similar
func TestXonXoff(t *testing.T) {
	// open
	portA, err := serial.Open(OpenOptions{
		PortName:              "CNCA0",
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
	defer portA.Close()

	portB, err := serial.Open(OpenOptions{
		PortName:              "CNCB0",
		BaudRate:              115200,
		DataBits:              8,
		StopBits:              1,
		ParityMode:            0,
		RTSCTSFlowControl:     false,
		InterCharacterTimeout: 100,
		MinimumReadSize:       1,

		//RTSCTSFlowControl: true,

		XFlowControl: &XFlowControl{
			TXContinueOnXOFF: true,
			InX:              true,
			OutX:             true,
		},
	})
	if err != nil {
		t.Fatalf("Error opening serial port: ", err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func(portA serial.Port) {
		buff := make([]byte, 512)

		for _, v := range sendBuf {
			v = 0x55
		}
		loopCntr := 0
		for {
			select {
			case <-c:
				os.Exit(1)
			default:
				w, err := portA.Write(out)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
				if w > 0 {
					s, err := portA.GetModemStatusBits()
					fmt.Printf("%v %v[%+v]: rx t (%v) (%v) % x\n", "A", loopCntr, s, len(out), w, out)
					loopCntr++
				}
			}
		}
	}(portA)

	go func(portB serial.Port) {
		buff := make([]byte, 4096)
		for {
			select {
			case <-c:
				os.Exit(1)
			default:
				n, err := portB.Read(buff)
				if n > 0 {
					out := buff[:n]
					s, err := portB.GetModemStatusBits()
					fmt.Printf("%v[%+v]: rx m (%v) (%v) % x\n", "B", s, len(out), w, out)
				} else if err != nil && err != io.EOF {
					fmt.Printf("%v\n", err)
				}

			}
		}
	}(portB)

	<-c

}

package serial

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

const (
	XOFF = 0x11
	XON  = 0x13
)

// requires com0com virtual port or similar
func TestXonXoff(t *testing.T) {
	// open
	portA, err := Open(OpenOptions{
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
			InX:              false,
			OutX:             true,
		},
	})
	if err != nil {
		t.Fatalf("Error opening serial port: %v", err)
	}
	defer portA.Close()

	portB, err := Open(OpenOptions{
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
		t.Fatalf("Error opening serial port: %v", err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	var transmitted []byte
	var received []byte
	go func(portA io.ReadWriteCloser) {
		txBuff := make([]byte, 512)

		for k, _ := range txBuff {
			txBuff[k] = 0x55
		}
		txBuff[10] = XON
		txBuff[20] = XOFF
		loopCntr := 0
	DONE:
		for {
			select {
			case <-c:
				os.Exit(1)
			default:
				if loopCntr > 5 {
					break DONE
				}
				in := []byte{0x00, 0x00}
				n, err := portB.Read(in)
				if n > 0 {
					fmt.Printf("tx in %x\n", in)
				}
				w, err := portA.Write(txBuff)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
				if w > 0 {
					transmitted = append(transmitted, txBuff...)
					//s, err := portA.GetModemStatusBits()
					fmt.Printf("%v: tx total:%v %v % x\n", "A", loopCntr*len(txBuff), w, txBuff)
					loopCntr++
				}
			}
		}
	}(portA)

	go func(portB io.ReadWriteCloser) {
		buff := make([]byte, 200)
		loopCntr := 0
		for {
			select {
			case <-c:
				os.Exit(1)
			default:
				// simulate busy slow reader
				time.Sleep(time.Millisecond * 100)
				n, err := portB.Read(buff)
				if n > 0 {
					out := buff[:n]
					received = append(received, out[:n]...)
					//s, err := portB.GetModemStatusBits()
					fmt.Printf("%v: rx total:%v %v % x\n", "B", loopCntr*len(out), n, out)
					loopCntr++
				} else if err != nil && err != io.EOF {
					fmt.Printf("%v\n", err)
				}

			}
		}
	}(portB)

	fmt.Printf("% x", transmitted)
	fmt.Printf("% x", received)
	<-c

}

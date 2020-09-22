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
		//https://docs.microsoft.com/en-us/windows/win32/api/winbase/ns-winbase-dcb
		XFlowControl: &XFlowControl{
			//fTXContinueOnXoff
			//If this member is TRUE, transmission continues after the input buffer has come within XoffLim bytes of being full
			// and the driver has transmitted the XoffChar character to stop receiving bytes. If this member is FALSE, transmission
			//does not continue until the input buffer is within XonLim bytes of being empty and the driver has transmitted the XonChar
			//character to resume reception.
			TXContinueOnXOFF: true,
			//fOutX
			//Indicates whether XON/XOFF flow control is used during transmission. If this member is TRUE,
			//transmission stops when the XoffChar character is received and starts again when the XonChar character is received.
			OutX: false,
			//fInX
			//Indicates whether XON/XOFF flow control is used during reception. If this member is TRUE,
			//the XoffChar character is sent when the input buffer comes within XoffLim bytes of being full,
			//and the XonChar character is sent when the input buffer comes within XonLim bytes of being empty.
			InX: false,
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
		txBuff[10] = XOFF
		txBuff[11] = XOFF
		txBuff[20] = XON
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

	<-c
	fmt.Printf("%v % x\n", len(transmitted), transmitted[10:13])
	fmt.Printf("missing, xon,xoff driver is intercepting %v % x\n", len(received), received[10:13])

}

// requires com0com virtual port or similar
func TestHardwareFlow(t *testing.T) {
	// open
	portA, err := Open(OpenOptions{
		PortName:              "COM12",
		BaudRate:              115200,
		DataBits:              8,
		StopBits:              1,
		ParityMode:            0,
		RTSCTSFlowControl:     true,
		InterCharacterTimeout: 100,
		MinimumReadSize:       1,
	})
	if err != nil {
		t.Fatalf("Error opening serial port: %v", err)
	}
	defer portA.Close()

	portB, err := Open(OpenOptions{
		PortName:              "COM13",
		BaudRate:              115200,
		DataBits:              8,
		StopBits:              1,
		ParityMode:            0,
		RTSCTSFlowControl:     true,
		InterCharacterTimeout: 100,
		MinimumReadSize:       1,
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
		txBuff[10] = XOFF
		txBuff[11] = XOFF
		txBuff[20] = XON
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

	<-c
	fmt.Printf("%v % x\n", len(transmitted), transmitted[10:13])
	fmt.Printf("missing, xon,xoff driver is intercepting %v % x\n", len(received), received[10:13])

}

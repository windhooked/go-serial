go-serial
=========

forked of github.com/jacobsa/go-serial/serial now with XON XOFF flow control on windows
This is a package that allows you to read from and write to serial ports in Go.


Use
---

Set up a `serial.OpenOptions` struct, then call `serial.Open`. For example:

````go
    import (
      "fmt"
      "log"
      "github.com/windhooked/go-serial/serial"
    )

    ...

    // Set up options.
    options := serial.OpenOptions{
      PortName: "COM4",
      BaudRate: 19200,
      DataBits: 8,
      StopBits: 1,
      MinimumReadSize: 4,
      XFlowControl: &XFlowControl{
			TXContinueOnXOFF: true,
			InX:              true,
			OutX:             true,
		},
    }

    // Open the port.
    port, err := serial.Open(options)
    if err != nil {
      log.Fatalf("serial.Open: %v", err)
    }

    // Make sure to close it later.
    defer port.Close()

    // Write 4 bytes to the port.
    b := []byte{0x00, 0x01, 0x02, 0x03}
    n, err := port.Write(b)
    if err != nil {
      log.Fatalf("port.Write: %v", err)
    }

    fmt.Println("Wrote", n, "bytes.")
````
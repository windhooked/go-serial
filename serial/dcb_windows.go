package serial

import "fmt"

const (
	dcbBinary                uint32 = 0x00000001
	dcbParity                       = 0x00000002
	dcbOutXCTSFlow                  = 0x00000004
	dcbOutXDSRFlow                  = 0x00000008
	dcbDTRControlDisableMask        = ^uint32(0x00000030)
	dcbDTRControlEnable             = 0x00000010
	dcbDTRControlHandshake          = 0x00000020
	dcbDSRSensitivity               = 0x00000040
	dcbTXContinueOnXOFF             = 0x00000080
	dcbOutX                         = 0x00000100
	dcbInX                          = 0x00000200
	dcbErrorChar                    = 0x00000400
	dcbNull                         = 0x00000800
	dcbRTSControlDisableMask        = ^uint32(0x00003000)
	dcbRTSControlEnable             = 0x00001000
	dcbRTSControlHandshake          = 0x00002000
	dcbRTSControlToggle             = 0x00003000
	dcbAbortOnError                 = 0x00004000
)

type Dcb struct {
	dcb
}

type dcb struct {
	DCBlength uint32
	BaudRate  uint32

	// Flags field is a bitfield
	//  fBinary            :1
	//  fParity            :1
	//  fOutxCtsFlow       :1
	//  fOutxDsrFlow       :1
	//  fDtrControl        :2
	//  fDsrSensitivity    :1
	//  fTXContinueOnXoff  :1
	//  fOutX              :1
	//  fInX               :1
	//  fErrorChar         :1
	//  fNull              :1
	//  fRtsControl        :2
	//  fAbortOnError      :1
	//  fDummy2            :17
	Flags uint32

	wReserved  uint16
	XonLim     uint16
	XoffLim    uint16
	ByteSize   byte
	Parity     byte
	StopBits   byte
	XonChar    byte
	XoffChar   byte
	ErrorChar  byte
	EOFChar    byte
	EvtChar    byte
	wReserved1 uint16
}

func NewDcb(in Dcb) *Dcb {
	dcb := &in
	dcb.SetBinary() // always
	return dcb
}

func (p *Dcb) String() string {
	return fmt.Sprintf("% b", p.Flags)
}

func (p *Dcb) SetBinary() {
	p.Flags |= 0x01
}
func (p *Dcb) SetBaudRate(b int) {
	p.BaudRate = uint32(b)
}

// Flags field is a bitfield
//  fBinary            :1
func (p *Dcb) SetParity() {
}
func (p *Dcb) SetOutxCtsFlow() {
}
func (p *Dcb) SetOutxDsrFlow() {
}
func (p *Dcb) SetDtrControl()       {}
func (p *Dcb) SetDsrSensitivity()   {}
func (p *Dcb) SetTXContinueOnXoff() {}
func (p *Dcb) SetOutX()             {}
func (p *Dcb) SetInX()              {}
func (p *Dcb) SetErrorChar()        {}
func (p *Dcb) SetNull()             {}
func (p *Dcb) SetRtsControl() {
}
func (p *Dcb) SetAbortOnError() {
}

func (p *Dcb) SetXonLim(l int) {
	p.XonLim = uint16(l)
}
func (p *Dcb) SetXoffLim(l int) {
	p.XoffLim = uint16(l)

}
func (p *Dcb) SetByteSize() {} //byte
func (p *Dcb) SetStopBits() {} //byte
func (p *Dcb) SetXonChar()  {} //byte
func (p *Dcb) SetXoffChar() {} //byte
func (p *Dcb) SetEOFChar()  {} //byte
func (p *Dcb) SetEvtChar()  {} //byte

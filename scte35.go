package scte35

import "fmt"
import "os"
import "encoding/base64"
import "github.com/futzu/bitter"

// PktSz is the size of an MPEG-TS packet in bytes.
const PktSz = 188

// BufferSize is the size of a read when parsing files.
const BufferSize = 384 * PktSz

// Chk is generic catchall error checking.
func Chk(e error) {
	if e != nil {
		panic(e)
	}
}

// DeB64 decodes base64 strings.
func DeB64(b64 string) []byte {
	deb64, err := base64.StdEncoding.DecodeString(b64)
	Chk(err)
	return deb64
}

// IsIn is a test for slice membership
func IsIn(slice []uint8, val uint8) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// SCTE35Parser parses a slice of bytes for SCTE 35 data.
func SCTE35Parser(bites []byte) {
	var bitn bitter.Bitn
	bitn.Load(bites)
	var spi SpInfo
	spi.Decode(&bitn)
	var cmd SpCmd
	cmd.Decode(&bitn, spi.SpliceCommandType)
	spi.DescriptorLoopLength = bitn.AsUInt64(16)
	fmt.Printf("%+v\n", spi)
	fmt.Printf("%+v\n", cmd)
}

// PktParser is a parser for an MPEG-TS SCTE 35 packet
func PktParser(pkt []byte) {
	// pkt is the mpegts packet, pld is the packet payload
	pld := pkt[5:PktSz]
	magicbytes := [4]uint8{252, 48, 0, 255}
	// cmds is an array of valid SCTE 35 command codes
	cmds := []uint8{0, 5, 6, 7, 255}
	// compare the pld bytes 0,1,3,and 10 to the magicbytes 
	pldbytes := [4]uint8{pld[0], pld[1], pld[3], pld[10]}
	if pldbytes == magicbytes {
		// if pld byte 13 is a valid scte 35 command code
		if IsIn(cmds, pld[13]) {
			// then this is a scte 35 packet, parse it
			SCTE35Parser(pld)
		}
	}
}

// FileParser is a parser for an MPEG-TS file.
func FileParser(fname string) {
	file, err := os.Open(fname)
	Chk(err)
	defer file.Close()
	buffer := make([]byte, BufferSize)
	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			break
		}
		for i := 1; i <= (bytesread / PktSz); i++ {
			end := i * PktSz
			start := end - PktSz
			pkt := buffer[start:end]
			PktParser(pkt)
		}
	}
}

// SpInfo is the splice info section of the SCTE 35 cue.
type SpInfo struct {
	Name                   string
	TableId                string
	SectionSyntaxIndicator bool
	Private                bool
	Reserved               string
	SectionLength          uint64
	ProtocolVersion        uint64
	EncryptedPacket        bool
	EncryptionAlgorithm    uint64
	PtsAdjustment          float64
	CwIndex                string
	Tier                   string
	SpliceCommandLength    uint64
	SpliceCommandType      uint64
	DescriptorLoopLength   uint64
}

// Decode extracts bits for the splice info section values.
func (spi *SpInfo) Decode(bitn *bitter.Bitn) {
	spi.Name = "Splice Info Section"
	spi.TableId = bitn.AsHex(8)
	spi.SectionSyntaxIndicator = bitn.AsBool()
	spi.Private = bitn.AsBool()
	spi.Reserved = bitn.AsHex(2)
	spi.SectionLength = bitn.AsUInt64(12)
	spi.ProtocolVersion = bitn.AsUInt64(8)
	spi.EncryptedPacket = bitn.AsBool()
	spi.EncryptionAlgorithm = bitn.AsUInt64(6)
	spi.PtsAdjustment = bitn.As90k(33)
	spi.CwIndex = bitn.AsHex(8)
	spi.Tier = bitn.AsHex(12)
	spi.SpliceCommandLength = bitn.AsUInt64(12)
	spi.SpliceCommandType = bitn.AsUInt64(8)
}

// SpCmd is the splice command for the SCTE35 cue.
type SpCmd struct {
	Name                       string
	SpliceEventId              string
	SpliceEventCancelIndicator bool
	OutOfNetworkIndicator      bool
	ProgramSpliceFlag          bool
	DurationFlag               bool
	BreakAutoReturn            bool
	BreakDuration              float64
	SpliceImmediateFlag        bool
	TimeSpecifiedFlag          bool
	PTS                        float64
	ComponentCount             uint64
	Components                 []uint64
	UniqueProgramId            uint64
	AvailNum                   uint64
	AvailExpected              uint64
	Identifier                 uint64
}

// Decode the splice command values.
func (cmd *SpCmd) Decode(bitn *bitter.Bitn, cmdtype uint64) {
	if cmdtype == 0 {
		cmd.SpliceNull()
	}
	//4: Splice_Schedule,
	if cmdtype == 5 {
		cmd.SpliceInsert(bitn)
	}
	if cmdtype == 6 {
		cmd.TimeSignal(bitn)
	}
	if cmdtype == 7 {
		cmd.BandwidthReservation(bitn)
	}
	if cmdtype == 255 {
		cmd.PrivateCommand(bitn)
	}
}

// ParseBreak parses out the ad break duration values.
func (cmd *SpCmd) ParseBreak(bitn *bitter.Bitn) {
	cmd.BreakAutoReturn = bitn.AsBool()
	bitn.Forward(6)
	cmd.BreakDuration = bitn.As90k(33)
}

// SpliceTime parses out the PTS value as needed.
func (cmd *SpCmd) SpliceTime(bitn *bitter.Bitn) {
	cmd.TimeSpecifiedFlag = bitn.AsBool()
	if cmd.TimeSpecifiedFlag {
		bitn.Forward(6)
		cmd.PTS = bitn.As90k(33)
	} else {
		bitn.Forward(7)
	}
}

// SpliceInsert handles SCTE 35 splice insert commands.
func (cmd *SpCmd) SpliceInsert(bitn *bitter.Bitn) {
	cmd.Name = "Splice Insert"
	cmd.SpliceEventId = bitn.AsHex(32)
	cmd.SpliceEventCancelIndicator = bitn.AsBool()
	bitn.Forward(7)
	if !(cmd.SpliceEventCancelIndicator) {
		cmd.OutOfNetworkIndicator = bitn.AsBool()
		cmd.ProgramSpliceFlag = bitn.AsBool()
		cmd.DurationFlag = bitn.AsBool()
		cmd.SpliceImmediateFlag = bitn.AsBool()
		bitn.Forward(4)
	}
	if cmd.ProgramSpliceFlag {
		if !(cmd.SpliceImmediateFlag) {
			cmd.SpliceTime(bitn)
		}
	} else {
		cmd.ComponentCount = bitn.AsUInt64(8)
		var Components [256]uint64
		cmd.Components = Components[0:cmd.ComponentCount]
		for i := range cmd.Components {
			cmd.Components[i] = bitn.AsUInt64(8)
		}
		if !(cmd.SpliceImmediateFlag) {
			cmd.SpliceTime(bitn)
		}
	}
	if cmd.DurationFlag {
		cmd.ParseBreak(bitn)
	}
	cmd.UniqueProgramId = bitn.AsUInt64(16)
	cmd.AvailNum = bitn.AsUInt64(8)
	cmd.AvailExpected = bitn.AsUInt64(8)
}

// SpliceNull is a No-Op command.
func (cmd *SpCmd) SpliceNull() {
	cmd.Name = "Splice Null"
}

// TimeSignal splice command is a wrapper for SpliceTime.
func (cmd *SpCmd) TimeSignal(bitn *bitter.Bitn) {
	cmd.Name = "Time Signal"
	cmd.SpliceTime(bitn)
}

// BandwidthReservation splice command.
func (cmd *SpCmd) BandwidthReservation(bitn *bitter.Bitn) {
	cmd.Name = "Bandwidth Reservation"
}

// PrivateCommand splice command.
func (cmd *SpCmd) PrivateCommand(bitn *bitter.Bitn) {
	cmd.Name = "Private Command"
	cmd.Identifier = bitn.AsUInt64(32)
}

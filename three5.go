package three5

import "fmt"
import "io"
import "log"
import "os"
import "encoding/base64"
import "github.com/futzu/gobit"

const PktSz = 188
const BufferSize = 384 * PktSz

// DeB64 decodes base64 strings
func DeB64(b64 string) []byte {
	deb64, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Fatal(err)
	}
	return deb64
}

// Find is a test for array inclusion
func Find(slice []uint8, val uint8) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// PktParser is a parser for an mpegts SCTE 35 packet
func PktParser(pkt []byte) {
	if pkt[5] == 0xfc {
		if pkt[6]>>4 == 3 {
			if pkt[8] == 0 {
				cmds := []uint8{0, 5, 6, 7, 255}
				_, found := Find(cmds, pkt[18])
				if found {
					var bitn gobit.Bitn
					bitn.Load(pkt[4:PktSz])
					var spi SpInfo
					spi.Decode(&bitn)
					var cmd SpCmd
					cmd.Decode(&bitn, spi.SpliceCommandType)
					fmt.Printf("%+v\n", spi)
					fmt.Printf("%+v\n", cmd)
				}
			}
		}
	}
}

// FileParser is a parser for an mpegts file
func FileParser(fname string) {
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	buffer := make([]byte, BufferSize)
	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
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
	SectionLength          uint16
	ProtocolVersion        uint8
	EncryptedPacket        bool
	EncryptionAlgorithm    uint8
	PtsAdjustment          float64
	CwIndex                string
	Tier                   string
	SpliceCommandLength    uint16
	SpliceCommandType      uint8
}

// Decode extracts bits for the splice info section values
func (spi *SpInfo) Decode(bitn *gobit.Bitn) {
	spi.Name = "Splice Info Section"
	spi.TableId = bitn.AsHex(8)
	spi.SectionSyntaxIndicator = bitn.AsBool()
	spi.Private = bitn.AsBool()
	spi.Reserved = bitn.AsHex(2)
	spi.SectionLength = bitn.AsUInt16(12)
	spi.ProtocolVersion = bitn.AsUInt8(8)
	spi.EncryptedPacket = bitn.AsBool()
	spi.EncryptionAlgorithm = bitn.AsUInt8(6)
	spi.PtsAdjustment = bitn.As90k(33)
	spi.CwIndex = bitn.AsHex(8)
	spi.Tier = bitn.AsHex(12)
	spi.SpliceCommandLength = bitn.AsUInt16(12)
	spi.SpliceCommandType = bitn.AsUInt8(8)
}

// SpCmd is the splice command for the SCTE35 cue
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
	ComponentCount             uint8
	Components                 []uint8
	UniqueProgramId            uint16
	AvailNum                   uint8
	AvailExpected              uint8
	Identifier                 uint32
}

// Decode the splice command values
func (cmd *SpCmd) Decode(bitn *gobit.Bitn, cmdtype uint8) {
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

// ParseBreak parses out the ad break duration values
func (cmd *SpCmd) ParseBreak(bitn *gobit.Bitn) {
	cmd.BreakAutoReturn = bitn.AsBool()
	bitn.Forward(6)
	cmd.BreakDuration = bitn.As90k(33)
}

// SpliceTime parses out the PTS value as needed
func (cmd *SpCmd) SpliceTime(bitn *gobit.Bitn) {
	cmd.TimeSpecifiedFlag = bitn.AsBool()
	if cmd.TimeSpecifiedFlag {
		bitn.Forward(6)
		cmd.PTS = bitn.As90k(33)
	} else {
		bitn.Forward(7)
	}
}

// SpliceInsert handles SCTE 35 splice insert commands
func (cmd *SpCmd) SpliceInsert(bitn *gobit.Bitn) {
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
	}
	if !(cmd.ProgramSpliceFlag) {
		cmd.ComponentCount = bitn.AsUInt8(8)
		var Components [256]uint8
		cmd.Components = Components[0:cmd.ComponentCount]
		for i := range cmd.Components {
			cmd.Components[i] = bitn.AsUInt8(8)
		}
		if !(cmd.SpliceImmediateFlag) {
			cmd.SpliceTime(bitn)
		}
	}
	if cmd.DurationFlag {
		cmd.ParseBreak(bitn)
	}
	cmd.UniqueProgramId = bitn.AsUInt16(16)
	cmd.AvailNum = bitn.AsUInt8(8)
	cmd.AvailExpected = bitn.AsUInt8(8)
}

// SpliceNull is a no op command.
func (cmd *SpCmd) SpliceNull() {
	cmd.Name = "Splice Null"
}

// TimeSignal splice command is a wrapper for SpliceTime
func (cmd *SpCmd) TimeSignal(bitn *gobit.Bitn) {
	cmd.Name = "Time Signal"
	cmd.SpliceTime(bitn)
}

// BandwidthReservation splice command
func (cmd *SpCmd) BandwidthReservation(bitn *gobit.Bitn) {
	cmd.Name = "Bandwidth Reservation"
}

// PrivateCommand splice command
func (cmd *SpCmd) PrivateCommand(bitn *gobit.Bitn) {
	cmd.Name = "Private Command"
	cmd.Identifier = bitn.AsUInt32(32)
}

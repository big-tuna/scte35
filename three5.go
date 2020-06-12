package three5

import "encoding/base64"
import "github.com/futzu/gobit"
import "log"

// DeB64 decodes base64 strings
func DeB64(b64 string) []byte {
	deb64, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Fatal(err)
	}
	return deb64
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
}

// Decode extracts bits for the splice info section values
func (spi *SpInfo) Decode(bitn *gobit.Bitn) {
	spi.Name = "Splice Info Section"
	spi.TableId = bitn.AsHex(8)
	spi.SectionSyntaxIndicator = bitn.AsBool()
	spi.Private = bitn.AsBool()
	spi.Reserved = bitn.AsHex(2)
	spi.SectionLength = bitn.AsInt(12)
	spi.ProtocolVersion = bitn.AsInt(8)
	spi.EncryptedPacket = bitn.AsBool()
	spi.EncryptionAlgorithm = bitn.AsInt(6)
	spi.PtsAdjustment = bitn.As90k(33)
	spi.CwIndex = bitn.AsHex(8)
	spi.Tier = bitn.AsHex(12)
	spi.SpliceCommandLength = bitn.AsInt(12)
	spi.SpliceCommandType = bitn.AsInt(8)
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
	ComponentCount             uint64
	Components                 []uint64
	UniqueProgramId            uint64
	AvailNum                   uint64
	AvailExpected              uint64
	Identifier                 uint64
}

// Decode the splice command values
func (cmd *SpCmd) Decode(bitn *gobit.Bitn, cmdtype uint64) {
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
		cmd.ComponentCount = bitn.AsInt(8)
		var Components [100]uint64
		for i := range Components {
			Components[i] = bitn.AsInt(8)
		}
		cmd.Components = Components[0:cmd.ComponentCount]
		if !(cmd.SpliceImmediateFlag) {
			cmd.SpliceTime(bitn)
		}
	}
	if cmd.DurationFlag {
		cmd.ParseBreak(bitn)
	}
	cmd.UniqueProgramId = bitn.AsInt(16)
	cmd.AvailNum = bitn.AsInt(8)
	cmd.AvailExpected = bitn.AsInt(8)
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
	cmd.Identifier = bitn.AsInt(32)
}

package three5

import "fmt"
import "github.com/futzu/gobit"

// SpInfo is the splice info section of the SCTE 35 cue.
type SpInfo struct {
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
func (spi *SpInfo) Decode(bitn gobit.Bitn) {
	spi.TableId = bitn.AsHex(8)
	spi.SectionSyntaxIndicator = bitn.AsBool()
	spi.Private = bitn.AsBool()
	spi.Reserved = bitn.AsHex(2)
	spi.SectionLength = bitn.AsInt(12)
	spi.ProtocolVersion = bitn.AsInt(8)
	spi.EncryptedPacket = bitn.AsBool()
	spi.EncryptionAlgorithm = bitn.AsInt(6)
	spi.PtsAdjustment = bitn.AsFloat(33)
	spi.CwIndex = bitn.AsHex(8)
	spi.Tier = bitn.AsHex(12)
	spi.SpliceCommandLength = bitn.AsInt(12)
	spi.SpliceCommandType = bitn.AsInt(8)
}

// SpCmd is the splice command for the SCTE35 cue
type SpCmd struct {
	Name                       string
	BreakAutoReturn            bool
	BreakDuration              float64
	TimeSpecifiedFlag          bool
	PTS                        float64
	SpliceEventId              uint64
	SpliceEventCancelIndicator bool
	OutOfNetworkIndicator      bool
	ProgramSpliceFlag          bool
	DurationFlag               bool
	SpliceImmediateFlag        bool
	ComponentCount             uint64
	SpliceImmediateFlag        bool
	UniqueProgramId            uint64
	AvailNum                   uint64
	AvailExpected              uint64
}

// Decode the splice command values
func (cmd *SpCmd) Decode(bitn gobit.Bitn, cmdtype uint64) {
	if cmdtype == 0 {
		cmd.SpliceNull()
	}
	//4: Splice_Schedule,
	//5: Splice_Insert,
	if cmdtype == 6 {
		cmd.TimeSignal(bitn)
	}
	// 7: Bandwidth_Reservation
	// 255: Private_Command
}

// ParseBreak parses out the ad break duration values
func (cmd *SpCmd) ParseBreak(bitn gobit.Bitn) {
	cmd.BreakAutoReturn = bitn.AsBool()
	bitn.Forward(6)
	cmd.BreakDuration = bitn.AsFloat(33)
}

// SpliceTime parses out the PTS value as needed
func (cmd *SpCmd) SpliceTime(bitn gobit.Bitn) {
	cmd.TimeSpecifiedFlag = bitn.AsBool()
	if cmd.TimeSpecifiedFlag == true {
		bitn.Forward(6)
		cmd.PTS = bitn.AsFloat(33)
	} else {
		bitn.Forward(7)
	}
}
func (cmd *SpCmd) SpliceInsert(bitn gobit.Bitn) {
	cmd.Name = "Splice Insert"
	cmd.SpliceEventId = bitbin.AsInt(32)
	cmd.SpliceEventCancelIndicator = bitbin.AsBool()
	bitbin.Forward(7)
	if cmd.SpliceEventCancelIndicator == false {
		cmd.OutOfNetworkIndicator = bitbin.AsBool()
		cmd.ProgramSpliceFlag = bitbin.AsBool()
		cmd.DurationFlag = bitbin.AsBool()
		cmd.SpliceImmediateFlag = bitbin.AsBool()
		bitbin.Forward(4)
		if cmd.ProgramSpliceFlag == true {
			if cmd.SpliceImmediateFlag == false {
				cmd.SpliceTime(bitbin)
			}
		}
		if cmd.ProgramSpliceFlag == false {
			cmd.ComponentCount = bitbin.AsInt(8)
			/**
			                cmd.Components = []
			                for i in range(0, cmd.ComponentCount):
			                    cmd.Components[i] = bitbin.AsInt(8)
						**/
			if cmd.SpliceImmediateFlag == false {
				cmd.SpliceTime(bitbin)
			}
		}
		if cmd.DurationFlag == true {
			cmd.ParseBreak(bitbin)
		}
		cmd.UniqueProgramId = bitbin.AsInt(16)
		cmd.AvailNum = bitbin.AsInt(8)
		cmd.AvailExpected = bitbin.AsInt(8)
	}
}

// SpliceNull is a no op command.
func (cmd *SpCmd) SpliceNull() {
	cmd.Name = "Splice Null"
}

// TimeSignal is a wrapper for Splicetime
func (cmd *SpCmd) TimeSignal(bitn gobit.Bitn) {
	cmd.Name = "Time Signal"
	cmd.SpliceTime(bitn)
}

package three5

import "fmt"
import "github.com/futzu/gobit"

// SpCmd is the splice command for the SCTE35 cue.
type SpCmd struct {
	Name              string
	bitn              gobit.Bitn
	BreakAutoReturn   bool
	BreakDuration     float64
	TimeSpecifiedFlag bool
	PTS               float64
}

// Load bites into cmd.bitn instance
func (cmd *SpCmd) Load(bites []byte) {
	cmd.bitn.Load(bites)
}

// ParseBreak parses out the break duration
func (cmd *SpCmd) ParseBreak() {
	cmd.BreakAutoReturn = cmd.bitn.AsBool()
	cmd.bitn.Forward(6)
	cmd.BreakDuration = cmd.bitn.AsFloat(33)
}

// SpliceTime parses out the PTS value if cmd.timeSpecifiedFlag
func (cmd *SpCmd) SpliceTime() {
	cmd.TimeSpecifiedFlag = cmd.bitbin.AsBool()
	if cmd.TimeSpecifiedFlag == true {
		cmd.bitbin.Forward(6)
		cmd.PTS = cmd.bitbin.AsFloat(33)
	} else {
		cmd.bitbin.Forward(7)
	}
}

// Splice Null Command
func (cmd *SpCmd) SpliceNull() {
	cmd.Name = "Splice Null"
}

// TimeSignal is a wrapper for Splicetime
func (cmd *SpCmd) TimeSignal() {
	cmd.Name = "Time Signal"
	cmd.SpliceTime()
}

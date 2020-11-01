package scte35

import "fmt"
import "os"
import "encoding/base64"
import "encoding/json"
import "github.com/futzu/bitter"

// PktSz is the size of an MPEG-TS packet in bytes.
const PktSz = 188

// BufferSize is the size of a read when parsing files.
const BufferSize = 384 * PktSz

// Generic catchall error checking
func Chk(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

// MkJson structs to JSON
func MkJson(i interface{}) string {
	jason, err := json.MarshalIndent(&i, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	return string(jason)
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

// PktParser is a parser for an MPEG-TS SCTE 35 packet
func PktParser(pkt []byte) {
	var hdr bitter.Bitn
	hdr.Load(pkt[0:5])
	hdr.Forward(8)
	hdr.Forward(3)
	PID := hdr.AsUInt64(12)
	pld := pkt[5:PktSz]
	magicbytes := [4]uint8{252, 48, 0, 255}
	pldbytes := [4]uint8{pld[0], pld[1], pld[3], pld[10]}
	if pldbytes == magicbytes {
		cmds := []uint8{0, 5, 6, 7, 255}
		if IsIn(cmds, pld[13]) {
			var cue Cue
			cue.Decode(pld)
			fmt.Println(PID)

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

type Cue struct {
	InfoSection SpInfo
	Command     SpCmd
	Descriptors []SpDscptr `json:",omitempty"`
}

// Decode extracts bits for the Cue values.
func (cue *Cue) Decode(bites []byte) {
	var bitn bitter.Bitn
	bitn.Load(bites)
	cue.InfoSection.Decode(&bitn)
	cue.Command.Decode(&bitn, cue.InfoSection.SpliceCommandType)
	cue.InfoSection.DescriptorLoopLength = bitn.AsUInt64(16)
	cue.DscptrLoop(&bitn)
	fmt.Println(MkJson(&cue))
}

func (cue *Cue) DscptrLoop(bitn *bitter.Bitn) {
	var i uint64
	i = 0
	for i < cue.InfoSection.DescriptorLoopLength {
		var sd SpDscptr
		sd.DescriptorType = bitn.AsUInt64(8)
		sd.DescriptorLen = bitn.AsUInt64(8)
		//sd.Decode(bitn *bitter.Bitn)
		i += sd.DescriptorLen + 2
		cue.Descriptors = append(cue.Descriptors, sd)
	}
}

/**
  SetSpDscptr(self):
        '''

        threefive.Cue.set_splice_descriptor
        is called by threefive.Cue.descriptorloop.
        '''
        # splice_descriptor_tag 8 uimsbf
        tag = self.bitbin.asint(8)
        desc_len = self.bitbin.asint(8)
        if tag in self.sd_tags:
            if tag == 2:
                sd = SegmentationDescriptor()
            else:
                sd = SpliceDescriptor()
            sd.parse(self.bitbin,tag)
            sd.descriptor_length = desc_len
            return sd
        else:
            return False
**/

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
	SpliceEventId              string   `json:"omitempty"`
	SpliceEventCancelIndicator bool     `json:"omitempty"`
	OutOfNetworkIndicator      bool     `json:"omitempty"`
	ProgramSpliceFlag          bool     `json:"omitempty"`
	DurationFlag               bool     `json:"omitempty"`
	BreakAutoReturn            bool     `json:"omitempty"`
	BreakDuration              float64  `json:"omitempty"`
	SpliceImmediateFlag        bool     `json:"omitempty"`
	TimeSpecifiedFlag          bool     `json:"omitempty"`
	PTS                        float64  `json:"omitempty"`
	ComponentCount             uint64   `json:"omitempty"`
	Components                 []uint64 `json:"omitempty"`
	UniqueProgramId            uint64   `json:"omitempty"`
	AvailNum                   uint64   `json:"omitempty"`
	AvailExpected              uint64   `json:"omitempty"`
	Identifier                 uint64   `json:"omitempty"`
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

type AudioCmpnt struct {
	component_tag   string `json:"omitempty"`
	ISO_code        uint64 `json:"omitempty"`
	bit_stream_mode uint64 `json:"omitempty"`
	num_channels    uint64 `json:"omitempty"`
	full_srvc_audio bool   `json:"omitempty"`
}

// Splice Descriptor
type SpDscptr struct {
	DescriptorType uint64
	DescriptorLen  uint64
	// identiﬁer 32 uimsbf == 0x43554549 (ASCII “CUEI”)
	Identifier      string
	Name            string
	ProviderAvailId uint64 `json:"omitempty"`
	PreRoll         uint64 `json:"omitempty"`
	DTMFCount       uint64 `json:"omitempty"`
	//dtmf_chars = [] `json:"omitempty"`
	TAISeconds uint64 `json:"omitempty"`
	TAINano    uint64 `json:"omitempty"`
	UTCOffset  uint64 `json:"omitempty"`
}

// AvailDscptr Avail Splice Descriptor
func (dscptr *SpDscptr) AvailDscptr(bitn *bitter.Bitn) {
	dscptr.Name = "Avail Descriptor"
	dscptr.ProviderAvailId = bitn.AsUInt64(32)
}

// DTMFDscptr DTMF Splice DSescriptor
func (dscptr *SpDscptr) DTMFDscptr(bitn *bitter.Bitn) {
	dscptr.Name = "DTMF Descriptor"
	dscptr.PreRoll = bitn.AsUInt64(8)
	dscptr.DTMFCount = bitn.AsUInt64(3)
	bitn.Forward(5)
	/**
	        dscptr.DTMFChars = []
	        for i in range(0, dscptr.DTMFCount):
	            dscptr.DTMFChars.append(bitbin.asint(8))
			**/
}

// TimeDscptr Time Splice DSescriptor
func (dscptr *SpDscptr) TimeDscptr(bitn *bitter.Bitn) {
	dscptr.Name = "Time Descriptor"
	dscptr.TAISeconds = bitn.AsUInt64(48)
	dscptr.TAINano = bitn.AsUInt64(32)
	dscptr.UTCOffset = bitn.AsUInt64(16)
}

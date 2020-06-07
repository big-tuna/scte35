package three5
import "fmt"
import "github.com/futzu/gobit"

func main() {
	bites := []byte("\xfc0H\x00\x00\x00\x00\x00\x00\xff\xff\xf0\x05\x06\xfe\x93.8\x0b\x002\x02\x17CUEIH\x00\x00\n\x7f\x9f\x08\x08\x00\x00\x00\x00,\xa0\xa1\xe3\x18\x00\x00\x02\x17CUEIH\x00\x00\t\x7f\x9f\x08\x08\x00\x00\x00\x00,\xa0\xa1\x8a\x11\x00\x00\xb4!~\xb0")
	var bitn gobit.Bitn
	bitn.Load(bites)
	tableid := bitn.AsHex(8)
	section_syntax_indicator := bitn.AsBool()
	private := bitn.AsBool()
	reserved := bitn.AsHex(2)
	section_length := bitn.AsInt(12)
	protocol_version := bitn.AsInt(8)
	encrypted_packet := bitn.AsBool()
	encryption_algorithm := bitn.AsInt(6)
	pts_adjustment := bitn.AsFloat(33)
	cw_index := bitn.AsHex(8)
	tier := bitn.AsHex(12)
	splice_command_length := bitn.AsInt(12)
	splice_command_type := bitn.AsInt(8)

	fmt.Println(tableid, section_syntax_indicator, private,reserved,
			section_length, protocol_version, encrypted_packet,
			encryption_algorithm, pts_adjustment, cw_index, tier,
			splice_command_length,  splice_command_type)



}

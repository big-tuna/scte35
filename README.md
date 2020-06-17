# three5
The SCTE 35 Parser in Go.
### Heads up, work in progress.

#### What is working:
✓ Parsing MPEG-TS files

✓ Parsing Base64 strings

✓ Splice Info Section 	
	
✓ Splice Null
	
✓ Time Signal
	
✓ Splice Insert   	 	
	
✓ Bandwidth Reservation
	
✓ Private Command		

#### Installation
```sh
go get -u github.com/futzu/three5
```

#### Parsing an MPEG-TS file 
```go
package main

import (
	"github.com/futzu/three5"
)

func main(){
    fname := "video.ts" 
    three5.FileParser(fname)
}   
```


#### Parsing a base64 string
```go
package main

import "github.com/futzu/three5"

func main() {
	b64 := "/DAvAAAAAAAA///wFAVIAACPf+/+c2nALv4AUsz1AAAAAAAKAAhDVUVJAAABNWLbowo="
	bites := three5.DeB64(b64)
	three5.SCTE35Parser(bites)
	}
```  
---
##### Output
```sh
{Name:Splice Info Section TableId:0xfc SectionSyntaxIndicator:false
Private:false Reserved:0x3 SectionLength:47 ProtocolVersion:0
EncryptedPacket:false EncryptionAlgorithm:0 PtsAdjustment:0 CwIndex:0xff 
Tier:0xfff SpliceCommandLength:20 SpliceCommandType:5}
{Name:Splice Insert SpliceEventId:0x4800008f SpliceEventCancelIndicator:false 
OutOfNetworkIndicator:true ProgramSpliceFlag:true DurationFlag:true BreakAutoReturn:true 
BreakDuration:60.293566 SpliceImmediateFlag:false TimeSpecifiedFlag:true PTS:21514.559088 
ComponentCount:0 Components:[] UniqueProgramId:0 AvailNum:0 AvailExpected:0 Identifier:0}

```


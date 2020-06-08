# three5
The SCTE 35 Parser in Go.
## Heads up, this is not yet functional. 
### I expect to have it mostly working within the next week or so. 


#### What is working:

 * Splice Info Section ✓
 * Splice Commands 
    * Splice Null     ✓
    * Time Signal     ✓

#### Here's how I'm running it.
```go
package main

import "github.com/futzu/three5"
import "fmt"
import "github.com/futzu/gobit"

func main() {
	bites := []byte("\xfc0H\x00\x00\x00\x00\x00\x00\xff\xff\xf0\x05\x06\xfe\x93.8\x0b\x002\x02\x17CUEIH\x00\x00\n\x7f\x9f\x08\x08\x00\x00\x00\x00,\xa0\xa1\xe3\x18\x00\x00\x02\x17CUEIH\x00\x00\t\x7f\x9f\x08\x08\x00\x00\x00\x00,\xa0\xa1\x8a\x11\x00\x00\xb4!~\xb0")
	var bitn gobit.Bitn
	bitn.Load(bites)
	var spi three5.SpInfo
	spi.Decode(bitn)
	var cmd three5.SpCmd
	cmd.Decode(bitn,spi.SpliceCommandType)
	fmt.Println(spi)
	fmt.Println(cmd)
	}
```  

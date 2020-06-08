# three5
The SCTE 35 Parser in Go.
## Heads up, this is not yet functional. 

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
	bites := []byte("\xfc0/\x00\x00\x00\x00\x00\x00\xff\xff\xf0\x14\x05H\x00\x00\x8f\x7f\xef\xfesi\xc0.\xfe\x00R\xcc\xf5\x00\x00\x00\x00\x00\n\x00\x08CUEI\x00\x00\x015b\xdb\xa3")
	var bitn gobit.Bitn
	bitn.Load(bites)
	var spi three5.SpInfo
	spi.Decode(&bitn)
	var cmd three5.SpCmd
	cmd.Decode(&bitn,spi.SpliceCommandType)
	fmt.Printf("%+v",spi)
	fmt.Printf("%+v",cmd)
	}
	
```  

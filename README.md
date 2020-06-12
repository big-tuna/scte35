# three5
The SCTE 35 Parser in Go.
## Heads up, this is not yet functional. 

#### What is working:
* Base64 encoded strings ✓
 * Splice Info Section 	 ✓
 * Splice Commands 
    * Splice Null     	 ✓
    * Time Signal     	 ✓
    * Splice Insert   	 ✓
    
#### Here's how I'm running it.
```go
package main

import "github.com/futzu/three5"
import "fmt"
import "github.com/futzu/gobit"

func main() {
	b64 := "/DAvAAAAAAAA///wFAVIAACPf+/+c2nALv4AUsz1AAAAAAAKAAhDVUVJAAABNWLbowo="
	bites := three5.DeB64(b64)
	var bitn gobit.Bitn
	bitn.Load(bites)
	var spi three5.SpInfo
	spi.Decode(&bitn)
	var cmd three5.SpCmd
	cmd.Decode(&bitn,spi.SpliceCommandType)
	fmt.Printf("%+v\n",spi)
	fmt.Printf("%+v\n",cmd)
	}
```  

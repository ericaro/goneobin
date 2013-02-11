/*neobin package implements a code generator that generates Reader and Writer object for a given neobin fileformat.
*/
package neobin

import (
"encoding/xml"
"io/ioutil"
"path/filepath"
"os"
)

// this is the main function
// it basically unmarshall the src, and generate to a target file

func Generate(srcfile, targetfile string) (err error) {
	data , err := ioutil.ReadFile( srcfile )
	pack := filepath.Base(  filepath.Dir(targetfile)) // this should be the package name
	
	if err != nil {
		return
	}
	n := &neobin{}
	n.GoPackage = pack
	err = xml.Unmarshal([]byte(data), n)
	if err != nil {
		return
	}
	
	m, err := compile(n)
	if err != nil {
		return
	}
	
	fid, err := os.Create( targetfile )
	
	err  = generateWriter(fid, m)
	return

}

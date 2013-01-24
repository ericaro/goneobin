package neobin

import (
"encoding/xml"
"io/ioutil"
"path/filepath"
"os"
)

func Generate(srcfile, targetfile string) (err error) {
	data , err := ioutil.ReadFile( srcfile )
	pack := filepath.Base(  filepath.Dir(targetfile)) // this should be the package name
	
	if err != nil {
		return
	}
	n := &neobin{}
	n.GoPackage = pack
	err = xml.Unmarshal([]byte(data), &n)
	if err != nil {
		return
	}
	fid, err := os.Create( targetfile )
	err  = generate(fid, n)
	return

}

package main

import (
	"nbw"
	"os"
)

func main() {
	w, err := os.Create("test.bin")
	if err != nil {
		panic(err.Error())
	}
	p := nbw.NewSimulationWriter(w, 1024)
	
	t := p.Name("toto").Timestamp(0).Ids([]int64{1,2}).Values([]float32{3.,4.})
	
	var i int64
	for ;i<100;i++{
		t = t.Timestamp(i+1).Ids([]int64{1,2}).Values([]float32{3.,4.})
	}
	t.Eof()
	w.Close()
	
}

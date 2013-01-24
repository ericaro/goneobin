package main

import (
	"ericaro.net/neobin"
	"flag"
)

var (
	

)

func main() {
	flag.Parse()
	src := flag.Arg(0)
	target := flag.Arg(1)
	neobin.Generate(src, target)

}

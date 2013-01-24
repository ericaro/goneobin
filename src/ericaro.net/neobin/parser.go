package neobin

import (
	"ericaro.net/gogrex"
)

type dataType string

var (
	byteType   = dataType("byte")   // byte       byte
	intType    = dataType("int")    // int         int
	longType   = dataType("long")   // long       long
	stringType = dataType("string") // string   string
	floatType  = dataType("float")  // float     float
	doubleType = dataType("double") // double   double

	manyMany   = "many"
	manySingle = "single"
)

// struct dedicated to parsing neobin xml
type neobin struct {
	Package     string  `xml:"package"`
	GoPackage   string  `xml:"-"` // the generated go package name
	Header      *string `xml:"header"`
	Name        *string `xml:"name"` //optional
	Expression  string  `xml:"expression"`
	grex        *gogrex.Grex
	Transitions []*transition `xml:"transitions>transition"`
	States      []*state      `xml:"states>state"`
}

type transition struct {
	Name      string      `xml:"name,attr"`
	Variables []*variable `xml:"var"`
}

type variable struct {
	Type dataType `xml:"type,attr"`
	Name string   `xml:"name,attr"`
	Many string  `xml:"many,attr"`
}

type state struct {
	id     int    `xml:"-"`
	Output bool   `xml:"-"`
	Input  bool   `xml:"-"`
	Choice bool   `xml:"-"`
	Path   string `xml:"path,attr"`
	State  string `xml:",chardata"`
}


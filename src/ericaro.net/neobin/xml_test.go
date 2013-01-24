package neobin

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func assertNotNil(t *testing.T, p interface{}, nature string) {

	if p == nil {
	t.Fatalf("%s: unexpected nil\n", nature)
	}
}

func assertString(t *testing.T, str, expected, nature string){
	
	fmt.Printf("%s: %s\n", nature, str)
	if str != expected {
	t.Fatalf("%s: \"%s\" !=   \"%s\"\n", nature, str, expected)
	}
	
}

func TestXYZ(t *testing.T) {

	data := `<?xml version="1.0" encoding="UTF-8"?>
<neobin xmlns="http://neobin.ericaro.net/v1/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
        <expression>name, (id, value)+, eof</expression>
        <header>toto</header>
        <name>Demo</name>
        <package>net.ericaro.demo</package>
        <transitions>

                <transition name="name">
                        <var type="string" name="name"/>
                </transition>
                <transition name="id">
                        <var type="int" name="id"/>
                </transition>
                <transition name="value">
                        <var type="float" name="x"/>
                        <var type="float" name="y"/>
                        <var type="float" name="z"/>
                </transition>
                <transition name="eof">
                </transition>
        </transitions>

        <states>
        <state path="">Perf</state>
        <state path="name">Named</state>
        <state path="name.id">Idified</state>
        <state path="name.id.value">Valued</state>
        </states>

</neobin>
`
	n := neobin{}
	err := xml.Unmarshal([]byte(data), &n)
	if err != nil {
		t.Fatalf("error: %v", err)
		return
	}
	assertString(t, n.Package, "net.ericaro.demo", "Package")
	
	assertNotNil(t, n.Header, "Header")
	assertString(t, *n.Header, "toto", "Header")
	
	assertNotNil(t, n.Name, "Name")
	assertString(t, *n.Name, "Demo", "Name")
	
	assertString(t, n.Expression, "name, (id, value)+, eof", "Expression")
	
	if len(n.Transitions)  != 4 {
	t.Fatalf("Was expecting 4 transitions, got %d", len(n.Transitions)  )
	}
	
	// loop over the 4 transitions
	var tr *transition
	
	tr = n.Transitions[0]
	fmt.Printf("transition\n")
	assertString(t, tr.Name, "name", "    Name")
	if len(tr.Variables)  != 1 {
		t.Fatalf("Was expecting 1 var, got %d ", len(tr.Variables)  )
	}
	
	var v *variable
	v = tr.Variables[0]
	
	assertString(t, string(v.Type), "string", "        Type")
	assertString(t, v.Name, "name", "        Name")
	if v.Many != nil {
		t.Fatalf("Was a nil Many attr")
	}
	
	
	if len(n.States)  != 4 {
	t.Fatalf("Was expecting 4 States, got %d", len(n.States)  )
	}
	var s *state
	s = n.States[0]
	assertString(t, s.Path, "", "        Path")
	assertString(t, s.State, "Perf", "        State")
	
	
	
	generate(&n)
	
	fmt.Printf("Success\n")
	
	t.Fail()

}

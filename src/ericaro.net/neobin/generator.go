package neobin

import (
	"ericaro.net/gogrex"
	"io"
	"sort"
	"strings"
	"text/template"
)

func UpperFirst(str string) string {
	return strings.ToUpper(str[0:1]) + str[1:]
}

var (
	funcMap = template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"uf": UpperFirst,
	}
	transitionTemplate = template.Must(template.New("transition").Funcs(funcMap).Parse(
`//
func (this *{{.Source.Name| uf}}) {{.Name | uf}}({{range $v := .Variables}}{{$v.Name}} {{if .IsMany}}[]{{end}}{{$v.TypeDef}}, {{end}}){{if not .Dest.Output}}  *{{.Dest.Name| uf}} {{end}}{
	{{if .Source.Choice}}this.encodeVar({{.Id}});{{end}}
	{{range $v := .Variables}}this.write{{if .IsMany}}Many{{end}}{{$v.EncodeDef | uf}}({{$v.Name}})
{{end}}
{{if .Dest.Output}}this.Flush(){{else}}return &{{.Dest.Name| uf}}{this.binWriter}{{end}}
}
	`))
	
	
	
	stateTemplate = template.Must(template.New("state").Funcs(funcMap).Parse(
`{{if not .Output}}//##################################################################################################
//{{.Name}}
type {{.Name | uf}} struct {
	*binWriter
}
{{end}}`))

	headTemplate = template.Must(template.New("head").Funcs(funcMap).Parse(
`package {{.GoPackage}}
	
import (
	"io"
	"math"
)

const (
	MAGIC =  0x14E0B114
)
	
{{with .InputState}}func New{{.Name}}Writer(w io.Writer, bufsize int) *{{.Name}} {
	state := &{{.Name}}{
		binWriter : &binWriter{
			w: w,
			buf: make([]byte, 0, bufsize),
		} ,
	}{{end}}
	state.writeInt32(MAGIC)
	state.writeString("{{.PackageName}}")
	return state
}
`))

	binWriterTemplate = template.Must(template.New("binwriter").Funcs(funcMap).Parse(
`//##################################################################################################
//   binary writers
//##################################################################################################
type binWriter struct {
	w io.Writer
	buf []byte
}
	
func (p *binWriter) append(b ...uint8) {
	p.buf = append(p.buf, b...)
}
func (p *binWriter) Flush() (err error) {
	n, err := p.w.Write(p.buf)
	p.buf = p.buf[n:]
	return
}

func (p *binWriter) encodeLen(x int) error {
	return p.encodeVar(uint64(x) )
}
func (p *binWriter) encodeVar(x uint64) error {
	for x >= 1<<7 {
		p.append(uint8(x&0x7f | 0x80))
		x >>= 7
	}
	p.append(uint8(x))
	return nil
}

func (p *binWriter) encodeFixed32(x uint32) error {
	p.append(
		uint8(x),
		uint8(x>>8),
		uint8(x>>16),
		uint8(x>>24))
	return nil
}
func (p *binWriter) encodeFixed64(x uint64) error {
	p.append(
		uint8(x),
		uint8(x>>8),
		uint8(x>>16),
		uint8(x>>24),
		uint8(x>>32),
		uint8(x>>40),
		uint8(x>>48),
		uint8(x>>56))
	return nil
}
func (p *binWriter) writeByte(x byte) error {
	p.append(x)
	return nil
}
func (p *binWriter) writeString(s string) error {
	p.encodeVar(uint64(len(s)))
	p.append([]byte(s)...)
	return nil
}
func (p *binWriter) writeInt32(x int32) error {
	return p.encodeFixed32( uint32(x) ) 
}
func (p *binWriter) writeInt64(x int64) error {
	return p.encodeFixed64( uint64(x) ) 
}
func (p *binWriter) writeFloat32(x float32) error {
	return p.encodeFixed32( math.Float32bits(x) )
}
func (p *binWriter) writeFloat64(x float64) error {
	return p.encodeFixed64( math.Float64bits(x) )
}

func (p *binWriter) writeManyByte(x []byte) error {
	p.encodeLen(len(x) )
	for _,b := range x {
		err := p.writeByte(b)
		if err != nil {return err}
	}
	return nil
}
func (p *binWriter) writeManyString(x []string) error {
	p.encodeLen(len(x) )
	for _,b := range x {
		err := p.writeString(b)
		if err != nil {return err}
	}
	return nil
}
func (p *binWriter) writeManyInt32(x []int32) error {
	p.encodeLen(len(x) )
	for _,b := range x {
		err := p.writeInt32(b)
		if err != nil {return err}
	}
	return nil
}
func (p *binWriter) writeManyInt64(x []int64) error {
	p.encodeLen(len(x) )
	for _,b := range x {
		err := p.writeInt64(b)
		if err != nil {return err}
	}
	return nil
}
func (p *binWriter) writeManyFloat32(x []float32) error {
	p.encodeLen(len(x) )
	for _,b := range x {
		err := p.writeFloat32(b)
		if err != nil {return err}
	}
	return nil
}
func (p *binWriter) writeManyFloat64(x []float64) error {
	p.encodeLen(len(x) )
	for _,b := range x {
		err := p.writeFloat64(b)
		if err != nil {return err}
	}
	return nil
}

`))
)

type neobinManager struct {
	id     int
	neobin *neobin
}

type transitionEdge struct {
	transition   *transition
	Id           int
	source, dest *state
}

func (t *transitionEdge) Name() string {
	return t.transition.Name
}
func (t *transitionEdge) Variables() []*variable {
	return t.transition.Variables
}
func (t *state) Name() string {
	return t.State
}

func (t *variable) IsMany() bool {
	return t.Many == manyMany
}
func (t *variable) TypeDef() string {
	switch t.Type {
	case byteType:
		return "byte" // byte       byte
	case intType:
		return "int32" // int         int
	case longType:
		return "int64" // long       long
	case stringType:
		return "string" // string   string
	case floatType:
		return "float32" // float     float
	case doubleType:
		return "float64" // double   double
	}
	panic("unexpected type " + t.Type)
}
func (t *variable) EncodeDef() string {
	switch t.Type {
	case byteType:
		return "byte" // byte       byte
	case intType:
		return "int32" // int         int
	case longType:
		return "int64" // long       long
	case stringType:
		return "string" // string   string
	case floatType:
		return "float32" // float     float
	case doubleType:
		return "float64" // double   double
	}
	panic("unexpected type " + t.Type)
}

func (t *transitionEdge) Source() *state {
	return t.source
}
func (t *transitionEdge) Dest() *state {
	return t.dest
}

func (m *neobinManager) getTransition(name string) *transition {
	for _, tr := range m.neobin.Transitions {
		if tr.Name == name {
			return tr
		}
	}
	panic("unknown transition " + name)
}

func (m *neobinManager) InputState() *state {
	in := m.neobin.grex.InputVertex()
	return (in).(*state)
}
func (m *neobinManager) GoPackage() string {
	return m.neobin.GoPackage
}
func (m *neobinManager) PackageName() string {
	return m.neobin.Package
}
func (m *neobinManager) StateByPath(path string) *state {
	return (m.neobin.grex.VertexByPath(path)).(*state)
}

func (m *neobinManager) NewVertex() gogrex.Vertex {
	m.id++
	name := "State"
	if m.neobin.Name != nil {
		name = *m.neobin.Name
	}
	return &state{
		id:    m.id,
		State: name,
	}
}
func (m *neobinManager) NewEdge(name string) gogrex.Edge {
	m.id++
	tr := m.getTransition(name)
	return &transitionEdge{Id: m.id, transition: tr}
}
func (m *neobinManager) CloneEdge(e gogrex.Edge) gogrex.Edge {
	t := e.(*transitionEdge)
	m.id++
	return &transitionEdge{Id: m.id, transition: t.transition}
}

type ByName struct {
	list []gogrex.Edge
}

func (s ByName) Less(i, j int) bool {
	ti := s.list[i].(*transitionEdge)
	tj := s.list[j].(*transitionEdge)
	return ti.transition.Name < tj.transition.Name
}
func (s ByName) Len() int      { return len(s.list) }
func (s ByName) Swap(i, j int) { s.list[i], s.list[j] = s.list[j], s.list[i] }

func generate(w io.Writer, n *neobin) (err error) {


	// parse the expression using grex
	m := &neobinManager{neobin: n}
	n.grex, err = gogrex.ParseGrex(m, n.Expression)

	// build state name // start with the input one
	name := "State"
	if m.neobin.Name != nil {
		name = *m.neobin.Name
	}

	in := m.InputState()
	in.State = name
	in.Input = true

	// for every one declared in the states section:
	for _, s := range n.States {
		vertex := m.StateByPath(s.Path)
		vertex.Path = s.Path
		vertex.State = s.State
	}

	// copy src and dest into the transitionEdge
	for edge, bounds := range n.grex.Edges() {
		t := edge.(*transitionEdge)
		t.source = bounds.Start().(*state)
		t.dest = bounds.End().(*state)
	}

	// mark outputs as so
	for _, s := range n.grex.OutputVertice() {
		st := s.(*state)
		st.Output = true
	}

	err = headTemplate.Execute(w, m)
	if err != nil {
		panic(err)
	}
	

	// parse all state, then all transitions
	for _, st := range n.grex.Vertices() {
		s := st.(*state)
		outEdges := n.grex.OutputEdges(st)

		//sort outEdges first
		sort.Sort(ByName{outEdges})

		//mark the boolean choice
		s.Choice = len(outEdges) > 1
		// run the state template
		err = stateTemplate.Execute(w, s)
		if err != nil {
			panic(err)
		}

		var id int // the unique id for the transition
		for _, edge := range outEdges {
			t := edge.(*transitionEdge)
			t.Id = id
			id++
			err = transitionTemplate.Execute(w, edge)
			if err != nil {
				panic(err)
			}
		}
	}
	
	err = binWriterTemplate.Execute(w, nil)
	if err != nil {
		panic(err)
	}
	
	return
}

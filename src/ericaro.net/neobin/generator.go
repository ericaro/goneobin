package neobin

import (
	"ericaro.net/gogrex"
	"fmt"
	"io"
	"sort"
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
func (t *state) Id() string {
	return fmt.Sprintf("%d", t.id)
}
func (t *state) Single() *transitionEdge {
	if len(t.Transitions) == 1 {
		return t.Transitions[0]
	}
	return nil
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
func (m *neobinManager) States() []*state {
	return m.neobin.States
}
func (m *neobinManager) StateByPath(path string) *state {
	return (m.neobin.grex.VertexByPath(path)).(*state)
}

func (m *neobinManager) UniqueTransitions() []*transition {
		return m.neobin.Transitions
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

//computes all that is needed for this neobin
func compile(n *neobin) (m *neobinManager, err error) {

	// parse the expression using grex
	m = &neobinManager{neobin: n}
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

	v := n.grex.Vertices()
	m.neobin.States = make([]*state, len(v))

	for i, st := range v {
		s := st.(*state)
		m.neobin.States[i] = s
	}

	// parse all state, then all transitions
	for _, s := range m.States() {
		//s := st.(*state)
		outEdges := n.grex.OutputEdges(s)

		//sort outEdges first
		sort.Sort(ByName{outEdges})

		s.Transitions = make([]*transitionEdge, len(outEdges))
		for i, v := range outEdges {
			t := v.(*transitionEdge)
			s.Transitions[i] = t
		}

		//mark the boolean choice
		s.Choice = len(outEdges) > 1
		// run the state template
		if err != nil {
			panic(err)
		}

		var id int // the unique id for the transition
		for _, t := range s.Transitions {
			t.Id = id
			id++
			if err != nil {
				panic(err)
			}
		}
	}

	return
}

func generateWriter(w io.Writer, m *neobinManager) (err error) {
	
	err = FileUnitTemplate.Execute(w, m)
	if err != nil {
		panic(err)
	}
	//appends Protocol
	err = ProtocolTemplate.Execute(w, nil)
	if err != nil {
		panic(err)
	}
	return
}

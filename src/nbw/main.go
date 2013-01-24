package nbw
	
import (
	"io"
	"math"
)

const (
	MAGIC =  0x14E0B114
)
	
func NewSimulationWriter(w io.Writer, bufsize int) *Simulation {
	state := &Simulation{
		binWriter : &binWriter{
			w: w,
			buf: make([]byte, 0, bufsize),
		} ,
	}
	state.writeInt32(MAGIC)
	state.writeString("net.ericaro.demo")
	return state
}
//##################################################################################################
//Ids
type Ids struct {
	*binWriter
}
//
func (this *Ids) Ids(id []int64, )  *Values {
	
	this.writeManyInt64(id)

return &Values{this.binWriter}
}
	//##################################################################################################
//Loop
type Loop struct {
	*binWriter
}
//
func (this *Loop) Eof(){
	this.encodeVar(0);
	
this.Flush()
}
	//
func (this *Loop) Timestamp(t int64, )  *Ids {
	this.encodeVar(1);
	this.writeInt64(t)

return &Ids{this.binWriter}
}
	//##################################################################################################
//Data
type Data struct {
	*binWriter
}
//
func (this *Data) Timestamp(t int64, )  *Ids {
	
	this.writeInt64(t)

return &Ids{this.binWriter}
}
	//##################################################################################################
//Simulation
type Simulation struct {
	*binWriter
}
//
func (this *Simulation) Name(name string, )  *Data {
	
	this.writeString(name)

return &Data{this.binWriter}
}
	//##################################################################################################
//Values
type Values struct {
	*binWriter
}
//
func (this *Values) Values(x []float32, )  *Loop {
	
	this.writeManyFloat32(x)

return &Loop{this.binWriter}
}
	//##################################################################################################
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


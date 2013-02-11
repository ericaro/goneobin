package main

import (
	"nbr"
	"os"
	"fmt"
)


type R struct {
}


func(r *R) 	Name(name string, address string) (err error){fmt.Printf("receiving Name(name %s %s\n",name, address)   ;return}
func(r *R)	Timestamp(t int64) (err error)               {fmt.Printf("receiving Timestamp %v\n"	,t				)   ;return}
func(r *R)	Ids(id []int64) (err error)                  {fmt.Printf("receiving Ids(id %v\n",id	                )   ;return}
func(r *R)	Values(x []float32) (err error)              {fmt.Printf("receiving Values %v \n",x                 )   ;return}
func(r *R)	Eof() (err error)                            {fmt.Printf("receiving Eof\n"			                )   ;return}



func main() {
	w, err := os.Create("test.bin")
	if err != nil {
		panic(err.Error())
	}
	p := nbr.NewSimulationWriter(w)

	t := p.Name("toto", "toto@rigolo.com").Timestamp(0).Ids([]int64{1, 2}).Values([]float32{3., 4.})

	var i int64
	for ; i < 10; i++ {
		f := float32(i)
		t = t.Timestamp(i + 1).Ids([]int64{i, i+1}).Values([]float32{1.0*f, f*0.1})
	}
	t.Eof()
	w.Close()

	f, err := os.Open("test.bin")
	if err != nil {
		panic(err.Error())
	}
	
	receiver := R{}		
	r := nbr.NewSimulationReader(f, &receiver)
	
	
	namespace, err  := r.ParseNamespace()
	 if err != nil {panic(err)}
	fmt.Printf("namespace %s\n", namespace)
	err = r.Parse()
	 if err != nil {panic(err)}

}

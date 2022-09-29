package main

import (
	"fmt"
	"unsafe"
)

type testStruct struct {
	a struct{
		i int
	}
}


func main() {
	t := testStruct{}
	t2 := t
	fmt.Printf("s1 first element :%v\n", unsafe.Pointer(&(t.a.i)))
	fmt.Printf("s2 first element :%v\n", unsafe.Pointer(&(t2.a.i)))

	s1 := make([]int,4,4)
	s2 := s1
	s3 := &s1
	fmt.Printf("s1 first element :%v\n", unsafe.Pointer(&s1[0]))
	fmt.Printf("s2 first element :%v\n", unsafe.Pointer(&s2[0]))

	*s3 = append(*s3,1,2,3,4)
	s2 = append(s1, 1,2,3,4)
	fmt.Printf("s1 first element :%v\n", unsafe.Pointer(&s1[0]))
	fmt.Printf("s2 first element :%v\n", unsafe.Pointer(&s2[0]))
	fmt.Printf("s3 first element :%v\n", unsafe.Pointer(&(*s3)[0]))
}

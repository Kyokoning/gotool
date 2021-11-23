package size

import (
	"fmt"
	"strconv"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Structure1 struct {
	A int32 // 4
	B [16]int32 // 4*16 = 64
	C [3][100]Structure4
	D bool // 1
}

type Structure2 struct {
	A struct{}
	B int32
}

type Structure3 struct {
	A int32
	B struct{}
}

type Structure4 struct {
	A [7] Structure2
}

type Structure5 struct {
	Structure1
	Structure2
	Id int32
	Flag bool
	Structure3
}

type Structure6 struct {
	C struct {}
}

type T1 struct {
	F7 [5][5][5]uint8
}

type T3 struct {
	T1R1F1 [16]uint8
	T1R1F2 uint32
	T1R1F3 [8]uint32
	T1R1F4 [8]uint32
	T1R1F5 [8]uint32
	T1
	T1R1F6 bool
}


type T5 struct {
	T  T3
	C1 struct{}
}


func Test_getAlignSize4Structure(k *testing.T) {
	s1 := Structure1{B: [16]int32{}, C: [3][100]Structure4{}}
	s2 := Structure2{}
	s3 := Structure3{}
	s4 := Structure4{A: [7]Structure2{}}
	s5 := Structure5{s1, s2, 0, true, s3}
	s6 := Structure6{}
	t5 := T5{}

	// case1
	alignData, _ := getAlignSize(s1)
	except := fmt.Sprintf("%d", unsafe.Sizeof(s1))
	assert.Equal(k, except, strconv.Itoa(alignData.size))

	//case2
	alignData, _ = getAlignSize(s2)
	except = fmt.Sprintf("%d", unsafe.Sizeof(s2))
	assert.Equal(k, except, strconv.Itoa(alignData.size))

	//kk := [3][4]Manager{{*test2}}
	//fmt.Println(reflect.TypeOf(kk).Elem())

	//case3
	alignData, _ = getAlignSize(s3)
	except = fmt.Sprintf("%d", unsafe.Sizeof(s3))
	assert.Equal(k, except, strconv.Itoa(alignData.size))

	//case4
	alignData, _ = getAlignSize(s4)
	except = fmt.Sprintf("%d", unsafe.Sizeof(s4))
	assert.Equal(k, except, strconv.Itoa(alignData.size))

	//case5
	alignData, _ = getAlignSize(s5)
	except = fmt.Sprintf("%d", unsafe.Sizeof(s5))
	assert.Equal(k, except, strconv.Itoa(alignData.size))

	//case6
	alignData, _ = getAlignSize(s6)
	except = fmt.Sprintf("%d", unsafe.Sizeof(s6))
	assert.Equal(k, except, strconv.Itoa(alignData.size))

	//case T5
	alignData, _ = getAlignSize(t5)
	except = fmt.Sprintf("%d", unsafe.Sizeof(t5))
	assert.Equal(k, except, strconv.Itoa(alignData.size))
}

package size

import (
	"math"
	"reflect"
	"strconv"
)

var BasicLitSize = map[reflect.Kind]int{
	reflect.Bool:    1,
	reflect.Int8:    1,
	reflect.Uint8:   1,
	reflect.Int16:   1,
	reflect.Uint16:  2,
	reflect.Int32:   4,
	reflect.Uint32:  4,
	reflect.Float32: 4,
	reflect.Int:     strconv.IntSize / 8,
	reflect.Uint:    strconv.IntSize / 8,
	reflect.Int64:   8,
	reflect.Uint64:  8,
	reflect.Float64: 8,
}

var systemAlign = strconv.IntSize / 8

type alignData struct {
	size         int
	align        int
	maxFieldSize int // 如果不是struct，等于size。否则是最长字段大小
}

func getAlignSize(target interface{}) (*alignData, bool) {
	value := reflect.ValueOf(target)
	kind := value.Kind()

	size := 0
	align := 1
	switch kind {

	case reflect.Struct:
		maxFieldSize, preFieldSize := 0, 0
		numField := value.NumField()

		for i := 0; i < numField; i++ {
			element := value.Field(i).Interface()
			fieldAlignData, ok := getAlignSize(element)
			if !ok {
				return nil, false
			}

			// 如果最后一个字段是空结构体字段struct {}，会进行特殊填充，最后一个字段会被填充对齐到前一个字段的大小。
			if i == numField-1 && reflect.ValueOf(element).Kind() == reflect.Struct && fieldAlignData.size == 0 {
				size += preFieldSize
				break
			}

			// 每个成员相对于结构体首地址的offset
			// 是该成员大小与有效对齐值中较小那个的整数倍
			fieldAlignData.align = int(math.Min(float64(fieldAlignData.align), float64(fieldAlignData.size)))
			size = alignSize(size, fieldAlignData.align)
			size += fieldAlignData.size

			// 结构体x的对齐值alignOf(x)是结构体中最长字段对应的对齐值
			// (如果是structure字段，找该字段最长对齐值)
			if fieldAlignData.maxFieldSize > maxFieldSize {
				maxFieldSize = fieldAlignData.maxFieldSize
				align = fieldAlignData.align
			}
			preFieldSize = fieldAlignData.maxFieldSize
		}

		// 结构的长度必须是编译器默认的对齐长度和成员中最长类型中最小的数据大小的倍数对齐。
		size = alignSize(size, align)
		return &alignData{size, align, maxFieldSize}, true

	case reflect.Array:
		newElement := reflect.New(value.Type().Elem()).Elem().Interface()
		elementAlignData, ok := getAlignSize(newElement)
		if !ok {
			return nil, false
		}
		size += elementAlignData.size * value.Len()
		return &alignData{size, elementAlignData.align, elementAlignData.maxFieldSize}, true

	default:
		if isBasicLit(kind) {
			size := BasicLitSize[kind]
			align := int(math.Min(float64(size), float64(systemAlign)))
			return &alignData{size, align, size}, true
		}
		return nil, false
	}
}

func isBasicLit(kind reflect.Kind) bool {
	return kind == reflect.Bool || kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 ||
		kind == reflect.Int32 || kind == reflect.Int64 || kind == reflect.Uint || kind == reflect.Uint8 ||
		kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 || kind == reflect.Float32 ||
		kind == reflect.Float64
}

func alignSize(size, align int) int {
	if align != 0 && size%align != 0 {
		return (size/align + 1) * align
	}
	return size
}

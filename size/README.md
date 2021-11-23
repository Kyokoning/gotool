[TOC]

### 目标

根据输入的结构，直接推断go的结构大小（即实现一个`unsafe.Sizeof()`）

### 参考资料

[Go白皮书:Size and alignment guarantees](https://golang.google.cn/ref/spec#Size_and_alignment_guarantees)

[Go101内存布局](https://gfw.go101.org/article/memory-layout.html)

[Go 最细节篇 — 空结构体是什么?](https://juejin.cn/post/6908733156707287048)

### 实现方法

#### 输入

输入都是定长结构：

- 基本类型：Bool, Int, Int8, Int16 , Int32, Int64, Uint, Uint8 , Uint16, Uint32, Uint64, Float32 , Float64
- 复合类型：Array
- 复杂类型：Struct

#### 基本情况

`C语言`的对齐规则与`Go`语言一样，所以`C语言`的对齐规则对`Go`同样适用：

- 对于结构体的各个成员，第一个成员位于偏移为`0`的位置，结构体第一个成员的偏移量(offset)为`0`，以后**每个成员相对于结构体首地址的`offset`**都是**该成员大小与有效对齐值中较小那个的整数倍**，如有需要编译器会在成员之间加上填充字节。
- 除了结构成员需要对齐，结构本身也需要对齐，**结构的长度**必须是编译器默认的对齐长度和**成员中最长类型中最小的数据大小**的倍数对齐。

白皮书规定了以下几种基本类型的大小和对齐：

<img src="https://github.com/Kyokoning/gotool/blob/main/images/image-20211123150918245.png?raw=true"  alt="无法载入！">

其他的基本类型：`uint`和`int`的大小取决于编译器实现，32位架构是4，64位是8。

#### 特殊情况

**空结构体字段对齐**：编译器在遇到空结构体 `struct {}` 在**最后一个字段**的场景，会进行特殊填充，`struct {}` 作为最后一个字段，会被填充对齐到前一个字段的大小，地址偏移对齐规则不变。

### 实现

`getAlignSize(target interface{}) (int, int, bool)`

输入：interface()

输出：实际大小，对齐值，成功flag

使用`reflect.ValueOf`方法反射得到接口值的值信息。

如果结构值的类型是：

- **基本类型**：返回定义的size、align和true flag。

- **Array**：使用`reflect.TypeOf`方法获得接口值的值类型，使用`Type.Elem()`获得Array元素的值类型。使用`reflect.New()`方法生成该值类型的新对象的指针，然后将该指针指向的对象递归调用`getAlignSize`方法，得到Array元素的size、align和flag。使用`Value`的`Len()`方法获得Array长度。array长度乘元素大小即为array尺寸。单个元素的align即为array的align。

- **Struct**：将Struct的每一个字段递归调用`getAlignSize`方法。调用之后，计算offset（基本情况中的第一条规则）。在计算完整个结构之后，根据基本情况的第二条规则，对齐结构体长度。整个结构体的对齐值是结构体中最长字段对应的对齐值。


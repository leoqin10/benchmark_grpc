## 2. 序列化对性能的影响
序列化在RPC中是十分重要的一环，其将入参对象和返回值对象转换成可以存储或传输的格式在网络中进行传输
这里主要对比protobuf，messagePack，Gencode三种不同序列化对性能的影响

### 测试环境
`VM-3-227-ubuntu`
| CPU   | MEM  |
| ----- | ---- |
| 4core | 8G   |

### 测试数据
所有测试的数据结构都是一致的，它包含三个字段，一个是int类型的字段Id,一个是string类型的字段Name，一个是[]string类型的字段Colors
protobuf.proto文件内容如下：
```proto
message ProtoColorGroup {
    required int32 id = 1;
    required string name = 2;
    repeated string colors = 3;
}
```
```go
type ColorGroup struct {
    Id     int      `msg:"id"`
    Name   string   `msg:"name"`
    Colors []string `msg:"colors"`
}
```
gencode.schema文件内容如下：
```schema
struct GencodeColorGroup {
	Id     int32
	Name   string
	Colors []string
}
```
### 测试方法

测试命令：

```shell
go test -bench=. -benchmem
```

使用`go`自带的基准测试，在bench_test.go文件中定义如下类型的函数
```go
//序列化
func BenchmarkMarshalByProtoBuf(b *testing.B) {
    bb := make([]byte, 0, 1024)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        bb, _ = proto.Marshal(&protobufGroup)
    }
    b.ReportMetric(float64(len(bb)), "marshaledBytes")
}

//反序列化
func BenchmarkUnmarshalByProtoBuf(b *testing.B) {
    bytes, _ := proto.Marshal(&protobufGroup)
    result := model.ProtoColorGroup{}
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        proto.Unmarshal(bytes, &result)
    }
}
```
基准函数会运行 `b.N` 次。在基准执行期间，程序会自动调整 b.N 直到基准测试函数持续足够长的时间

### 测试结果

序列化：
| benchmark name    | 基准测试的迭代总次数 b.N | 平均每次迭代所消耗的纳秒数 | 平均每次迭代内存所分配的字节数 | 平均每次迭代的内存分配次数 | CPU（3%~5%） |
| ----------------- | ------------------------ | -------------------------- | ------------------------------ | -------------------------- | ------------ |
| MarshalByProtoBuf | 3164437                  | 380.2 ns/op                | 64 B/op                        | 2 allocs/op                | 34.2%        |
| MarshalByMsgp     | 13111096                 | 92.40 ns/op                | 80 B/op                        | 1 allocs/op                | 33.4%        |
| MarshalByGencode  | 28184382                 | **`42.94 ns/op`**          | 0 B/op                         | 0 allocs/op                | 31.0%        |


反序列化：
| benchmark name      | 基准测试的迭代总次数 b.N | 平均每次迭代所消耗的纳秒数 | 平均每次迭代内存所分配的字节数 | 平均每次迭代的内存分配次数 | CPU（3%~5%） |
| ------------------- | ------------------------ | -------------------------- | ------------------------------ | -------------------------- | ------------ |
| UnmarshalByProtoBuf | 1651257                  | 719.1 ns/op                | 176 B/op                       | 11 allocs/op               | 35.1%        |
| UnmarshalByMsgp     | 6710724                  | 163.9 ns/op                | 32 B/op                        | 5 allocs/op                | 30.6%        |
| UnmarshalByGencode  | 9098287                  | **`121.5 ns/op`**          | 32 B/op                        | 5 allocs/op                | 30.8%        |


### 结果分析
#### Protobuf
从结果中我们可以看到，Google官方版的Protobuf无论是在迭代时间，还是资源消耗上，都比较差，官方的protobuf是通过反射(reflect)实现序列化反序列化的，所以效率不是很高。
在调研中发现一种基于Protobuf的第三方实现的库[gogo/protobuf](https://github.com/gogo/protobuf), 该库在维持了原版的易用性上还提升了性能和增加了很多功能。从下表中可以看到其性能较原版有明显的提升。

| benchmark name          | 基准测试的迭代总次数 b.N | 平均每次迭代所消耗的纳秒数 | 平均每次迭代内存所分配的字节数 | 平均每次迭代的内存分配次数 | CPU（3%~5%） |
| ----------------------- | ------------------------ | -------------------------- | ------------------------------ | -------------------------- | ------------ |
| MarshalByGogoProtoBuf   | 11601102                 | 105.1 ns/op                | 48 B/op                        | 1 allocs/op                | 32.5%        |
| UnmarshalByGogoProtoBuf | 2608179                  | 469.8 ns/op                | 160 B/op                       | 10 allocs/op               | 31.7%        |

#### MessagePack

MessagePack在测试中也表现优异，它的序列化和反序列化时间比protobuf要少很多，同时它对于资源的消耗也比protobuf小，同时它也是跨语言的
#### Gencode
Gencode在这次测试中表现最为出色，同时由于 gencode 在序列化时不会写入字段的名字，所以生成的字节体积很小
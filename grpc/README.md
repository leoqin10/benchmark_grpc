## 多链接对性能的影响 grpc
### 测试环境
服务器和客户端均在`VM-3-227-ubuntu`上进行

| CPU   | MEM  |
| ----- | ---- |
| 4core | 8G   |

### 测试场景
模拟简单业务，处理主要是将反序列化的对象设置几个字段后返回，在保证相同流量下，不断增加连接数量
### 测试数据
使用protobuf进行序列化和反序列化
数据格式： 序列化之后大小为581字节
```proto
message BenchmarkMessage {
  required string field1 = 1;
  optional string field9 = 9;
  optional string field18 = 18;
  optional bool field80 = 80 [default=false];
  optional bool field81 = 81 [default=true];
  required int32 field2 = 2;
  required int32 field3 = 3;
  optional int32 field280 = 280;
  optional int32 field6 = 6 [default=0];
  optional int64 field22 = 22;
  optional string field4 = 4;
  repeated fixed64 field5 = 5;
  optional bool field59 = 59 [default=false];
  optional string field7 = 7;
  optional int32 field16 = 16;
  optional int32 field130 = 130 [default=0];
  optional bool field12 = 12 [default=true];
  optional bool field17 = 17 [default=true];
  optional bool field13 = 13 [default=true];
  optional bool field14 = 14 [default=true];
  optional int32 field104 = 104 [default=0];
  optional int32 field100 = 100 [default=0];
  optional int32 field101 = 101 [default=0];
  optional string field102 = 102;
  optional string field103 = 103;
  optional int32 field29 = 29 [default=0];
  optional bool field30 = 30 [default=false];
  optional int32 field60 = 60 [default=-1];
  optional int32 field271 = 271 [default=-1];
  optional int32 field272 = 272 [default=-1];
  optional int32 field150 = 150;
  optional int32 field23 = 23 [default=0];
  optional bool field24 = 24 [default=false];
  optional int32 field25 = 25 [default=0];
  optional bool field78 = 78;
  optional int32 field67 = 67 [default=0];
  optional int32 field68 = 68;
  optional int32 field128 = 128 [default=0];
  optional string field129 = 129 [default="xxxxxxxxxxxxxxxxxxxxx"];
  optional int32 field131 = 131 [default=0];
}
```
### 测试client
通过`-c`来控制并发数量和`-n`来控制请求数量
- `-c`:  指定并发的client，每个client都是独立的，我通过循环调用`grpc.Dial()`来创建多个连接，再将这些连接放入到容器中来形成连接池
- `-n`:  要发送的总请求数，由client平分，为了更清晰的看到性能变化情况，我们使用较多的请求数量，以此来看到较大的性能变化，经过反复测试，选定请求总数为20000个，故总流量大约在`12M`

### 测试指标
- cpu占用率：通过netdata来观察cpu曲线图变化情况，在不开启连接的情况下，cpu占用率在3%~5%
- 可用内存：通过`free -s 1 -k`命令来采集，在不开启连接情况下，内存占用率在35.4%左右
- TPS：每秒完成的请求数量，请求总数/总共所消耗的时间
- 平均延迟：服务发出到收到response所需的时间的平均值

### 测试结果
每个结果均为三次测试的平均值，以此来保证准确性


![](https://github.com/leoqin10/benchmark_grpc/blob/main/grpc/image/%E5%A4%9A%E8%BF%9E%E6%8E%A5cpu%E3%80%81%E5%86%85%E5%AD%98%E6%8A%98%E7%BA%BF%E5%9B%BE.png)

![](https://github.com/leoqin10/benchmark_grpc/blob/main/grpc/image/%E5%B9%B3%E5%9D%87%E5%BB%B6%E8%BF%9F.png)


| 项目/连接数量 | cpu（3%-5%） | 内存（35.4%） | tps   | 平均延迟ns |
| ------------- | ------------ | ------------- | ----- | ---------- |
| 20            | 63.3%        | 35.802%       | 23254 | 844926     |
| 100           | 67.83%       | 36.080%       | 22557 | 2983850    |
| 500           | 70.62%       | 35.805%       | 24332 | 13721681   |
| 1000          | 73.48%       | 35.994%       | 23714 | 27231117   |
| 10000         | 80.48%       | 36.596%       | 29183 | 215628634  |


### 结果分析
从图中可以看出，多连接使得cpu占用率明显上升，对于内存的影响不大，cpu上升的主要原因在于不断在创建和断开连接的过程中对cpu的消耗，吞吐率有所下降但并不明显，响应的平均延迟有明显的上升。

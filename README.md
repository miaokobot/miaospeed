# MiaoSpeed 4.0

---

> miaospeed 于 4.0.0 与 miaoko 分离，正式成为独立的开源项目。一般来说，miaospeed 依然被认为是 miaoko 的专用后端，但也能成为一个通用型后端。

## 版权与协议

miaospeed 采用 AGPLv3 协议开源，您可以按照 AGPLv3 协议对 miaospeed 进行修改、贡献、分发、乃至商用。但请切记，您必须遵守 AGPLv3 协议下的一切义务，以免发生不必要的法律纠纷。

### 主要开源依赖公示

miaospeed 采用了如下的开源项目:

- Dreamacro/clash [GPLv3]
- juju/ratelimit [LGPLv3]
- dop251/goja [MIT]
- json-iterator/go [MIT]
- pion/stun [MIT]
- go-yaml/yaml [MIT]
- gorilla/websocket [BSD]

## 抽象与模块

如果您想贡献 miaospeed，您可以参考以下 miaospeed 的抽象设计:

- **Matrix**: 数据矩阵 [interfaces/matrix.go]。即用户想要获取的某个数据的最小颗粒度。例如，用户希望了解某个节点的 RTT 延迟，则 TA 可以要求 miaospeed 对 `TEST_PING_RTT` [例如: service/matrices/httpping/matrix.go] 进行测试。
- **Macro**: 运行时宏任务 [interfaces/macro.go]。如果用户希望批量运行数据矩阵，他们往往会做重复的事情。例如 `TEST_PING_RTT` 与 `TEST_PING_HTTP` 大多数时间都在做相同的事情。如果将两个 _Matrix_ 独立运行，则会浪费大量资源。因此，我们定义了 _Macro_ 最为一个最小颗粒度的执行体。由 _Macro_ 并行完成一系列耗时的操作，随后，_Matrix_ 将解析 _Macro_ 运行得到的数据，以填充自己的内容。
- **Vendor**: 服务提供商 [interfaces/vendor.go]。miaospeed 本身只是一个测试工具，**它不具备任何代理能力**。因此，_Vendor_ 作为一个接口，为 miaospeed 提供了链接能力。

## 基本使用方式

### 二进制

关于二进制的使用，本 _README_ 不做赘述，请手动编译后执行 `./miaospeed -help` 查看。

### 对接方法

由于 _miaoko_ 是闭源软件/服务，如果您想在其他服务内对接 miaospeed，可能没有现成的案例。但是，您依然可以参考如下文件:

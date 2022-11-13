# MiaoSpeed 4.0

---

> miaospeed 于 4.0.0 与 miaoko 分离，正式成为独立的开源项目。一般来说，miaospeed 依然被认为是 miaoko 的专用后端，但也能成为一个通用型后端。

## 基本使用方式

### 二进制

关于二进制的使用，本 _README_ 不做赘述，请手动编译或下载预编译文件后执行 `./miaospeed` 查看。

### 编译

由于 miaospeed 中含有部分证书与脚本并未开源，您需要补齐以下文件以成功编译:

1. `./utils/embeded/BUILDTOKEN.key`: 这是 `编译TOKEN`，它用于签名 miaospeed request 的结构体，以防止您的客户端使用不合规的 miaospeed 造成数据不真实的纠纷。您可以随便定义它，例如: `1111|2222|33333333`，不同段用 `|` 切开。
2. `./preconfigs/embeded/miaokoCA/miaoko.crt`: 当 `-mtls` 启用时，miaospeed 会读取这里的证书让客户端做 TLS 验证。
3. `./preconfigs/embeded/miaokoCA/miaoko.key`: 同上，这是私钥。(对于这两个您可以自己用 openssl 签一个证书，但它不能用于 miaoko。)
4. `./preconfigs/embeded/ca-certificates.crt`: miaospeed 自带的根证书集，防止有恶意用户修改系统更证书以作假 TLS RTT。（对于 debian 用户，您可以在安装 `ca-certificates` 包后，在 `/etc/ssl/certs/ca-certificates.crt` 获取这个文件）
5. `./engine/embeded/predefined.js`: 这个文件定义了 `JavaScript` (流媒体)脚本中一些通用方法，例如 `get()`, `safeStringify()`, `safeParse()`, `println()`，您可以自己实现它们，或者只是新建一个空文件。
6. `./engine/embeded/default_geoip.js`: 默认的 `geoip` 脚本，需要提供一个 `handler()` 入口函数。如果您不想提供默认的 `geoip` 脚本，则可以直接新建空文件（或者您可以直接拷贝 miaoko 提供的 `geoip` 脚本）。
7. `./engine/embeded/default_ip.js`: 默认的 `ip_resolve` 脚本，需要提供一个 `ip_resolve_default()` 入口函数，用于获取入口、出口的 IP。如果您不想提供默认的 `ip_resolve` 脚本，则可以直接新建空文件（或者您可以直接拷贝 miaoko 提供的 `geoip` 脚本里的 ip_resolve 函数）。

当您新建好以上文件后，就可以运行 `go build .` 构建 `miaospeed` 了。

### 对接方法

由于 _miaoko_ 是闭源软件/服务，如果您想在其他服务内对接 miaospeed，可能没有现成的案例。但是，您依然可以参考如下思路:

0. miaospeed 对接本质是通过 ws 通道发送指令、传递信息。一般来说，您只需要连接 ws，构建请求结构体，签名请求，接收结果即可。
1. 连接 ws，这一步很简单，也就不用赘述了。（如果您在客户端强制断开链接，则任务会被自动中止）
2. 构建请求结构体，参考: https://github.com/miaokobot/miaospeed/blob/fd7abecc2d36a0f18b08f048f9a53b7c0a26bd9e/interfaces/api_request.go#L50
3. 签名，参考: https://github.com/miaokobot/miaospeed/blob/df6202409e87c5d944ab756608fd31d35390b5c0/utils/challenge.go#L39 其中需要传入两个参数。第一个参数是 `启动TOKEN` （即您启动 miaospeed 时传入的 -token 后的内容），第二个就是在第二步中您构建的结构体 `req`。签名的方法，通俗一些说明就是将结构体转换为 JSON String 然后与 `启动TOKEN` 和 `编译TOKEN` 切片分别累积做 SHA512 HASH。最后，将签名的字符串写入 `req.Challenge` 即可。
4. 发送完成签名后的请求，您就可以接收返回值了。服务器返回的结构体统一为 https://github.com/miaokobot/miaospeed/blob/fd7abecc2d36a0f18b08f048f9a53b7c0a26bd9e/interfaces/api_response.go#L28

## 版权与协议

miaospeed 采用 AGPLv3 协议开源，您可以按照 AGPLv3 协议对 miaospeed 进行修改、贡献、分发、乃至商用。但请切记，您必须遵守 AGPLv3 协议下的一切义务，以免发生不必要的法律纠纷。

### 主要开源依赖公示

miaospeed 采用了如下的开源项目:

- Dreamacro/clash [GPLv3]
- MetaCubeX/Clash.Meta [GPLv3]
- juju/ratelimit [LGPLv3]
- dop251/goja [MIT]
- json-iterator/go [MIT]
- pion/stun [MIT]
- go-yaml/yaml [MIT]
- gorilla/websocket [BSD]

## 抽象设计

如果您想贡献 miaospeed，您可以参考以下 miaospeed 的抽象设计:

- **Matrix**: 数据矩阵 [interfaces/matrix.go]。即用户想要获取的某个数据的最小颗粒度。例如，用户希望了解某个节点的 RTT 延迟，则 TA 可以要求 miaospeed 对 `TEST_PING_RTT` [例如: service/matrices/httpping/matrix.go] 进行测试。
- **Macro**: 运行时宏任务 [interfaces/macro.go]。如果用户希望批量运行数据矩阵，他们往往会做重复的事情。例如 `TEST_PING_RTT` 与 `TEST_PING_HTTP` 大多数时间都在做相同的事情。如果将两个 _Matrix_ 独立运行，则会浪费大量资源。因此，我们定义了 _Macro_ 最为一个最小颗粒度的执行体。由 _Macro_ 并行完成一系列耗时的操作，随后，_Matrix_ 将解析 _Macro_ 运行得到的数据，以填充自己的内容。
- **Vendor**: 服务提供商 [interfaces/vendor.go]。miaospeed 本身只是一个测试工具，**它不具备任何代理能力**。因此，_Vendor_ 作为一个接口，为 miaospeed 提供了链接能力。

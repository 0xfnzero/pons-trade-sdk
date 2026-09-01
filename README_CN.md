<div align="center">
    <h1>Pons Trade SDK for Go</h1>
    <h3><em>面向生产环境的 pons v2 发射、Bonding Curve 交易、手续费与毕业流程 Go SDK</em></h3>
</div>

<p align="center">
    <strong>基于官方 ponsfamily v2 文档与 Solidity 源码，为 Robinhood Chain 上的 pons v2 提供类型化读写与 calldata 构建能力。</strong>
</p>

<p align="center">
    <a href="https://pkg.go.dev/github.com/0xfnzero/pons-trade-sdk">
        <img src="https://pkg.go.dev/badge/github.com/0xfnzero/pons-trade-sdk.svg" alt="Go Reference">
    </a>
    <a href="https://github.com/0xfnzero/pons-trade-sdk/actions/workflows/ci.yml">
        <img src="https://github.com/0xfnzero/pons-trade-sdk/actions/workflows/ci.yml/badge.svg" alt="CI">
    </a>
    <a href="https://github.com/0xfnzero/pons-trade-sdk/releases/latest">
        <img src="https://img.shields.io/github/v/release/0xfnzero/pons-trade-sdk" alt="Release">
    </a>
    <a href="LICENSE">
        <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License">
    </a>
</p>

<p align="center">
    <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go">
    <img src="https://img.shields.io/badge/EVM-3C3C3D?style=for-the-badge&logo=ethereum&logoColor=white" alt="EVM">
    <img src="https://img.shields.io/badge/Pons-v2-4B8BBE?style=for-the-badge" alt="Pons v2">
</p>

<p align="center">
    <a href="./README_CN.md">中文</a> |
    <a href="./README.md">English</a> |
    <a href="https://fnzero.dev/">Website</a> |
    <a href="https://t.me/fnzero_group">Telegram</a> |
    <a href="https://discord.gg/vuazbGkqQE">Discord</a>
</p>

---

## SDK

| SDK | Module |
|-----|--------|
| Trade | [`github.com/0xfnzero/pons-trade-sdk`](https://github.com/0xfnzero/pons-trade-sdk) |
| Parser | [`github.com/0xfnzero/pons-parser-sdk`](https://github.com/0xfnzero/pons-parser-sdk) |

## 这个 SDK 适合什么场景

`pons-trade-sdk` 适合需要接入 Robinhood Chain 上 pons v2 的 Go 服务、交易机器人、发射工具、keeper 和索引后端。集成行为以官网文档为第一优先级；官网没有写完整的 ABI、事件、错误或状态字段，再由公开 V2 Solidity 源码补齐。

| 范围 | 覆盖内容 |
|------|----------|
| 链 | Robinhood Chain，chain id `4663` |
| 官方参考 | [docs.ponsfamily.com/v2](https://docs.ponsfamily.com/v2)、[ponsdotdev/ponsfamily](https://github.com/ponsdotdev/ponsfamily) |
| 交易流程 | Launch、launch-and-buy、buy、sell、approve、手续费领取、手续费 sweep、毕业 |
| 读取能力 | Factory 配置、launch 记录、curve 储备、hook 费用、vault vesting、token 元数据 |
| 运行环境 | Go 1.25+；生产构建请使用仍受支持且已安装安全补丁的 Go 版本 |

## 功能

1. 内置 Robinhood Chain 上官方 pons v2 合约地址。
2. 基于 ABI 的 calldata builder，覆盖 launch、buy、sell、claim、sweep、graduation。
3. 基于 `bind.ContractBackend` 的 `Client`，提供类型化读取与交易方法。
4. 确定性 bonding curve 报价工具，可用于 buy/sell 预估和滑点保护。
5. 支持重建毕业后的 Uniswap v4 pool key 和 pool id。
6. 支持手续费 escrow 与 buyback vesting 读取。
7. 支持解析 pons v2、OpenZeppelin ownership、SafeERC20、ReentrancyGuard 的 revert data。
8. 同仓提供 `pons-parser-sdk`，用于 pons v2 日志和 receipt 解析。

## 安装

### 源码克隆

在项目根目录中克隆 SDK，并固定到已发布的 `v0.1.0`：

```bash
cd your_project_root_directory
git clone --branch v0.1.0 --depth 1 https://github.com/0xfnzero/pons-trade-sdk
```

在您项目的 `go.mod` 中添加本地模块依赖：

```go
require github.com/0xfnzero/pons-trade-sdk v0.1.0

replace github.com/0xfnzero/pons-trade-sdk => ./pons-trade-sdk
```

然后解析依赖，并按正常方式导入 SDK package：

```bash
go mod tidy
```

```go
import "github.com/0xfnzero/pons-trade-sdk/ponstrade"
```

### Go Modules

```bash
go get github.com/0xfnzero/pons-trade-sdk@v0.1.0
```

## 使用示例

### 读取 Factory 状态

```go
client, eth, err := ponstrade.Dial(ctx, "https://your-robinhood-chain-rpc")
if err != nil {
    panic(err)
}
defer eth.Close()

fee, err := client.LaunchFee(ctx, nil)
enabled, err := client.LaunchEnabled(ctx, nil)
configs, err := client.OpenLaunchConfigs(ctx, nil)
```

`Dial` 返回前会校验 Robinhood Chain id `4663` 和 SDK 使用的全部合约目标地址。注入自定义或本地链 `bind.ContractBackend` 时请使用 `NewClientChecked`；`NewClient` 继续支持只配置部分能力的离线 client，各方法会在调用时校验自身需要的目标地址。

`OpenLaunchConfigs` 默认最多读取 4096 个配置，并发 RPC 上限为 16。需要其他边界时，请通过 `OpenLaunchConfigsWithOptions` 显式设置：

```go
configs, err := client.OpenLaunchConfigsWithOptions(ctx, ponstrade.LaunchConfigQueryOptions{
    MaxConfigs: 8192,
    Concurrency: 32,
}, nil)
```

### 预估买入

```go
quote, err := client.QuoteBuy(ctx, curve, quoteIn, recipient, nil)
if err != nil {
    panic(err)
}

minOut, err := ponstrade.MinTokensOutForBuy(quoteIn, quote, 500) // 5%
if err != nil {
    panic(err)
}
```

Curve 买入请使用 `MinTokensOutForBuy`。它会在临近毕业发生部分成交时继续按报价成交率保护价格。`MinOutputWithSlippage` 仍适用于卖出报价等普通输出金额；金额无效或滑点超过 100% 时会返回错误。

曲线已毕业或即将毕业时，`QuoteBuy` 和 `QuoteSell` 返回 `ErrCurveClosed`，不会继续给出不可执行报价。状态无效或输出舍入为零时返回 `ErrInvalidQuote` 或 `ErrUnquotableTrade`，不要据此构建交易。

### 构建原生 quote 买入交易

```go
auth.Value = quoteIn // 原生 quote 买入要求 msg.value == quoteIn

tx, err := client.Buy(auth, curve, quoteIn, minOut, recipient)
if err != nil {
    revertData := []byte{} // 从 RPC / eth_call / estimate error 中拿到 revert data 后填入
    if decoded, ok, parseErr := ponstrade.ParseContractError(revertData); parseErr == nil && ok {
        panic(decoded)
    }
    panic(err)
}
```

ERC-20 quote 买入需要先授权 curve 花费 quote token，并保持 `auth.Value == 0`。

### 发射 Token

```go
salt, err := ponstrade.RandomSalt()
if err != nil {
    panic(err)
}

expected, err := client.PreviewLaunchEconomics(ctx, big.NewInt(0), ponstrade.NativeQuote, nil)
if err != nil {
    panic(err)
}

params := ponstrade.TokenParams{
    Name:              "Example",
    Symbol:            "EXMPL",
    CreatorTaxBps:     100,
    BuybackEnabled:    true,
    ExpectedEconomics: ponstrade.HashToBytes32(expected),
    Salt:              salt,
}

auth.Value = launchFee
tx, err := client.LaunchToken(auth, params, big.NewInt(0), ponstrade.NativeQuote, nil)
```

### 原子发射并买入

```go
auth.Value = ponstrade.LaunchAndBuyValue(ponstrade.NativeQuote, launchFee, quoteIn)

tx, err := client.LaunchAndBuy(
    auth,
    params,       // router 路径里 creatorFeeRecipient 建议显式设置
    big.NewInt(0),
    ponstrade.NativeQuote,
    quoteIn,
    minTokensOut,
    recipient,
    nil,
)
```

ERC-20 pair launch 需要先授权 launch-and-buy router 花费 `quoteIn`，然后使用 `LaunchAndBuyValue(pairToken, launchFee, quoteIn)`；该 helper 会只返回原生 launch fee。

### 毕业与手续费

```go
launched, err := client.GetLaunchedToken(ctx, token, nil)
poolID, err := ponstrade.PoolID(ponstrade.BuildPoolKey(launched, client.Addresses().MemeHook))

tx, err := client.Graduate(auth, token)
tx, err = client.CreateGraduatedPool(auth, token)
tx, err = client.SweepPoolFees(auth, poolID, minConversionQuoteOut, minBuybackTokensOut)
```

## 项目结构

```text
.
├── ponstrade/                 # Trade SDK package
├── examples/basic/            # 基础报价示例
├── .github/workflows/ci.yml   # GitHub Actions 验证
├── go.mod
└── go.sum
```

## 开发

```bash
go test ./...
```

## 重要说明

1. 集成流程和部署地址以官网文档为第一优先级。官网缺失的完整 ABI、事件、错误和状态字段，由 `ponsdotdev/ponsfamily` 的 `contractsV2/src/v2` 补齐。
2. 原生 quote 买入要求 `msg.value == quoteIn`；ERC-20 quote 买入要求 `msg.value == 0`。
3. `launchAndBuy` 按官网公开的地址和 ABI 实现。公开 V2 源码通过 `launchForwarder` / `launchTokenFor` 引用了这个 trusted router，但没有包含 router 合约源码。
4. SDK 通过 go-ethereum 抽象构建和发送交易；nonce、gas、费用策略、私钥托管、RPC 重试策略由调用方控制。
5. 生产资金交易前请先在 fork 或测试环境做模拟。
6. 复合读取会并发发起 RPC，并在 backend 支持 `BlockNumber` 时固定到同一区块。自定义 backend 不支持该方法时，请显式传入 `bind.CallOpts.BlockNumber`。
7. `BPS` 和 `OnePctBPS` 是不可变数值常量。需要 `*big.Int` 时使用 `new(big.Int).SetUint64(ponstrade.BPS)` 转换。
8. `WithAddresses` 会替换完整 SDK 地址集合。`Dial` 和 `NewClientChecked` 会拒绝固定合约目标为零地址；零地址只在表示 `NativeQuote` 的参数位置合法。
9. 构建 calldata 前会按 UTF-8 字节数校验 launch metadata：name 64、symbol 16、logo 512、description 2,048、每个 social 字段 256 字节，与 V2 launch deployer 一致。
10. 导出的 ABI 变量是兼容性快照。SDK 内部使用独立的不可变解析结果，修改导出 ABI 不会改变交易或读取行为。
11. non-payable 交易 helper 遇到非零 `auth.Value` 时会返回 `ErrUnexpectedTransactionValue`，避免误带原生资产并在链上回滚。launch 和 buy 属于 payable 路径，准确 value 取决于实时 launch fee 与 quote asset，仍由调用方设置。
12. `ParseContractError` 对已识别 revert data 设置 1 MiB 上限，并拒绝尾随数据等非 canonical ABI 编码。若两个 custom error 发生 4-byte selector 冲突，会返回歧义错误，不会任意套用错误名。

可使用以下命令对 Robinhood Chain 已部署合约执行只读检查：

```bash
ROBINHOOD_RPC_URL=https://your-robinhood-chain-rpc go test ./ponstrade -run TestRobinhoodDeployment
```

## 许可证

MIT

## 联系

- Website: [fnzero.dev](https://fnzero.dev/)
- Telegram: [fnzero_group](https://t.me/fnzero_group)
- Discord: [FnZero Discord](https://discord.gg/vuazbGkqQE)

<div align="center">
    <h1>Pons Trade SDK for Go</h1>
    <h3><em>Production-oriented Go SDK for pons v2 launches, bonding-curve trades, fees, and graduation</em></h3>
</div>

<p align="center">
    <strong>Typed helpers for building and reading pons v2 transactions on Robinhood Chain, aligned with the official ponsfamily v2 docs and Solidity source.</strong>
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

## SDKs

| SDK | Module |
|-----|--------|
| Trade | [`github.com/0xfnzero/pons-trade-sdk`](https://github.com/0xfnzero/pons-trade-sdk) |
| Parser | [`github.com/0xfnzero/pons-parser-sdk`](https://github.com/0xfnzero/pons-parser-sdk) |

## What This SDK Is For

`pons-trade-sdk` is a Go SDK for services, bots, launch tooling, keepers, and indexer backends that interact with pons v2 on Robinhood Chain. Integration behavior follows the official docs first; when the docs omit a complete ABI/event/error/state detail, the SDK fills that gap from the public V2 Solidity source.

| Area | Coverage |
|------|----------|
| Chain | Robinhood Chain, chain id `4663` |
| Official references | [docs.ponsfamily.com/v2](https://docs.ponsfamily.com/v2), [ponsdotdev/ponsfamily](https://github.com/ponsdotdev/ponsfamily) |
| Trading workflows | Launch, launch-and-buy, buy, sell, approve, fee claim, fee sweep, graduation |
| Reads | Factory config, launch records, curve reserves, hook fees, vault vesting, token metadata |
| Runtime | Go 1.25+; use a supported, fully patched Go release for production builds |

## Features

1. Official pons v2 contract addresses for Robinhood Chain.
2. ABI-backed calldata builders for launch, buy, sell, claim, sweep, and graduation calls.
3. `Client` helpers over `bind.ContractBackend` for typed reads and transactions.
4. Deterministic bonding-curve quote helpers for buy/sell previews and slippage bounds.
5. Uniswap v4 pool key and pool id reconstruction for graduated launches.
6. Fee escrow and buyback vesting helpers.
7. Custom error decoding for pons v2, OpenZeppelin ownership, SafeERC20, and reentrancy errors.
8. Sibling parser SDK for pons v2 logs and receipt parsing.

## Installation

### Direct Clone

Clone this SDK into your project directory at the released `v0.1.0` tag:

```bash
cd your_project_root_directory
git clone --branch v0.1.0 --depth 1 https://github.com/0xfnzero/pons-trade-sdk
```

Add the local module to your project's `go.mod`:

```go
require github.com/0xfnzero/pons-trade-sdk v0.1.0

replace github.com/0xfnzero/pons-trade-sdk => ./pons-trade-sdk
```

Then resolve dependencies and import the SDK package normally:

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

## Usage Examples

### Read Factory State

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

`Dial` verifies Robinhood Chain id `4663` and all SDK contract destinations before returning. When injecting a custom or local-chain `bind.ContractBackend`, use `NewClientChecked`; `NewClient` remains available for partial/offline clients whose methods validate their own required destination at call time.

`OpenLaunchConfigs` defaults to at most 4,096 configurations and 16 concurrent RPC calls. Use `OpenLaunchConfigsWithOptions` to choose a different explicit bound:

```go
configs, err := client.OpenLaunchConfigsWithOptions(ctx, ponstrade.LaunchConfigQueryOptions{
    MaxConfigs: 8192,
    Concurrency: 32,
}, nil)
```

### Quote A Curve Buy

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

Use `MinTokensOutForBuy` for curve buys. It preserves the quoted price bound when the curve partially fills a buy near graduation. `MinOutputWithSlippage` remains suitable for ordinary output amounts such as sell quotes and returns an error for invalid amounts or slippage above 100%.

`QuoteBuy` and `QuoteSell` return `ErrCurveClosed` instead of quoting a curve that is graduated or ready to graduate. Invalid, zero-output states return `ErrInvalidQuote` or `ErrUnquotableTrade`; do not build a transaction from those results.

### Build A Native Buy Transaction

```go
auth.Value = quoteIn // native quote buys require msg.value == quoteIn

tx, err := client.Buy(auth, curve, quoteIn, minOut, recipient)
if err != nil {
    revertData := []byte{} // fill from your RPC / eth_call / estimate error when available
    if decoded, ok, parseErr := ponstrade.ParseContractError(revertData); parseErr == nil && ok {
        panic(decoded)
    }
    panic(err)
}
```

For ERC-20 quote buys, approve the curve to spend the quote token and keep `auth.Value` at zero.

### Launch A Token

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

### Launch And Buy Atomically

```go
auth.Value = ponstrade.LaunchAndBuyValue(ponstrade.NativeQuote, launchFee, quoteIn)

tx, err := client.LaunchAndBuy(
    auth,
    params,       // creatorFeeRecipient should be explicit for the router path
    big.NewInt(0),
    ponstrade.NativeQuote,
    quoteIn,
    minTokensOut,
    recipient,
    nil,
)
```

For ERC-20 pair launches, approve the launch-and-buy router for `quoteIn` first and use `LaunchAndBuyValue(pairToken, launchFee, quoteIn)`, which returns only the native launch fee.

### Graduation And Fees

```go
launched, err := client.GetLaunchedToken(ctx, token, nil)
poolID, err := ponstrade.PoolID(ponstrade.BuildPoolKey(launched, client.Addresses().MemeHook))

tx, err := client.Graduate(auth, token)
tx, err = client.CreateGraduatedPool(auth, token)
tx, err = client.SweepPoolFees(auth, poolID, minConversionQuoteOut, minBuybackTokensOut)
```

## Project Structure

```text
.
├── ponstrade/                 # Trade SDK package
├── examples/basic/            # Basic quote example
├── .github/workflows/ci.yml   # GitHub Actions verification
├── go.mod
└── go.sum
```

## Development

```bash
go test ./...
```

## Important Notes

1. Integration flow and deployed addresses follow the official docs first. Missing full ABI, event, error, and state surfaces are filled from `contractsV2/src/v2` in `ponsdotdev/ponsfamily`.
2. Native quote buys require `msg.value == quoteIn`; ERC-20 quote buys require `msg.value == 0`.
3. `launchAndBuy` is implemented from the official docs address and ABI. The public V2 source tree references this trusted router through `launchForwarder` / `launchTokenFor`, but does not include the router contract source.
4. This SDK builds and submits through go-ethereum abstractions; nonce, gas, fee strategy, private key custody, and RPC retry policy stay with the caller.
5. Run fork or testnet simulations before sending production value.
6. Composite reads use concurrent RPC calls and pin them to one block when the backend exposes `BlockNumber`. Pass an explicit `bind.CallOpts.BlockNumber` when using a custom backend that does not expose it.
7. `BPS` and `OnePctBPS` are immutable numeric constants. Convert them with `new(big.Int).SetUint64(ponstrade.BPS)` when a `*big.Int` is required.
8. `WithAddresses` replaces the complete SDK address set. `Dial` and `NewClientChecked` reject a zero fixed-contract destination; the zero address remains valid only where it denotes `NativeQuote`.
9. Launch metadata is validated by UTF-8 byte length before calldata is built: name 64, symbol 16, logo 512, description 2,048, and each social field 256 bytes, matching the V2 launch deployer.
10. Exported ABI variables are compatibility snapshots. SDK internals use independent immutable parses, so changing an exported ABI does not alter transaction or read behavior.
11. Transaction helpers for non-payable methods reject a non-zero `auth.Value` with `ErrUnexpectedTransactionValue`, preventing accidental native value from reaching a reverting call. Payable launch and buy helpers leave the exact value policy to the caller because it depends on the live launch fee and quote asset.
12. `ParseContractError` accepts at most 1 MiB of recognized revert data and rejects trailing or otherwise non-canonical ABI encodings. A colliding 4-byte custom-error selector is reported as ambiguous instead of being assigned an arbitrary name.

Run the read-only deployment checks against Robinhood Chain with:

```bash
ROBINHOOD_RPC_URL=https://your-robinhood-chain-rpc go test ./ponstrade -run TestRobinhoodDeployment
```

## License

MIT

## Contact

- Website: [fnzero.dev](https://fnzero.dev/)
- Telegram: [fnzero_group](https://t.me/fnzero_group)
- Discord: [FnZero Discord](https://discord.gg/vuazbGkqQE)

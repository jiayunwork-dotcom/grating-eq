# grating-eq: Go headless HTTP API and CLI for plane transmission grating orders

用户给出光栅槽距 d（或线密度）、波长 λ、入射角 θi 与级次 m，工具按光栅方程 mλ=d(sinθm−sinθi) 算出各级衍射角、是否存在、角色散 dθ/dλ、自由光谱范围 FSR≈λ/|m| 与分辨本领 R=mN。必须同时成立：|sinθm|≤1 的级次才存在；正入射时 ±m 角度对称；d 与线密度互为倒数。 d≤0、λ≤0、|sinθi|>1 或 N<1 必须报错。

## How to run

```
go run . serve -addr :8080
go run . orders example/600lpmm.json
```

## API

- `POST /api/orders` — grating JSON returns the order table.
- `POST /api/checks` — wavelength / density / zero-order / ±m / R-linear cross-checks.
- `GET /api/health` — liveness.

## Build

```
go build ./...
go test ./...
```

Module targets Go 1.21.

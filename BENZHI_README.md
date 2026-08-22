# grating-eq — 平面衍射光栅级次与色散核算

grating-eq 是平面衍射光栅级次与色散核算命令行工具：给定槽距（或线密度）、波长、入射角与可选有限缝数的 JSON 文件，按 mλ=d(sinθ_m−sinθ_i) 算出各级次衍射角、判定可见级次，并给出角色散、自由光谱范围、分辨本领与最小可分波长间隔。纯标准库，无网络依赖，无 cgo。

## 构建 / 运行 / 测试

```text
go build ./...
go run . orders example/600lpmm.json   # CLI：核算算例并打印各级次 θ、是否存在、色散
go run . checks example/600lpmm.json   # 交叉规则自检
go test ./...                          # 单元测试（grating / disperse / report）
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```

进容器后运行 `go build ./... && go test ./...`，再用 `go run . orders example/600lpmm.json` 验证 CLI 输出 ±1 级 θ≈19.27° 的核算结果。

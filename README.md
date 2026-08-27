# grating-eq — 平面衍射光栅核算

grating-eq 以 headless HTTP 服务加少量 CLI 交付。给定平面透射光栅的槽距 d（或线密度）、照明波长 λ、入射角 θi 与可选光照缝数 N（JSON），它按光栅方程 mλ = d(sinθm − sinθi) 逐级次算出衍射角 θm，判定每个级次是否落在可见区（|sinθm|≤1），并给出角色散 dθ/dλ = m/(d·cosθm)、自由光谱范围 FSR ≈ λ/|m|，以及给定 N 时的分辨本领 R = mN 与最小可分波长 Δλ = λ/R。有限缝阵列因子的主峰半宽经角色散换回的 Rayleigh 间隔必须与 λ/R 一致。能力边界为平面透射光栅的一维衍射：角度一律从法线计，钉死透射约定，不做反射光栅的符号反转、不做薄膜干涉、不算任何材料折射率与吸光度。

## HTTP API

```text
go run . serve -addr :8080
```

- `GET /api/health` 返回 `{"status":"ok"}`
- `POST /api/orders` 接受与 CLI 相同的光栅 JSON，返回各级次 θ、是否存在、角色散、FSR 以及可选的 R / Δλ
- `POST /api/checks` 运行波长增大 / 更密光栅 / 0 级 / ±m 对称 / R 线性 五条交叉规则

非法输入返回 HTTP 422。

## CLI 用法

```text
go run . orders example/600lpmm.json
```

打印内置算例（600 线/mm、λ=550 nm、正入射）从 −64 到 +64 级中每个级次的 sinθ、θ、是否存在、角色散与 FSR。±1 级落在 ±arcsin(λ/d) = ±19.27°，m=±3 接近掠射角 81.89°，m=±4 因 sin 超 1 不存在；0 级沿入射方向直行。

其他子命令与选项：

```text
go run . orders example/600lpmm.json --limit 5      # 只扫 −5..+5 级
go run . orders example/600lpmm.json --existing     # 只打印存在的级次
go run . orders example/600lpmm.json --json         # 结构化 JSON 输出
go run . orders example/coarse-300lpmm.json         # 300 线/mm + N=2000：多出 R 与 Δλ 列
go run . checks example/600lpmm.json                # 交叉规则自检
go run . density example/600lpmm.json               # 打印 d 与线密度的倒数关系
go run . serve [-addr :8080]                        # 启动 HTTP API
go run . help
```

算例：`example/600lpmm.json`（600 线/mm、550 nm、正入射，主算例）、`example/coarse-300lpmm.json`（300 线/mm、550 nm、正入射、N=2000，展示分辨本领列）。

JSON 输入字段（未知字段会被拒绝）：

```json
{
  "grooves_per_mm": 600,
  "d_nm": 1666.67,
  "wavelength_nm": 550,
  "incident_angle_deg": 0,
  "slits": 2000
}
```

槽距与线密度二选一即可；两个都写时必须互为倒数。波长单位 nm，角度单位度（从法线计，正方向与正级次同侧）。

## 关键约定

- **光栅方程**：mλ = d(sinθm − sinθi)，θ 一律从法线计，透射约定钉死。把方程写成加号（mλ = d(sinθm + sinθi)）会破坏正入射下 ±m 的反对称性与色散符号。
- **存在性**：|sinθm|≤1 才判该级次存在；arcsin 的自变量先做定义域检查，越界级次打印「—」而不是 NaN。λ/d 超过 1 的过密光栅只剩 0 级。
- **d 与线密度**：d（nm）与线密度（线/mm）是精确倒数，任何路径都不会绕过这条换算。
- **交叉规则**：波长增大 → 同级 |θ| 增大；d 减小（更密）→ 同级 |θ| 增大直到级次消失；正入射时 0 级 θ=0；±m 角度反对称；R=mN 随 N 线性增长。
- **有限缝包络**：主峰半宽 Δθ = λ/(N d cosθm)，经角色散换回的 Rayleigh Δλ 必须等于 λ/|m|N；占空比 a/d = 1/q 时偶数级次（q=2）的单缝包络为零，但光栅方程仍判定该级次几何上存在。
- **非法输入**：d≤0、λ≤0、|sinθi|>1、显式 N<1、未知 JSON 字段、文件缺失一律在 stderr 打印明文并非零退出。

## 构建与测试

```text
go build ./...
go test ./...
```

纯标准库，无第三方依赖，无 cgo。

## 许可

MIT，见 [LICENSE](./LICENSE)。

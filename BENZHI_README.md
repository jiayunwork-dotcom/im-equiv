# im-equiv：三相感应电机 T 形等效电路核算（输入电流 / 电磁转矩 / 机械功率 / 最大转矩转差）

Go 命令行工具：给定定子电阻漏抗、励磁支路、转子归算电阻漏抗与转差，核算输入电流、电磁转矩、机械功率与最大转矩转差，并输出同步速/转速、气隙功率、铜耗、铁耗、效率等配套量。s∈(0,1] 电动、s<0 发电、s=0 同步（转矩为 0）。

## 构建 / 运行 / 测试

```text
go build ./...
go run . torque example/4pole-50hz.json
go test ./...
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

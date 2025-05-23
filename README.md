go init test

go mod tidy  # 下载依赖并清理 go.mod


go mod download  # 仅下载依赖
go mod verify    # 验证依赖完整性
go mod vendor    # 将依赖复制到 vendor/ 目录（可选）


go build -o myapp  # 生成可执行文件（Windows 会生成 .exe）

go run main.go  # 无需编译，直接运行

# build.ps1
$env:CGO_ENABLED="0"
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o  minio_system


$env:CGO_ENABLED="0"
$env:GOARCH="amd64"
$env:GOOS="windows"
go build -o  niuzheng.exe


GOTMPDIR  设置go run 临时目录






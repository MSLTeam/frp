SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build github.com/fatedier/frp/cmd/frps
pause
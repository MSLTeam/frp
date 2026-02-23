echo 构建 frpc web UI
cd web\frpc
npm install
npx vite build
cd ..\..

echo 构建 frps web UI  
cd web\frps
npm install
npx vite build
cd ..\..

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64

echo 构建 Go 程序
$env:CGO_ENABLED=0
go build -trimpath -ldflags "-s -w" -tags frpc -o bin/frpc.exe ./cmd/frpc
go build -trimpath -ldflags "-s -w" -tags frps -o bin/frps.exe ./cmd/frps

pause
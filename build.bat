@echo off
echo Building Go binaries...

:: Windows
set GOOS=windows
set GOARCH=amd64
go build -o mocli.exe

:: Linux
set GOOS=linux
set GOARCH=amd64
go build -o mocli

echo Done!
pause

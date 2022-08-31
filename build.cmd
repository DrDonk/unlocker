@echo off

echo Building release executables
if [%1]==[] (
  echo Product version not found
  exit /b 1
)

pushd .\commands\check
echo Building check
go-winres make --arch amd64 --product-version %1 --file-version %1
set GOOS=windows
set GOARCH=amd64
go build -o ..\..\dist\windows\check.exe
set GOOS=linux
set GOARCH=amd64
go build -o ..\..\dist\linux\check
del /q rsrc_windows_amd64.syso
popd

pushd .\commands\relock
echo Building relock
go-winres make --arch amd64 --product-version %1 --file-version %1
set GOOS=windows
set GOARCH=amd64
go build -o ..\..\dist\windows\relock.exe
set GOOS=linux
set GOARCH=amd64
go build -o ..\..\dist\linux\relock
del /q rsrc_windows_amd64.syso
popd

pushd .\commands\unlock
echo Building unlock
go-winres make --arch amd64 --product-version %1 --file-version %1
set GOOS=windows
set GOARCH=amd64
go build -o ..\..\dist\windows\unlock.exe
set GOOS=linux
set GOARCH=amd64
go build -o ..\..\dist\linux\unlock
del /q rsrc_windows_amd64.syso
popd

xcopy /R /Y LICENSE dist\
xcopy /R /Y *.md dist\
xcopy /E /F /I /R /Y ISO dist\ISO
xcopy /E /F /I /R /Y ISO dist\templates

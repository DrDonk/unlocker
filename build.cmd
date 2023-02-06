@echo off

echo Building release executables
if [%1]==[] (
  echo Product version not found: x.y.z (e.g. 1.2.3)
  exit /b 1
)

pushd .\commands\check
echo Building check
go-winres make --arch amd64 --product-version %1 --file-version %1
set GOOS=windows
set GOARCH=amd64
go build -o ..\..\build\windows\check.exe
set GOOS=linux
set GOARCH=amd64
go build -o ..\..\build\linux\check
del /q rsrc_windows_amd64.syso
popd

pushd .\commands\relock
echo Building relock
go-winres make --arch amd64 --product-version %1 --file-version %1
set GOOS=windows
set GOARCH=amd64
go build -o ..\..\build\windows\relock.exe
set GOOS=linux
set GOARCH=amd64
go build -o ..\..\build\linux\relock
del /q rsrc_windows_amd64.syso
popd

pushd .\commands\dumpsmc
echo Building dumpsmc
go-winres make --arch amd64 --product-version %1 --file-version %1
set GOOS=windows
set GOARCH=amd64
go build -o ..\..\build\windows\dumpsmc.exe
set GOOS=linux
set GOARCH=amd64
go build -o ..\..\build\linux\dumpsmc
del /q rsrc_windows_amd64.syso
popd

pushd .\commands\hostcaps
echo Building hostcaps
go-winres make --arch amd64 --product-version %1 --file-version %1
set GOOS=windows
set GOARCH=amd64
go build -o ..\..\build\windows\hostcaps.exe
set GOOS=linux
set GOARCH=amd64
go build -o ..\..\build\linux\hostcaps
del /q rsrc_windows_amd64.syso
popd

xcopy /R /Y LICENSE .\build\
xcopy /R /Y *.md .\build\
xcopy /R /Y cpuid\linux\cpuid .\build\linux\cpuid
xcopy /R /Y cpuid\windows\cpuid.exe .\build\windows\cpuid.exe
xcopy /E /F /I /R /Y ISO .\build\ISO
xcopy /E /F /I /R /Y ISO .\'dist\templates

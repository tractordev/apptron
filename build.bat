@echo off

set script_path=%~dp0%
set project_root=%script_path%

:: Folders
set build_folder=%project_root%\lib\hostbridge
set lib_folder=%project_root%\lib

:: Paths
set go_path="C:\Program Files\Go\bin"
:: make sure gcc.exe is in your path (for CGO build)
set gcc_path="C:\ProgramData\chocolatey\bin\"
set PATH=%PATH%;%gcc_path%

:: 32-Bit mode
::set GOOS=windows
::set GOARCH=386

pushd %project_root%
  ::%go_path%\go.exe build -tags pkg -o ./debug-pkg.exe ./cmd/debug

  ::%go_path%\go.exe run ./schema/main/main.go
  %go_path%\go.exe build -tags cmd -o ./debug-cmd.exe ./cmd/debug

  IF %errorlevel% NEQ 0 (popd && goto end)

  .\debug-cmd.exe
popd


:end
exit /B %errorlevel%

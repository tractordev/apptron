@echo off

set script_path=%~dp0%
set project_root=%script_path%

:: Folders
set build_folder=%project_root%\lib\hostbridge
set lib_folder=%project_root%\lib

:: Paths
set cargo_exe=%USERPROFILE%\.cargo\bin
set go_path="C:\Program Files\Go\bin"
:: make sure gcc.exe is in your path (for CGO build)
set gcc_path="C:\ProgramData\chocolatey\bin\"
set PATH=%PATH%;%gcc_path%

:: Build
pushd %build_folder%
  %cargo_exe%\cargo.exe build --release
popd

:: NOTE(nick): these need to be in the same directory as the final executable
copy /y %build_folder%\target\release\hostbridge.dll %project_root%
copy /y %build_folder%\target\release\hostbridge.lib %project_root%

set CGO_LDFLAGS=%project_root%\hostbridge.dll

pushd %project_root%

  %go_path%\go.exe build -a -o ./ffi-debug.exe ./cmd/ffi-debug/main_static.go
popd

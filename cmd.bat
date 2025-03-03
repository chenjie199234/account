@echo off
REM      Warning!!!!!!!!!!!This file is readonly!Don't modify this file!

for /f "tokens=2 delims=:." %%a in ('chcp') do set codepage=%%a

cd %~dp0

where /q git.exe
if %errorlevel% == 1 (
	echo "missing dependence: git"
	exit /b 1
)

where /q go.exe
if %errorlevel% == 1 (
	echo "missing dependence: golang"
	exit /b 1
)

where /q protoc.exe
if %errorlevel% == 1 (
	echo "missing dependence: protoc"
	exit /b 1
)

where /q protoc-gen-go.exe
if %errorlevel% == 1 (
	echo "missing dependence: protoc-gen-go"
	exit /b 1
)

where /q codegen.exe
if %errorlevel% == 1 (
	echo "missing dependence: codegen"
	exit /b 1
)

if "%1" == "" (
	goto :help
)
if %1 == "" (
	goto :help
)
if %1 == "h" (
	goto :help
)
if "%1" == "h" (
	goto :help
)
if %1 == "-h" (
	goto :help
)
if "%1" == "-h" (
	goto :help
)
if %1 == "help" (
	goto :help
)
if "%1" == "help" (
	goto :help
)
if %1 == "-help" (
	goto :help
)
if "%1" == "-help" (
	goto :help
)
if %1 == "pb" (
	goto :pb
)
if "%1" == "pb" (
	goto :pb
)
if %1 == "kube" (
	goto :kube
)
if "%1" ==  "kube" (
	goto :kube
)
if %1 == "html" (
	goto :html
)
if "%1" ==  "html" (
	goto :html
)
if %1 == "sub" (
	if "%2" == "" (
		goto :help
	)
	if %2 == "" (
		goto :help
	)
	goto :sub
)
if "%1" == "sub" (
	if "%2" == "" (
		goto :help
	)
	if %2 == "" (
		goto :help
	)
	goto :sub
)

goto :help

:pb
	del >nul 2>nul .\api\*.pb.go
	del >nul 2>nul .\api\*.md
	del >nul 2>nul .\api\*.ts
	go mod tidy
	if %errorlevel% == 1 (
		echo "go mod tidy failed"
		exit /b 1
	)
	call :update
	for /f %%a in ('go list -m -f "{{.Dir}}" github.com/chenjie199234/Corelib') do set corelib=%%a
	protoc -I ./ -I %corelib% --go_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I %corelib% --go-pbex_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I %corelib% --go-cgrpc_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I %corelib% --go-crpc_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I %corelib% --go-web_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I %corelib% --browser_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I %corelib% --markdown_out=paths=source_relative:. ./api/*.proto
	go mod tidy
goto :end

:kube
	go mod tidy
	if %errorlevel% == 1 (
		echo "go mod tidy failed"
		exit /b 1
	)
	call :update
	codegen -n account -p github.com/chenjie199234/account -kube
goto :end

:html
	go mod tidy
	if %errorlevel% == 1 (
		echo "go mod tidy failed"
		exit /b 1
	)
	call :update
	codegen -n account -p github.com/chenjie199234/account -html
goto :end

:sub
	go mod tidy
	if %errorlevel% == 1 (
		echo "go mod tidy failed"
		exit /b 1
	)
	call :update
	codegen -n account -p github.com/chenjie199234/account -sub %2
goto :end

:update
	chcp 65001 >NUL 2>&1
	echo "start update corelib tools"
	for /f %%a in ('go list -m -f "{{.Dir}}" github.com/chenjie199234/Corelib') do set corelib=%%a
	set workdir=%cd%
	cd %corelib%
	go install ./...
	cd %workdir%
	echo "update corelib tools success"
	pause
	chcp %codepage% >NUL 2>&1
goto :eof

:help
	echo cmd.bat - every thing you need
	echo           please install git
	echo           please install golang(1.24+)
	echo           please install protoc           (github.com/protocolbuffers/protobuf)
	echo           please install protoc-gen-go    (github.com/protocolbuffers/protobuf-go)
	echo           please install codegen          (github.com/chenjie199234/Corelib)
	echo.
	echo Usage:
	echo    ./cmd.bat ^<option^>
	echo.
	echo Options:
	echo    pb                        Generate the proto in this program.
	echo    sub ^<sub service name^>    Create a new sub service.
	echo    kube                      Update kubernetes config.
	echo    html                      Create html template.
	echo    h/-h/help/-help/--help    Show this message.
:end
exit /b 0

#!/bin/sh
#      Warning!!!!!!!!!!!This file is readonly!Don't modify this file!

cd "$(dirname "$0")" || exit 1

help() {
	echo "cmd.sh — every thing you need"
	echo "         please install git"
	echo "         please install golang(1.26.1+)"
	echo "         please install protoc           (github.com/protocolbuffers/protobuf)"
	echo "         please install protoc-gen-go    (github.com/protocolbuffers/protobuf-go)"
	echo "         please install codegen          (github.com/chenjie199234/Corelib)"
	echo ""
	echo "Usage:"
	echo "   ./cmd.sh <option>"
	echo ""
	echo "Options:"
	echo "   run                       go run with -ldflags."
	echo "   build                     go build with -ldflags."
	echo "   pb                        Generate the proto in this program."
	echo "   sub <sub service name>    Create a new sub service."
	echo "   kube                      Update kubernetes config."
	echo "   html                      Create html template."
	echo "   h/-h/help/-help/--help    Show this message."
}

run(){
	dt=$(date -u '+%Y-%m-%d %H:%M:%S')
	go run -ldflags="-X 'main.version=${dt}'" main.go
}

build(){
	dt=$(date -u '+%Y-%m-%d %H:%M:%S')
	go build -ldflags="-X 'main.version=${dt}'" main.go
}

pb() {
	rm -f ./api/*.pb.go
	rm -f ./api/*.md
	rm -f ./api/*.ts
	if ! go mod tidy;then
		echo "go mod tidy failed"
		exit 1
	fi
	update
	corelib=$(go list -m -f "{{.Dir}}" github.com/chenjie199234/Corelib)
	protoc -I ./ -I "$corelib" --go_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I "$corelib" --go-pbex_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I "$corelib" --go-cgrpc_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I "$corelib" --go-crpc_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I "$corelib" --go-web_out=paths=source_relative:. ./api/*.proto
	protoc -I ./ -I "$corelib" --browser_out=outdir=api:. ./api/*.proto
	go mod tidy
}

sub() {
	if ! go mod tidy;then
		echo "go mod tidy failed"
		exit 1
	fi
	update
	codegen -n account -p github.com/chenjie199234/account -sub "$1"
}

kube() {
	if ! go mod tidy;then
		echo "go mod tidy failed"
		exit 1
	fi
	update
	codegen -n account -p github.com/chenjie199234/account -kube
}

html() {
	if ! go mod tidy;then
		echo "go mod tidy failed"
		exit 1
	fi
	update
	codegen -n account -p github.com/chenjie199234/account -html
}

update() {
	corelib=$(go list -m -f "{{.Dir}}" github.com/chenjie199234/Corelib)
	workdir=$(pwd)
	cd "$corelib" || exit 1
	go install ./...
	cd "$workdir" || exit 1
}

if ! command -v git >/dev/null 2>&1;then
	echo "missing dependence: git"
	exit 1
fi

if ! command -v go >/dev/null 2>&1;then
	echo "missing dependence: golang"
	exit 1
fi

if ! command -v protoc >/dev/null 2>&1;then
	echo "missing dependence: protoc"
	exit 1
fi

if ! command -v protoc-gen-go >/dev/null 2>&1;then
	echo "missing dependence: protoc-gen-go"
	exit 1
fi

if ! command -v codegen >/dev/null 2>&1;then
	echo "missing dependence: codegen"
	exit 1
fi

if [ $# -eq 0 ] || [ "$1" = "h" ] || [ "$1" = "help" ] || [ "$1" = "-h" ] || [ "$1" = "-help" ] || [ "$1" = "--help" ]; then
	help
	exit 0
fi

if [ "$1" = "run" ];then
	run
	exit 0
fi

if [ "$1" = "build" ];then
	build
	exit 0
fi

if [ "$1" = "pb" ];then
	pb
	exit 0
fi

if [ "$1" = "kube" ];then
	kube
	exit 0
fi

if [ "$1" = "html" ];then
	html
	exit 0
fi

if [ $# -eq 2 ] && [ "$1" = "sub" ];then
	sub "$2"
	exit 0
fi

help

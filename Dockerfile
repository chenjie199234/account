FROM golang:1.24.1 as builder
ENV GOSUMDB='off' \
	GOOS='linux' \
	GOARCH='amd64' \
	CGO_ENABLED=0

RUN mkdir /code
ADD . /code
WORKDIR /code
RUN echo "start build" && go mod tidy && go build -o main && echo "end build"

FROM debian:buster
RUN apt-get update && apt-get install -y ca-certificates curl inetutils-telnet inetutils-ping inetutils-traceroute dnsutils iproute2 procps net-tools neovim && mkdir /root/app
WORKDIR /root/app
EXPOSE 6060 7000 8000 9000 10000
COPY --from=builder /code/main /code/AppConfig.json /code/SourceConfig.json ./
ENTRYPOINT ["./main"]

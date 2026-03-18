FROM golang:1.26.1 as builder
ENV GOSUMDB='off' \
	GOOS='linux' \
	GOARCH='amd64' \
	CGO_ENABLED=0

RUN mkdir /code
ADD . /code
WORKDIR /code
RUN echo "start build" && go mod tidy && go build -ldflags="-X 'main.version=$(date -u '+%Y-%m-%d %H:%M:%S')'" main.go && echo "end build"

FROM alpine:3.23.3
RUN mkdir /root/app
WORKDIR /root/app
EXPOSE 6060 7000 8000 9000 10000
COPY --from=builder /code/main /code/AppConfig.json /code/SourceConfig.json ./
ENTRYPOINT ["./main"]

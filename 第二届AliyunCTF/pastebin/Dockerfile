FROM golang:1.22
ENV key test
ENV FLAG aliyunctf{test}

COPY ./src /usr/src/app
WORKDIR /usr/src/app

RUN go env -w  GOPROXY=https://goproxy.cn,direct && go mod download && go mod verify

RUN go build -v -o pastebin

CMD ["/usr/src/app/pastebin"]
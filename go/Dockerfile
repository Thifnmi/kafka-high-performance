FROM golang:1.22.5 AS build_base
LABEL maintainer="thifnmi <tuthin2k@gmail.com>"

RUN apt-get update && apt-get install -y git pkg-config

# stage 2
FROM build_base AS build_go

ENV GO111MODULE=on

WORKDIR $GOPATH/router
COPY go.mod .
COPY go.sum .
RUN go mod download

# stage 3
FROM build_go AS server_builder

ENV GO111MODULE=on

COPY . .
# RUN go build -o /bin/router ./main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/router ./main.go

FROM debian AS router

ENV TZ 'Asia/Ho_Chi_Minh'
RUN echo $TZ > /etc/timezone && \
    apt-get update && apt-get install -y tzdata && \
    rm /etc/localtime && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    apt-get clean


COPY --from=server_builder /bin/router /bin/router


CMD ["/bin/router"]



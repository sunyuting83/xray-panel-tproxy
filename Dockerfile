FROM node:latest as buildJS
RUN mkdir -p /usr/src/app
COPY ./xpanel-web /usr/src/app/
WORKDIR /usr/src/app

# RUN mkdir /lib64
# RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
# RUN ln -s /usr/lib/libc.so /usr/lib/libresolv.so.2

RUN npm config set registry https://registry.npmmirror.com/
RUN npm install
RUN npm run build

FROM golang:alpine AS builder
RUN mkdir /app
RUN mkdir -p /app/static
COPY . /app/
COPY --from=buildJS /usr/src/app/build /app/static/
WORKDIR /app

RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache musl-dev
RUN mkdir /lib64
RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN ln -s /usr/lib/libc.so /usr/lib/libresolv.so.2

ADD https://github.com/upx/upx/releases/download/v4.2.1/upx-4.2.1-amd64_linux.tar.xz /usr/local
RUN tar -xf /usr/local/upx-4.2.1-amd64_linux.tar.xz -C /usr/local && mv /usr/local/upx-4.2.1-amd64_linux/upx /bin/upx && \
    chmod a+x /bin/upx

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN go build -tags=jsoniter -trimpath -ldflags "-s -w -buildid=" -o server main.go
RUN upx --lzma server

FROM alpine:latest

RUN mkdir -p /xpanel
WORKDIR /xpanel

COPY --from=builder /app/server /xpanel/
COPY --from=builder /app/entrypoint.sh /usr/bin/
COPY --from=builder /app/run.sh /xpanel/run.sh
COPY --from=builder /app/data /xpanel/data/
COPY --from=builder /app/template /xpanel/template/
COPY --from=builder /app/static /xpanel/static/

RUN chmod a+x /xpanel/run.sh

RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache \
 ca-certificates  \
 iptables \
 musl-dev
RUN mkdir /lib64
RUN ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN ln -s /usr/lib/libc.so /usr/lib/libresolv.so.2
RUN rm -rf /var/cache/apk/*
RUN chmod a+x /usr/bin/entrypoint.sh

ENTRYPOINT ["/usr/bin/entrypoint.sh"]
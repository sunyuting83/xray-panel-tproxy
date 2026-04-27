# --- 第一阶段：前端构建 ---
FROM --platform=$BUILDPLATFORM node:latest as buildJS
WORKDIR /usr/src/app
COPY ./xpanel-web ./
# RUN npm config set registry https://registry.npmmirror.com/ && npm install && npm run build
RUN npm install && npm run build

# --- 第二阶段：后端编译 ---
FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY . .
COPY --from=buildJS /usr/src/app/build ./static/

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    apk add --no-cache musl-dev

# 关键：根据 GitHub Action 传进来的架构参数进行交叉编译
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -tags=jsoniter -trimpath -ldflags "-s -w -buildid=" -o server main.go

# --- 第三阶段：UPX 压缩 (仅限 amd64，arm64 建议跳过或使用对应版本) ---
FROM --platform=$BUILDPLATFORM gruebel/upx:latest AS upx-processor
ARG TARGETARCH
COPY --from=builder /app/server /server
# UPX 在某些架构下可能有兼容性问题，这里加个判断，如果是 arm64 则可选压缩或不压
RUN upx --lzma /server

# --- 第四阶段：最终运行环境 ---
FROM alpine:latest
ARG TARGETARCH

RUN mkdir -p /xpanel
WORKDIR /xpanel

COPY --from=upx-processor /server /xpanel/server
COPY --from=builder /app/entrypoint.sh /usr/bin/
COPY --from=builder /app/run.sh /xpanel/run.sh
COPY --from=builder /app/data /xpanel/data/
COPY --from=builder /app/template /xpanel/template/
COPY --from=builder /app/static /xpanel/static/

RUN set -eux && \
    apk add --no-cache ca-certificates iptables ip6tables && \
    rm -rf /var/cache/apk/*

# 自动适配架构的软链接
RUN mkdir -p /lib64 && \
    if [ "$TARGETARCH" = "amd64" ]; then ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2; \
    elif [ "$TARGETARCH" = "arm64" ]; then ln -s /lib/libc.musl-aarch64.so.1 /lib64/ld-linux-aarch64.so.2; fi

RUN chmod a+x /xpanel/run.sh /usr/bin/entrypoint.sh
ENTRYPOINT ["/usr/bin/entrypoint.sh"]
FROM registry.cn-shanghai.aliyuncs.com/bitypes/golang:1.23.2-alpine3.20 AS builder
WORKDIR /app
ARG app
ENV app=${app}

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,https://proxy.golang.com.cn,direct
RUN go mod tidy
RUN go build -buildvcs=false -ldflags "-s -w" -trimpath -o /app/${app}

FROM registry.cn-shanghai.aliyuncs.com/bitypes/alpine:3.20
WORKDIR /app
ARG app
ARG imgVer
ENV app=${app}
ENV imgVer=${imgVer}
ENV TZ=Asia/Shanghai

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
  apk update && \
  apk upgrade && \
  apk add ca-certificates && update-ca-certificates && \
  apk add --update tzdata && \
  rm -rf /var/cache/apk/*

COPY --from=builder /app/${app} /app/${app}
EXPOSE 8080

ENTRYPOINT sh -c 'exec /app/${app}'  
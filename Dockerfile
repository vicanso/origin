FROM node:12-alpine as webbuilder

COPY . /origin
RUN cd /origin/web \
  && yarn \
  && yarn build \
  && rm -rf node_module

FROM golang:1.14-alpine as builder

COPY --from=webbuilder /origin /origin

RUN apk update \
  && apk add git make \
  && go get -u github.com/gobuffalo/packr/v2/packr2 \
  && cd /origin \
  && make build

FROM alpine 

EXPOSE 7001

# tzdata 安装所有时区配置或可根据需要只添加所需时区

RUN addgroup -g 1000 go \
  && adduser -u 1000 -G go -s /bin/sh -D go \
  && apk add --no-cache ca-certificates tzdata

COPY --from=builder /origin/origin /usr/local/bin/origin
COPY --from=builder /origin/entrypoint.sh /entrypoint.sh

USER go

WORKDIR /home/go

HEALTHCHECK --timeout=10s --interval=10s CMD [ "wget", "http://127.0.0.1:7001/ping", "-q", "-O", "-"]

CMD ["origin"]

ENTRYPOINT ["/entrypoint.sh"]

FROM golang:1.9 as builder
WORKDIR /go/src/github.com/xwjdsh/2048-ai
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o 2048-ai .

FROM alpine:latest  
LABEL maintainer="iwendellsun@gmail.com"
RUN apk --no-cache add ca-certificates tzdata \
			&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
			&& echo "Asia/Shanghai" >  /etc/timezone \
			&& apk del tzdata
WORKDIR /root/2048-ai
COPY --from=builder /go/src/github.com/xwjdsh/2048-ai/2048-ai .
COPY --from=builder /go/src/github.com/xwjdsh/2048-ai/2048 ./2048
EXPOSE 8080
CMD ["./2048-ai"]

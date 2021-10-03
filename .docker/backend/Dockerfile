#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache build-base
WORKDIR /go/src/app
COPY backend/. .

RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./main.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /usr/bin/yclogui/app
COPY --from=builder /go/src/app/.env /usr/bin/yclogui/.env
WORKDIR /usr/bin/yclogui/
ENTRYPOINT /usr/bin/yclogui/app
LABEL Name=yclogui Version=1.0.0
EXPOSE 1323

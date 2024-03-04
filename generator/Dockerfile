FROM golang:1.21.4-alpine3.18 as builder

WORKDIR /opt/build/

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY ./ ./

RUN go mod download && go mod verify

RUN go build -v -o /opt/app/ ./main.go



FROM alpine:3.19.0

COPY  --from=builder /opt/app /opt/app

WORKDIR /opt/app/


# non root user
RUN adduser -D golang
RUN chown -R golang:golang /opt/app

USER golang

CMD ./main
FROM golang:1.20.0-bullseye

WORKDIR /opt/build/

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY ./ ./

RUN go mod download && go mod verify

RUN go build -v -o /opt/app/ ./main.go

WORKDIR /opt/app/

CMD ./main

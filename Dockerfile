# Build Container
FROM golang:1.16 as builder
WORKDIR /go/src/github.com/swallowarc/tictactoe_battle_backend
COPY . .

# Set Environment Variable
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Build
RUN make

# runtime image
FROM alpine:3.13.5
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/swallowarc/tictactoe_battle_backend/bin /bin
ENTRYPOINT ["/bin/tictactoe_battle_backend"]

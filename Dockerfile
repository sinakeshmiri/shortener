FROM golang:1.20 AS builder

COPY . /src
WORKDIR /src
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build  -o /app /src/cmd/main.go

FROM scratch
COPY --from=builder /app ./
EXPOSE 3000
EXPOSE 2112
ENTRYPOINT ["./app"]
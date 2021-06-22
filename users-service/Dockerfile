FROM golang:latest as builder
WORKDIR /app/backend
ADD . /app/backend
RUN ls -a
RUN go test -v ./... -v -short -cover -coverprofile cover.out
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /userssrv /app/backend/main.go

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /userssrv ./
RUN chmod +x ./userssrv
ENTRYPOINT ["./userssrv Prod"]
EXPOSE 8080
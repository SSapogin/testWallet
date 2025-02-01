FROM golang:latest as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o run cmd/*.go


FROM scratch

COPY --from=builder /app/run /run
COPY .env ./

EXPOSE 8080

CMD ["/run"]
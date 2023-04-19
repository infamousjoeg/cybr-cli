FROM golang as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/cybr .

FROM ubuntu:22.10
COPY --from=builder /app/cybr /app/
RUN useradd -ms /bin/bash cybr && \
    chmod 777 /home/cybr && \
    mkdir -p /home/cybr/.cybr && \
    chown -R cybr /home/cybr/.cybr
USER cybr
ENV PATH="/app:${PATH}"
ENTRYPOINT ["/app/cybr"]

FROM golang:1-bullseye

ENV GO111MODULE=on

WORKDIR /app

ENV PATH="/app:${PATH}"

RUN useradd -ms /bin/bash cybr && \
    chmod 777 /home/cybr && \
    mkdir -p /home/cybr/.cybr && \
    chown -R cybr /home/cybr/.cybr 

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/cybr .

USER cybr
ENTRYPOINT ["/app/cybr"]
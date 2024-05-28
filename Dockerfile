FROM golang:1.18
RUN mkdir /home/go
WORKDIR /home/go
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY login.go ./
RUN go build -o login login.go
ENTRYPOINT ["./login"]
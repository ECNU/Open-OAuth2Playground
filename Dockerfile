FROM golang:1.20

WORKDIR /app/Open-OAuth2Playground

COPY ./go.mod .
COPY ./go.sum .

ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . .
RUN go build -o OAuth2Playground .

EXPOSE 80

CMD ["./OAuth2Playground"]
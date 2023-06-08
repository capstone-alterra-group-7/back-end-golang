FROM golang:alpine

RUN apk update && apk add git

RUN apk add --no-cache nginx

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary

# EXPOSE 443

ENTRYPOINT [ "./binary" ]

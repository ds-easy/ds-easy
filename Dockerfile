FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache bash git make sqlite

COPY . .

# Install dependencies
RUN go mod download

CMD ["make", "watch"]

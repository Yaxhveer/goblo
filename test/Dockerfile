FROM golang:1.22.3

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /test

EXPOSE 3331 3332 3333 3334 3335 3336 3337

CMD ["/test"]
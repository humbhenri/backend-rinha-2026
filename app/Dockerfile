FROM golang:1.26.3

COPY go.mod ./
#RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /rinha

CMD ["/rinha"]
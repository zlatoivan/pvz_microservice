FROM golang:1.22
LABEL authors="ivan"

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/pressly/goose/cmd/goose

COPY . .

RUN go build -o main cmd/server/main.go

EXPOSE 9000
EXPOSE 9001

CMD ["./main"]



#WORKDIR /Homework
#RUN go get github.com/pressly/goose/cmd/goose
#ENV CONFIG_PATH=config/config.yaml
#RUN go get -d -v ./...
#RUN go install -v ./...
#RUN CGO_ENABLED=0 GOOS=linux go build -o /main

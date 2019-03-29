FROM golang

WORKDIR /go/src/filmtracker
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["filmtracker"]
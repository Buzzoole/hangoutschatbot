FROM golang:1.10.0-stretch


COPY . /go/src/github.com/AppsterdamMilan/hangoutschatbot

# install dep and download dependencies
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN cd /go/src/github.com/AppsterdamMilan/hangoutschatbot && dep ensure

# build main
RUN go build -o /app github.com/AppsterdamMilan/hangoutschatbot

WORKDIR /

EXPOSE 8080
ENTRYPOINT ["/app"]
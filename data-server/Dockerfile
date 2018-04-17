FROM golang:alpine

WORKDIR /go/src/github.com/fcrisciani/data
COPY . ./

RUN mkdir -p /go/bin/
RUN go install

CMD ["ls", "/go/bin/"]

FROM golang:alpine

COPY --from=0 /go/bin/data /usr/local/bin/data-server
ARG SERVER_PORT=10000
ENV PORT=${SERVER_PORT}
ENTRYPOINT ["/usr/local/bin/data-server"]

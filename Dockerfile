FROM golang:1.13-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go mod download

RUN go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/app /

ENV PORT 8080

CMD ["/app"]


FROM golang:1.14 as build-env

# Docs on how to use container: https://github.com/levibostian/purslane/

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...

RUN go build -o /go/bin/purslane

RUN touch /.purslane.yaml

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/purslane /purslane
COPY --from=build-env /.purslane.yaml /root/.purslane.yaml
CMD ["/purslane", "run"]
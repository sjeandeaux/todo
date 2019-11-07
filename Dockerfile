FROM golang:1.13.4-alpine as base

RUN apk --no-cache add ca-certificates make git g++ protobuf protobuf-dev

ENV GO111MODULE on

#get source
WORKDIR /go/src/github.com/sjeandeaux/todo

#copy the source
COPY . .

RUN make tools
RUN make dependencies
RUN make generate

## test
FROM base AS test
RUN make test

#Build the application
FROM base AS build
RUN make build

FROM scratch AS release

ARG BUILD_VERSION=undefined
ARG BUILD_DATE=undefined

#http://label-schema.org/rc1/
LABEL "maintainer"="stephane.jeandeaux@gmail.com" \
      "org.label-schema.vendor"="sjeandeaux" \
      "org.label-schema.schema-version"="1.0.0-rc.1" \
      "org.label-schema.applications.todod.version"=${BUILD_VERSION} \
      "org.label-schema.build-date"=${BUILD_DATE}

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/sjeandeaux/todo/target/todod /todod
COPY --from=build /go/bin/grpc-health-probe /grpc-health-probe

##TODO add HEALTHCHECK
EXPOSE 8080
ENTRYPOINT ["/todod"]
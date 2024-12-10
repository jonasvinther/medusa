FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG VERSION
ARG TARGETARCH

RUN apk update && apk add --no-cache git

ENV GO111MODULE=on

WORKDIR /app

ADD . .

RUN go mod download
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build \
    -ldflags="-X 'github.com/jonasvinther/medusa/cmd.Version=${VERSION}'" \
    -o /go/bin/medusa
RUN apk add ca-certificates && update-ca-certificates

RUN adduser -S scratchuser
RUN chown scratchuser /go/bin/medusa

FROM scratch
COPY --from=builder /go/bin/medusa /medusa
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/share/ca-certificates/mozilla/* /etc/ssl/certs/

USER scratchuser
ENTRYPOINT ["/medusa"]

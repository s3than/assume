FROM golang:alpine as builder
MAINTAINER Tim Colbert <admin@tcolbert.net>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

COPY . /app

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/s3than/assume \
	&& make static \
	&& mv assume /usr/bin/assume \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/assume /usr/bin/assume

ENTRYPOINT [ "assume" ]
CMD [ "--help" ]
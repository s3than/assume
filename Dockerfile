FROM golang:alpine as builder

WORKDIR /app

COPY . /app

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& make static \
	&& mv assume /usr/bin/assume \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/assume /usr/bin/assume

ENTRYPOINT [ "assume" ]
CMD [ "--help" ]
FROM alpine:latest
RUN apk add --no-cache moreutils
ADD https://github.com/orls/envtemplate/releases/download/0.0.3/envtemplate /usr/bin/envtemplate
RUN chmod +x /usr/bin/envtemplate
COPY entrypoint.sh /usr/bin/
ENTRYPOINT [ "entrypoint.sh" ]

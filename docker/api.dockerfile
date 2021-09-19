FROM alpine:latest

RUN apk add --update ca-certificates sqlite tzdata curl httpie && \
    rm -rf /var/cache/apk/* && \
    cp /usr/share/zoneinfo/Asia/Taipei /etc/localtime && \
    echo "Asia/Taipei" > /etc/timezone

COPY build/api /usr/local/bin/api

CMD /usr/local/bin/api

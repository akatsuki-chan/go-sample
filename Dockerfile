FROM alpine:3.7

ARG ENV_MODE=${ENV_MODE:-prod}

# HTTPSでエラーが出るときなど
RUN apk add --no-cache ca-certificates

WORKDIR /work

COPY ./build/sample /work
COPY ./config.${ENV_MODE}.yml /work/config.yml

CMD ["./sample"]
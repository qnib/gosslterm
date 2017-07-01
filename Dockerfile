ARG DOCKER_REG=docker.io
FROM ${DOCKER_REG}/qnib/uplain-golang AS build

WORKDIR /usr/local/src/github.com/qnib/gosslterm
COPY main.go ./main.go
COPY vendor/ vendor/
RUN govendor install


FROM ubuntu AS ssl

WORKDIR /opt/qnib/ssl/
RUN apt-get update \
 && apt-get install -y openssl
RUN openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem \
  -days 365 -nodes -subj "/C=DE/ST=Berlin/L=Berlin/O=QNIB Solutions/OU=IT Department/CN=qnib.org"

## Build final image

FROM ${DOCKER_REG}/qnib/uplain-init

COPY --from=build /usr/local/bin/gosslterm /usr/local/bin/
COPY --from=ssl /opt/qnib/ssl/cert.pem /opt/qnib/ssl/key.pem /opt/qnib/ssl/
ENV SKIP_ENTRYPOINTS=true \
    GOSSLTERM_CERT=/opt/qnib/ssl/cert.pem \
    GOSSLTERM_KEY=/opt/qnib/ssl/key.pem \
    GOSSLTERM_FRONTEND_ADDR=:8081 \
    GOSSLTERM_BACKEND_ADDR=127.0.0.1:8080
CMD ["gosslterm"]

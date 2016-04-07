FROM deeva/manager:latest

USER root

COPY . /home/deeva/go/src/github.com/deevatech/manager/
RUN chown -R deeva:deeva /home/deeva/go \
  && mkdir -p /etc/docker/certs
VOLUME /etc/docker/certs

USER deeva
ENV DOCKER_CERT_PATH=/etc/docker/certs/ \
    DOCKER_TLS_VERIFY=1 \
    GOPATH=/home/deeva/go/

WORKDIR /home/deeva/go/src/github.com/deevatech/manager/
RUN glide install \
  && go build -o manager

EXPOSE 8080
CMD ./manager


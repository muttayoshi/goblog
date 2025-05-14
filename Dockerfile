# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-goblog"
LABEL REPO="https://github.com/muttayoshi/goblog"

ENV PROJPATH=/go/src/github.com/muttayoshi/goblog

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/muttayoshi/goblog
WORKDIR /go/src/github.com/muttayoshi/goblog

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/muttayoshi/goblog"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/goblog/bin

WORKDIR /opt/goblog/bin

COPY --from=build-stage /go/src/github.com/muttayoshi/goblog/bin/goblog /opt/goblog/bin/
RUN chmod +x /opt/goblog/bin/goblog

# Create appuser
RUN adduser -D -g '' goblog
USER goblog

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/goblog/bin/goblog"]

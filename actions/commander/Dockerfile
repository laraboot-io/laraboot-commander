FROM golang:1.16.2-alpine3.13 AS build
WORKDIR /src
ENV BUILDER_WORKBENCH /workbench
COPY . .
RUN mkdir -p $BUILDER_WORKBENCH
RUN apk upgrade -U -a && \
	apk add --update docker openrc && \
	rc-update add docker boot && \
    apk add wget ca-certificates openssl-dev --update-cache && \
    update-ca-certificates && \
    go get github.com/markbates/pkger/cmd/pkger  && \
    apk add bash jq yq curl gcc g++ musl-dev && \
    wget -4 https://github.com/mikefarah/yq/releases/download/v4.2.0/yq_linux_amd64 -O /usr/bin/yq && \
    chmod +x /usr/bin/yq && \
    (curl -sSL "https://github.com/buildpacks/pack/releases/download/v0.19.0/pack-v0.19.0-linux.tgz" | tar -C /usr/local/bin/ --no-same-owner -xzv pack) && \
    wget -4 https://github.com/paketo-buildpacks/packit/releases/download/v0.8.0/jam-linux -O /usr/bin/jam
WORKDIR $BUILDER_WORKBENCH
ENTRYPOINT ["/src/scripts/entrypoint.sh"]
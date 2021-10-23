FROM golang:1.16.2-alpine3.13 AS build
ENV BUILDER_WORKBENCH /workbench
WORKDIR /src
COPY . .
RUN mkdir -p $BUILDER_WORKBENCH
RUN apk upgrade -U -a && \
	apk add --update docker openrc && \
	rc-update add docker boot && \
    apk add wget ca-certificates openssl-dev --update-cache && \
    update-ca-certificates && \
    apk add bash jq yq curl gcc g++ musl-dev && \
    (curl -sSL "https://github.com/buildpacks/pack/releases/download/v0.19.0/pack-v0.19.0-linux.tgz" | tar -C /usr/local/bin/ --no-same-owner -xzv pack) && \
    wget -4 https://github.com/paketo-buildpacks/packit/releases/download/v0.8.0/jam-linux -O /usr/bin/jam
RUN chmod +x /src/scripts/entrypoint.sh
ENTRYPOINT ["/src/scripts/entrypoint.sh"]
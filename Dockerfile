ARG BUILDER_IMG=chadgrant/base:golang-1.13.5-alpine
ARG RUNTIME_IMG=chadgrant/base:alpine-3.11.2

FROM $BUILDER_IMG AS builder

RUN install-deps make

WORKDIR /go/src/github.com/$VENDOR/$SERVICE/

COPY vendor ./vendor/ 
COPY go.mod go.sum ./
RUN go mod download

COPY Makefile main.go ./
COPY store ./store/

ARG VENDOR
ARG GROUP
ARG SERVICE
ARG SERVICE_FRIENDLY
ARG SERVICE_DESCRIPTION
ARG SERVICE_URL
ARG BUILD_HASH
ARG BUILD_BRANCH
ARG BUILD_USER
ARG BUILD_NUMBER
ARG BUILD_REPO
ARG BUILD_DATE
ENV CGO_ENABLED=0 \
    VENDOR="${VENDOR}" GROUP="${GROUP}" \
    SERVICE="${SERVICE}" SERVICE_FRIENDLY="${SERVICE_FRIENDLY}" SERVICE_DESCRIPTION="${SERVICE_DESCRIPTION}" SERVICE_URL="${SERVICE_URL}" \
    BUILD_HASH="${BUILD_HASH}" BUILD_BRANCH="${BUILD_BRANCH}" BUILD_USER="${BUILD_USER}" \
    BUILD_NUMBER="${BUILD_NUMBER}" BUILD_REPO="${BUILD_REPO}" BUILD_DATE="${BUILD_DATE}" 

RUN BINARY_NAME=goapp OUT_DIR=/go/bin/ make build

FROM $RUNTIME_IMG

RUN install-deps ca-certificates libc6-compat 

RUN addgroup -S app && \
    adduser -S app -G app
USER app

WORKDIR /app

COPY docs /app/docs/
COPY schema /app/schema/
COPY --from=builder /go/bin/goapp /app/

CMD ["/app/goapp"]

ARG VENDOR
ARG GROUP
ARG SERVICE
ARG SERVICE_FRIENDLY
ARG SERVICE_DESCRIPTION
ARG SERVICE_URL
ARG BUILD_HASH
ARG BUILD_BRANCH
ARG BUILD_USER
ARG BUILD_NUMBER
ARG BUILD_REPO
ARG BUILD_DATE

## http://label-schema.org/rc1/
LABEL org.label-schema.schema-version="1.0" \
    org.label-schema.vendor="${VENDOR}" \
    org.label-schema.build-group="${GROUP}" \
    org.label-schema.application-name="${SERVICE}" \
    org.label-schema.name="${SERVICE_FRIENDLY}" \
    org.label-schema.description="${SERVICE_DESCRIPTION}" \
    org.label-schema.url="${SERVICE_URL}" \
    org.label-schema.version="${BUILD_NUMBER}" \
    org.label-schema.build-user="${BUILD_USER}" \
    org.label-schema.build-date="${BUILD_DATE}" \
    org.label-schema.vcs-branch="${BUILD_BRANCH}" \
    org.label-schema.vcs-ref="${BUILD_HASH}" \
    org.label-schema.vcs-url="${BUILD_REPO}"
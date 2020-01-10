ARG builder_img
ARG runtime_img

FROM $builder_img AS builder

ARG application
ARG friendly
ARG description
ARG build_hash
ARG build_branch
ARG build_user
ARG build_number
ARG build_group
ARG repo_url
ARG vendor
ARG build_date
ENV CGO_ENABLED=0 APPLICATION=$application FRIENDLY=$friendly DESCRIPTION="${description}" BUILD_HASH=$build_hash BUILD_BRANCH=$build_branch BUILD_USER=$build_user BUILD_NUMER=$build_number BUILD_GROUP=$build_group REPO_URL=$repo_url

RUN install-deps make

WORKDIR /go/src/github.com/chadgrant/$application/

COPY vendor ./vendor/ 
COPY go.mod go.sum ./
RUN go mod download

COPY Makefile main.go ./
COPY store ./store/

RUN BINARY_NAME=goapp OUT_DIR=/go/bin/ make build

FROM $runtime_img
ARG application
ARG friendly
ARG description
ARG build_hash
ARG build_branch
ARG build_user
ARG build_number
ARG build_group
ARG repo_url
ARG vendor
ARG build_date

RUN install-deps ca-certificates libc6-compat 
RUN addgroup -S app && \
    adduser -S app -G app
USER app
WORKDIR /app
COPY docs /app/docs/
COPY --from=builder /go/bin/goapp /app/
CMD ["/app/goapp"]

## http://label-schema.org/rc1/
LABEL org.label-schema.schema-version="1.0" \
    org.label-schema.version="${build_number}" \
    org.label-schema.name="${friendly}" \
    org.label-schema.description="${description}" \
    org.label-schema.application-name="${application}" \
    org.label-schema.build-group="${build_group}" \
    org.label-schema.build-user="${build_user}" \
    org.label-schema.build-date="${build_date}" \
    org.label-schema.vcs-branch="${build_branch}" \
    org.label-schema.vcs-ref="${build_hash}" \
    org.label-schema.vcs-url="${repo_url}" \
    org.label-schema.url="${repo_url}" \
    org.label-schema.vendor="${vendor}"
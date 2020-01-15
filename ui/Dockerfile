ARG BUILDER_IMG="node:12.13.1-alpine"
ARG RUNTIME_IMG="nginx:1.16.0-alpine"

FROM $BUILDER_IMG as builder
ARG UI="typescript-react"

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY ${UI}/*.json /app/
RUN rm *-lock.json \
    && npm install \
    && npm install react-scripts@3.3.0 -g

COPY ${UI}/public public
COPY ${UI}/src src
RUN npm run build

FROM $RUNTIME_IMG
ARG SERVICE
ARG SERVICE_FRIENDLY
ARG SERVICE_DESCRIPTION
ARG SERVICE_URL
ARG BUILD_HASH
ARG BUILD_BRANCH
ARG BUILD_USER
ARG BUILD_NUMBER
ARG BUILD_GROUP
ARG BUILD_REPO
ARG VENDOR
ARG BUILD_DATE

COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 80
ENTRYPOINT ["/entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]

## http://label-schema.org/rc1/
LABEL org.label-schema.schema-version="1.0" \
    org.label-schema.version="${BUILD_NUMBER}" \
    org.label-schema.name="${SERVICE_FRIENDLY}" \
    org.label-schema.description="${SERVICE_DESCRIPTION}" \
    org.label-schema.application-name="${SERVICE}" \
    org.label-schema.build-group="${BUILD_GROUP}" \
    org.label-schema.build-user="${BUILD_USER}" \
    org.label-schema.build-date="${BUILD_DATE}" \
    org.label-schema.vcs-branch="${BUILD_BRANCH}" \
    org.label-schema.vcs-ref="${BUILD_HASH}" \
    org.label-schema.vcs-url="${BUILD_REPO}" \
    org.label-schema.url="${SERVICE_URL}" \
    org.label-schema.vendor="${VENDOR}"
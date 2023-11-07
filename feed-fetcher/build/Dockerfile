FROM golang:1.21 AS builder

ARG GOMOD_PATH=go.mod
ARG PROJECT_PATH=.
ARG OUTPUT_BINARY
ARG BUILD_DIR
ARG APP_PORT=8080

WORKDIR /app

COPY ${GOMOD_PATH} ./
RUN go mod download && go mod verify

COPY ${PROJECT_PATH} .

RUN mkdir -p $(dirname ${OUTPUT_BINARY}) && \
    CGO_ENABLED=0 go build -o ${OUTPUT_BINARY} ${BUILD_DIR}

FROM alpine:3.18.3

RUN apk add tzdata libc6-compat

WORKDIR /opt

ARG OUTPUT_BINARY
ARG APP_NAME
ARG APP_VERSION

ENV APP_NAME=${APP_NAME}
ENV APP_VERSION=${APP_VERSION}

COPY --from=builder /app/${OUTPUT_BINARY} binary
COPY --from=builder /app/config/docker.yaml ./config/docker.yaml

RUN adduser --disabled-password --gecos "" --home "$(pwd)" \
    --ingroup "users" --no-create-home --uid "888" "nonroot" \
 && chown -R -f 888:888 . \
 && chmod +x ./binary
USER 888

EXPOSE $APP_PORT

CMD ["./binary"]
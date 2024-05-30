ARG GOLANG_VERSION=alpine
FROM golang:${GOLANG_VERSION} AS base

WORKDIR /app

FROM base as builder

COPY --link go.* .

RUN go mod download

COPY --link . .

RUN go test ./...

RUN --mount=type=cache,target=/root/.cache/go-build \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o score.app ./main.go


FROM alpine:latest as runner

LABEL app="simility-score"
ARG ENVIRONMENT=development
ENV GO_ENV=${ENVIRONMENT}
# github related envs
ARG GITHUB_SHA
ENV GITHUB_SHA=${GITHUB_SHA}
ARG GITHUB_REF
ENV GITHUB_REF=${GITHUB_REF}
ARG GITHUB_REF_NAME
ENV GITHUB_REF_NAME=${GITHUB_REF_NAME}

RUN adduser -D -u 1000 appuser

COPY --link --from=builder --chown=1000:1000 /app/score.app .

RUN chmod +x ./score.app

USER appuser

CMD [ "/score.app"]
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


FROM scratch as runner

WORKDIR /app

LABEL app="simility-score"

COPY --link --from=builder /app/score.app /score.app

CMD [ "/score.app"]
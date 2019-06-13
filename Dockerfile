FROM golang:1.12-alpine AS build_base
WORKDIR /go/src/github.com/troyspencer/launch-wasm
RUN apk add bash git
ENV GO111MODULE=on
ENV GOOS=js
ENV GOARCH=wasm
COPY game/go.mod game/go.sum game/
RUN cd game && go mod download

FROM build_base AS game_builder
COPY game game
RUN cd game && go build -o main.wasm main.go 
RUN gzip -k ./game/main.wasm

FROM gcr.io/cloud-builders/npm AS react_builder
WORKDIR /go/src/github.com/troyspencer/launch-wasm
COPY react/package.json react/yarn.lock react/
RUN cd react && yarn --pure-lockfile
COPY react react
COPY --from=game_builder /go/src/github.com/troyspencer/launch-wasm/game/main.wasm /go/src/github.com/troyspencer/launch-wasm/game/main.wasm.gz ./react/static/
RUN cd react && yarn build

FROM gcr.io/cloud-builders/gcloud
WORKDIR /go/src/github.com/troyspencer/launch-wasm
COPY server server
COPY --from=react_builder /go/src/github.com/troyspencer/launch-wasm/server/dist ./server/dist/
RUN bash -c set -e && if [[ "$BRANCH_NAME" == "master" ]]; then gcloud app deploy server/app.yaml; else gcloud app deploy server/app.yaml --no-promote; fi
FROM golang:1.12-alpine AS build_base

RUN apk add bash git

WORKDIR /go/src/github.com/troyspencer/launch-wasm

ENV GO111MODULE=on
ENV GOOS=js
ENV GOARCH=wasm

COPY game/go.mod game/go.sum game/

RUN cd game && go mod download

FROM build_base AS game_builder

COPY game game

RUN cd game && go build -o ../react/static/main.wasm main.go 
RUN gzip -k react/static/main.wasm

FROM node:latest

COPY package.json yarn.lock ./
COPY react/package.json react/yarn.lock react/

RUN yarn --pure-lockfile
RUN yarn --cwd react --pure-lockfile

COPY react react

RUN cd react && yarn build
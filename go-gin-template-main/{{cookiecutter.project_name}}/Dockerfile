FROM docker-registry.fmobile.kz/golang:1.19-alpine

ARG GIT_TOKEN
WORKDIR /app

ENV TZ=Asia/Almaty

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN git config --global url."https://oauth2:$GIT_TOKEN@git-ffd.kz/".insteadOf "https://git-ffd.kz/"

RUN go env -w GOPRIVATE=git-ffd.kz/*

COPY go.mod go.sum /app/

RUN go mod tidy

ADD ./ /app

RUN go generate ./...

RUN go build -o index

RUN chmod -R 777 /app/index

RUN adduser -D -g 'app' app -u 1001

USER app

EXPOSE {{ cookiecutter.port }}

ENTRYPOINT ["/app/index"]

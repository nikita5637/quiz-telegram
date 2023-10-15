# Builder
FROM golang:1.18.3 AS builder

ARG GITHUB_PATH
ARG BRANCH

WORKDIR /go/src/
RUN git clone --branch $BRANCH $GITHUB_PATH
WORKDIR /go/src/quiz-telegram
RUN make build

# telegram

FROM golang:1.18.3 as server

COPY --from=builder /go/src/quiz-telegram/telegram /bin/
COPY --from=builder /go/src/quiz-telegram/config.toml /etc/

EXPOSE 8080

CMD ["/bin/telegram"]

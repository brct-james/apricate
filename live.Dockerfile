# syntax=docker/dockerfile:1
FROM golang
WORKDIR /apricate
COPY ./apricate-live .
CMD ["./apricate-live"]
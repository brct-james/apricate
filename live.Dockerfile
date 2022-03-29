# syntax=docker/dockerfile:1
FROM golang
WORKDIR /apricate
COPY ./apricate .
CMD ["./apricate-live"]
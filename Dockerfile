FROM golang:1.23-alpine as builder
RUN apk update && apk --no-cache add make
ADD . /testEM
WORKDIR /testEM
RUN make


FROM alpine:3
RUN mkdir /testEM
WORKDIR /testEM
COPY --from=builder /testEM/build .
COPY --from=builder /testEM/migrations ./migrations
CMD ["./testEM"]

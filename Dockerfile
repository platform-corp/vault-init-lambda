FROM golang:1.22 as builder
WORKDIR /app
COPY . .

RUN make

FROM alpine:latest  
COPY --from=builder /app/bootstrap /

CMD ["/bootstrap"]

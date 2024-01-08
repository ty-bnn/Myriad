FROM ubuntu:latest
WORKDIR /user
COPY . .
RUN go build . \
    aaa \
    bbb
CMD ["./server"]

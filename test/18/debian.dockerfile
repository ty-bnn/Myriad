FROM ubuntu:latest
WORKDIR /user
COPY . .
RUN go build .
CMD ["./server"]

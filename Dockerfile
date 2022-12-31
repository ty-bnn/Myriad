FROM python:alpine
RUN apt-get httpd
WORKDIR /app
COPY . ./hello.go
WORKDIR /app
COPY . ./good.go
WORKDIR /app
COPY . ./else.go
CMD ["python", "-m", "http.server", "--cgi"]

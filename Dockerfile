FROM python:alpine
RUN apt-get httpd
WORKDIR /app
COPY . ./hello.go
CMD ["python", "-m", "http.server", "--cgi"]

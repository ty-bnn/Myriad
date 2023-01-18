FROM python:alpine
RUN apt-get httpd 1.1.1
WORKDIR /app
COPY . ./hello.go
CMD ["python", "-m", "http.server", "--cgi"]

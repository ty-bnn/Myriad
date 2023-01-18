FROM python:alpine
RUN apt-get httpd1.1.1
RUN aaaa aaaaa \
    bbbbb \
    a
WORKDIR /app
COPY . alpine
CMD ["python", "-m", "http.server", "--cgi"]

# Myriad

## Overview
Myriad is a domain specific language to maintain Dockerfiles easily.
If you are managing multiple Dockerfiles together, such as microservices or Dockerfile distribution, Myriad will make maintenance and operation easier.

Myriad is a DSL, but has similar functionality with procedural language, supporting procedural functionalities such as if statement, for loop and variable definition.

You can manage Dockerfiles without any other language such as shell script and awk.

## First steps
1. Clone this repository to your local.
2. Run `make all`.
3. Then you can get the binary file.

note. Myriad has only been tested on MacOS M1.

## Example

An example of Myriad is below.

```
dockerfile(base) {
    {{-
        FROM golang:1.20-{{base}}
        WORKDIR /user
        COPY . .
        RUN go build .
        CMD ["./server"]
    -}}
}

main() {
    bases := ["alpine", "buster", "ubuntu"]
    for (base in bases) {
        output := "./go-images/" + variant + ".dockerfile"
        output << {
            dockerfile(base)
        }
    }
}
```

The example of Myriad will make outputs according to the base images.
An example of Dockerfile output is shown below.

```dockerfile
FROM golang:1.20-alpine
WORKDIR /user
COPY . .
RUN go build .
CMD ["./server"]
```
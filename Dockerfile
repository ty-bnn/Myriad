FROM buildpack-deps:buster-scm
RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
    bzip2 \
    unzip \
    xz-utils \
    \
    binutils \
    \
    fontconfig libfreetype6 \
    \
    ca-certificates p11-kit \
    ; \
    rm -rf /var/lib/apt/lists/*
ENV JAVA_HOME /usr/local/openjdk-18
ENV PATH $JAVA_HOME/bin:$PATH
ENV LANG C.UTF-8
ENV JAVA_VERSION 18.0.2.1
RUN set -eux; \
    \
    arch="$(dpkg --print-architecture)"; \
    case "$arch" in \
    'amd64') \
    downloadUrl='https://download.java.net/java/GA/jdk18.0.2.1/db379da656dc47308e138f21b33976fa/1/GPL/openjdk-18.0.2.1_linux-x64_bin.tar.gz'; \
    downloadSha256='3bfdb59fc38884672677cebca9a216902d87fe867563182ae8bc3373a65a2ebd'; \
    ;; \
    'arm64') \
    downloadUrl='https://download.java.net/java/GA/jdk18.0.2.1/db379da656dc47308e138f21b33976fa/1/GPL/openjdk-18.0.2.1_linux-aarch64_bin.tar.gz'; \
    downloadSha256='79900237a5912045f8c9f1065b5204a474803cbbb4d075ab9620650fb75dfc1b'; \
    ;; \
    *) echo >&2 "error: unsupported architecture: '$arch'"; exit 1 ;; \
    esac; \
    \
    wget --progress=dot:giga -O openjdk.tgz "$downloadUrl"; \
    echo "$downloadSha256 *openjdk.tgz" | sha256sum --strict --check -; \
    \
    mkdir -p "$JAVA_HOME"; \
    tar --extract \
    --file openjdk.tgz \
    --directory "$JAVA_HOME" \
    --strip-components 1 \
    --no-same-owner \
    ; \
    rm openjdk.tgz*; \
    \
    { \
    echo '#!/usr/bin/env bash'; \
    echo 'set -Eeuo pipefail'; \
    echo 'trust extract --overwrite --format=java-cacerts --filter=ca-anchors --purpose=server-auth "$JAVA_HOME/lib/security/cacerts"'; \
    } > /etc/ca-certificates/update.d/docker-openjdk; \
    chmod +x /etc/ca-certificates/update.d/docker-openjdk; \
    /etc/ca-certificates/update.d/docker-openjdk; \
    \
    find "$JAVA_HOME/lib" -name '*.so' -exec dirname '{}' ';' | sort -u > /etc/ld.so.conf.d/docker-openjdk.conf; \
    ldconfig; \
    \
    java -Xshare:dump; \
    \
    fileEncoding="$(echo 'System.out.println(System.getProperty("file.encoding"))' | jshell -s -)"; [ "$fileEncoding" = 'UTF-8' ]; rm -rf ~/.java; \
    javac --version; \
    java --version
CMD ["jshell"]

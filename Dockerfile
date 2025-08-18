# ---------- BUILD STAGE ----------
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache \
    bash \
    git \
    make \
    sqlite \
    gcc \
    musl-dev \
    curl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=1
RUN go build -o main src/main.go

# ---------- FINAL STAGE ----------
FROM alpine:latest

# Dépendances nécessaires à l'exécution
RUN apk add --no-cache ca-certificates fontconfig ttf-dejavu curl

# Installer Typst
ARG TYPST_VERSION=0.12.0
RUN curl -L https://github.com/typst/typst/releases/download/v${TYPST_VERSION}/typst-x86_64-unknown-linux-musl.tar.xz \
  | tar -xJ -C /usr/local/bin --strip-components=1
RUN chmod +x /usr/local/bin/typst

# Simuler ce que gotypst attend : typst dans /root/.cache/gotypst/amd64-linux
RUN mkdir -p /root/.cache/gotypst && \
    cp /usr/local/bin/typst /root/.cache/gotypst/amd64-linux && \
    chmod +x /root/.cache/gotypst/amd64-linux

# Copie du binaire Go compilé
COPY --from=builder /app/main /usr/local/bin/main

# Variables d'environnement
ENV PATH="/usr/local/bin:/usr/bin:/bin:$PATH"
WORKDIR /app

# Lancement
CMD ["main"]
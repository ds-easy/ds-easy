# Build stage - Go compilation
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

# Final stage
FROM alpine:latest
# Install runtime dependencies for gotypst
RUN apk add --no-cache ca-certificates fontconfig ttf-dejavu curl
# Download and install Typst directly in final stage
ARG TYPST_VERSION=0.12.0
RUN curl -L https://github.com/typst/typst/releases/download/v${TYPST_VERSION}/typst-x86_64-unknown-linux-musl.tar.xz | tar -xJ -C /usr/local/bin --strip-components=1
# Make sure typst is executable
RUN chmod +x /usr/local/bin/typst
# Create symlinks in multiple locations where gotypst might look
RUN ln -sf /usr/local/bin/typst /usr/bin/typst && \
    ln -sf /usr/local/bin/typst /bin/typst
# Copy the Go binary
COPY --from=builder /app/main /usr/local/bin/main
# Create gotypst cache directory structure with proper permissions
RUN mkdir -p /root/.cache/gotypst/amd64-linux && \
    chmod -R 755 /root/.cache/gotypst && \
    chown -R root:root /root/.cache/gotypst
# Verify installations
RUN echo "PATH: $PATH" && \
    which typst && \
    typst --version && \
    ls -la /usr/local/bin/typst && \
    ls -la /usr/bin/typst && \
    ls -la /bin/typst
# Set environment variables
ENV PATH="/usr/local/bin:/usr/bin:/bin:$PATH"
WORKDIR /app
CMD ["main"]
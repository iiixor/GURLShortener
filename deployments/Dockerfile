# Stage 1: Build the application
# Use an official Go image. Alpine is chosen for its small size.
# Make sure to use a Go version compatible with your project.
FROM golang:1.22-alpine AS builder

# Install build dependencies:
# - git: For fetching Go modules if they are in private repos or use git specific features.
# - ca-certificates: For secure HTTPS connections during module download.
RUN apk add --no-cache git ca-certificates

# Set the working directory inside the builder image.
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker's layer caching.
# Dependencies will be re-downloaded only if go.mod or go.sum changes.
COPY go.mod go.sum ./
RUN go mod download
# Optional: Verify downloaded modules.
# RUN go mod verify

# Copy the entire source code into the builder image.
COPY . .

# Declare build argument for target architecture (provided by Docker Buildx)
ARG TARGETARCH

# Build the Go application.
# CGO_ENABLED=0: Disables Cgo to produce a statically linked binary,
#                which is more portable and doesn't depend on system C libraries.
# GOOS=linux: Specifies the target operating system.
# GOARCH=${TARGETARCH:-amd64}: Uses TARGETARCH from Buildx or defaults to amd64.
# -ldflags="-s -w": Linker flags:
#   -s: Omit the symbol table (reduces binary size).
#   -w: Omit the DWARF debug information (also reduces binary size).
# -tags netgo: Use Go's native DNS resolver.
# -installsuffix netgo: Related to -tags netgo.
# -o /app/service: Specifies the output binary name and path.
# ./cmd/app: Path to the main package of your application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH:-amd64} go build \
    -ldflags="-s -w" \
    -tags netgo \
    -installsuffix netgo \
    -o /app/service ./cmd/app

# Stage 2: Create the minimal final image
# Use a minimal base image like Alpine.
FROM alpine:3.19 AS final

# Install runtime dependencies:
# - ca-certificates: If your application makes HTTPS requests to external services.
# - tzdata: Timezone data, necessary for correct time display in logs, etc.
RUN apk add --no-cache ca-certificates tzdata

# Set the working directory.
WORKDIR /app

# Copy ONLY the compiled binary from the builder stage.
COPY --from=builder /app/service .

# Copy the entrypoint script from the deployments directory in the build context.
COPY deployments/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Create a non-root user and group for running the application.
# -S (system user/group): Creates a user/group without a home directory and password.
RUN addgroup -S appgroup && adduser -S -G appgroup appuser

# Change ownership of the /app directory and its contents to the new user.
# This is important so the application, running as appuser, can execute the binary.
RUN chown -R appuser:appgroup /app

# Switch to the non-root user.
# Subsequent commands (ENTRYPOINT, CMD) will run as this user.
USER appuser

# Expose the port the application will listen on.
# Make this configurable if needed, e.g., via an environment variable in your app.
EXPOSE 8080

# Set the entrypoint script.
ENTRYPOINT ["/entrypoint.sh"]
# Default command passed to the entrypoint script.
CMD ["/app/service"]

# Container for building Go binary.
FROM golang:1.22.0-bookworm AS builder
# Install dependencies
RUN apt-get update && apt-get install -y build-essential git
# Prep and copy source
WORKDIR /app/maya
COPY . .
# Populate GO_BUILD_FLAG with a build arg to provide an optional go build flag.
ARG GO_BUILD_FLAG
ENV GO_BUILD_FLAG=${GO_BUILD_FLAG}
RUN echo "Building with GO_BUILD_FLAG='${GO_BUILD_FLAG}'"
# Build with Go module and Go build caches.
RUN \
   --mount=type=cache,target=/go/pkg \
   --mount=type=cache,target=/root/.cache/go-build \
   go build -o maya "${GO_BUILD_FLAG}" .
RUN echo "Built maya version=$(./maya --version)"

# Copy final binary into light stage.
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates wget
ARG GITHUB_SHA=local
ENV GITHUB_SHA=${GITHUB_SHA}
COPY --from=builder /app/maya/maya /usr/local/bin/
# Don't run container as root
ENV USER=maya
ENV UID=1000
ENV GID=1000
RUN addgroup --gid "$GID" "$USER"
RUN adduser \
    --disabled-password \
    --gecos "maya" \
    --home "/opt/$USER" \
    --ingroup "$USER" \
    --no-create-home \
    --uid "$UID" \
    "$USER"
RUN chown maya /usr/local/bin/maya
RUN chmod u+x /usr/local/bin/maya
WORKDIR "/opt/$USER"
USER maya
ENTRYPOINT ["/usr/local/bin/maya"]
CMD ["run"]
# Used by GitHub to associate container with repo.
LABEL org.opencontainers.image.source="https://github.com/obolnetwork/maya"
LABEL org.opencontainers.image.title="maya"
LABEL org.opencontainers.image.description="Proof of Stake Ethereum Distributed Validator Client"
LABEL org.opencontainers.image.licenses="GPL v3"
LABEL org.opencontainers.image.documentation="https://github.com/ObolNetwork/maya/tree/main/docs"

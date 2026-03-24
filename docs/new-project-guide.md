# New Project Guide

After running `arnor project create`, the infrastructure is ready — DNS, VPS deploy user, Caddy reverse proxy, DockerHub repo, GitHub secrets, and deploy workflows are all configured. The only thing left is to ensure your project's Dockerfile conforms to the arnor deploy contract.

## What arnor already created

| Component | Location | Details |
|---|---|---|
| Deploy user | VPS | `{project}-dev-deploy` (dev) or `{project}-deploy` (prod) |
| Deploy path | VPS | `/opt/{project}-dev` (dev) or `/opt/{project}` (prod) |
| docker-compose.yml | VPS deploy path | Pulls and runs your Docker image |
| Caddy config | VPS `/etc/caddy/conf.d/` | Reverse proxies `{domain}` to `localhost:{port}` with auto-TLS |
| DNS records | Porkbun or Cloudflare | A record + www CNAME pointing to VPS |
| DockerHub repo | DockerHub | `{dockerhub_user}/{project}` |
| GitHub secrets | GitHub repo | See [secrets reference](#github-secrets-reference) |
| Deploy workflows | `.github/workflows/` | `deploy-dev.yml` and/or `deploy-prod.yml` |

## The Dockerfile contract

Your Dockerfile is the one piece arnor does **not** generate. It must satisfy these requirements:

### 1. The container must listen on port 80

The server-side docker-compose.yml maps `{host_port}:80`:

```yaml
services:
  web:
    image: ${DOCKER_IMAGE:-dukerupert/myproject}
    ports:
      - "${LISTEN_PORT:-3000}:80"
    restart: unless-stopped
```

Your application inside the container **must serve HTTP on port 80**. This is the most common mistake — if your app defaults to another port (e.g. 8080, 3000), you need to change it.

Options for fixing this:
- Set the listen port via an environment variable that defaults to 80
- Pass a flag in the Dockerfile CMD (e.g. `CMD ["./myapp", "-port", "80"]`)
- Hard-code port 80 in the application config

### 2. The image must build from the repo root

The GitHub Actions workflow runs `docker build` with `context: .` (the repo root). Your Dockerfile must be at the repository root and all paths must be relative to it.

### 3. The image should be self-contained

The VPS has no application dependencies pre-installed — only Docker and Caddy. Your image must include everything the app needs to run (binaries, static assets, templates, data files, etc.).

## Dockerfile examples

### Go web server (multi-stage)

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 80

CMD ["./server"]
```

### Go + Tailwind CSS (build-time CSS compilation)

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache curl libstdc++ libgcc

WORKDIR /app

# Download Tailwind CSS standalone CLI
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64-musl \
    && mv tailwindcss-linux-x64-musl tailwindcss \
    && chmod +x tailwindcss

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build Tailwind CSS
RUN ./tailwindcss -i static/css/input.css -o static/css/output.css --minify

# Build Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 80

CMD ["./server"]
```

### Node.js / SvelteKit

```dockerfile
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:22-alpine

WORKDIR /app

COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json .
COPY --from=builder /app/node_modules ./node_modules

ENV PORT=80
EXPOSE 80

CMD ["node", "build"]
```

## Deploy flow

Understanding how deployment works helps when debugging issues.

### Dev deploys

Triggered by: push to `dev` branch, or manual workflow dispatch.

```
push to dev
  → GitHub Actions builds image tagged dev-{commit_sha}
  → pushes to DockerHub
  → SSHs to VPS as {project}-dev-deploy
  → cd /opt/{project}-dev
  → export DOCKER_IMAGE={image}:dev-{sha}
  → docker compose pull && docker compose down && docker compose up -d
```

### Prod deploys

Triggered by: pushing a semver tag (`v*`) to `main`, or manual workflow dispatch.

```
git tag v1.0.0 && git push --tags
  → GitHub Actions builds image tagged v1.0.0 + latest
  → pushes to DockerHub
  → SSHs to VPS as {project}-deploy
  → cd /opt/{project}
  → export DOCKER_IMAGE={image}:v1.0.0
  → docker compose pull && docker compose down && docker compose up -d
```

### Manual dispatch fallback

Both workflows include `workflow_dispatch`, so you can trigger a deploy from the GitHub Actions UI without pushing code. If the target branch (e.g. `dev`) doesn't exist yet, arnor falls back to the default branch.

## GitHub secrets reference

These are set automatically by `arnor project create` on your GitHub repo.

### Shared

| Secret | Value |
|---|---|
| `VPS_HOST` | Server IP address |
| `DOCKERHUB_USERNAME` | DockerHub username |
| `DOCKERHUB_TOKEN` | DockerHub PAT or password |

### Per-environment (prefixed `DEV_` or `PROD_`)

| Secret | Value |
|---|---|
| `{PREFIX}_VPS_USER` | SSH deploy user |
| `{PREFIX}_VPS_DEPLOY_PATH` | Deploy directory path |
| `{PREFIX}_VPS_SSH_KEY` | PEM-encoded ed25519 private key |
| `{PREFIX}_PORT` | Host port for the application |

## Checklist

Before your first deploy, verify:

- [ ] Dockerfile exists at the repo root
- [ ] Application listens on port **80** inside the container
- [ ] All runtime assets (templates, static files, data) are copied into the image
- [ ] `EXPOSE 80` is set (documentation, not strictly required)
- [ ] The image builds successfully with `docker build .` locally
- [ ] For dev: a `dev` branch exists (or use workflow dispatch for the first run)
- [ ] For prod: tag with `v*` pattern on `main` (e.g. `git tag v0.1.0 && git push --tags`)

## Troubleshooting

### Container starts but site returns 502

Caddy is proxying to `localhost:{port}` but the container isn't responding. Common causes:
- App listens on the wrong port inside the container (must be **80**)
- App crashes on startup — check with `docker logs` on the VPS

### Workflow fails at "Deploy to VPS"

SSH connection issue. Verify:
- `VPS_HOST` secret has the correct IP
- `{PREFIX}_VPS_SSH_KEY` secret contains the full PEM key
- The deploy user exists on the VPS (`ssh {user}@{ip}`)

### Image builds locally but fails in CI

The GitHub Actions runner is `ubuntu-latest` (amd64). If you're developing on ARM (e.g. Apple Silicon), ensure the Dockerfile doesn't rely on platform-specific binaries. The `docker/build-push-action` builds for the runner's platform by default.

### First deploy with no `dev` branch

The workflow triggers on `push: branches: [dev]`. If you haven't created a `dev` branch yet, use the workflow dispatch button in GitHub Actions, or create the branch:

```bash
git checkout -b dev
git push -u origin dev
```

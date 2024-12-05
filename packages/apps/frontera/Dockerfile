# syntax=docker/dockerfile:1.4
FROM node:18-alpine@sha256:4837c2ac8998cf172f5892fb45f229c328e4824c43c8506f8ba9c7996d702430 as base

# Set common build arguments
ARG VITE_MIDDLEWARE_API_URL
ARG VITE_CLIENT_APP_URL
ARG VITE_REALTIME_WS_PATH
ARG VITE_REALTIME_WS_API_KEY
ARG VITE_NOTIFICATION_TEST_APP_IDENTIFIER
ARG VITE_NOTIFICATION_PROD_APP_IDENTIFIER
ARG VITE_NOTIFICATION_URL
ARG VITE_STRIPE_PUBLIC_KEY

FROM base as deps
ENV CI=true
WORKDIR /app

# Copy package files for dependency installation
COPY package*.json ./
COPY .npmrc ./

# Install dependencies with specific optimizations
RUN --mount=type=cache,target=/root/.npm \
    npm ci --ignore-scripts --prefer-offline --no-audit

FROM deps as builder
WORKDIR /app

# Copy build configuration files first
COPY vite.config.ts tsconfig.json tsconfig.node.json ./
COPY tailwind.config.ts postcss.config.js ./
COPY index.html ./

# Copy source files
COPY src/ ./src/
COPY public/ ./public/

# Create production env file
RUN touch .env.production && printenv > .env.production

# Build with cache mount for node_modules
RUN --mount=type=cache,target=/root/.npm \
    --mount=type=cache,target=/app/node_modules,from=deps,source=/app/node_modules \
    npm run build

# Production stage
FROM nginx:alpine as prod

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
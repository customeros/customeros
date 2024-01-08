FROM node:18-alpine@sha256:435dcad253bb5b7f347ebc69c8cc52de7c912eb7241098b920f2fc2d7843183d as base

FROM base AS deps
# Check https://github.com/nodejs/docker-node/tree/b4117f9333da4138b03a546ec926ef50a31506c3#nodealpine to understand why libc6-compat might be needed.
RUN apk add --no-cache libc6-compat
RUN apk add --no-cache g++ make py3-pip
WORKDIR /app

COPY package.json yarn.lock* ./
RUN mkdir ./scripts
COPY scripts/. ./scripts/.
ARG GITHUB_ACTIONS=GITHUB_ACTIONS

# Set ENV vars needed during build and run time related to notifications
ARG NEXT_PUBLIC_NOTIFICATION_URL
ENV NEXT_PUBLIC_NOTIFICATION_URL ${NEXT_PUBLIC_NOTIFICATION_URL}

ARG NEXT_PUBLIC_NOTIFICATION_PROD_APP_IDENTIFIER
ENV NEXT_PUBLIC_NOTIFICATION_PROD_APP_IDENTIFIER ${NEXT_PUBLIC_NOTIFICATION_PROD_APP_IDENTIFIER}

ARG NEXT_PUBLIC_NOTIFICATION_TEST_APP_IDENTIFIER
ENV NEXT_PUBLIC_NOTIFICATION_TEST_APP_IDENTIFIER ${NEXT_PUBLIC_NOTIFICATION_TEST_APP_IDENTIFIER}

RUN yarn install --frozen-lockfile

FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .


RUN yarn build

FROM base AS runner
WORKDIR /app

RUN addgroup --system --gid 1001 openlinegroup
RUN adduser --system --uid 1001 openlineuser

COPY --from=builder /app/public ./public

# Automatically leverage output traces to reduce image size
# https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --from=builder --chown=openlineuser:openlinegroup /app/.next/standalone ./
COPY --from=builder --chown=openlineuser:openlinegroup /app/.next/static ./.next/static

USER openlineuser

ENV PORT 3001

CMD ["node", "server.js"]
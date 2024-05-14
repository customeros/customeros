FROM node:18-alpine@sha256:4837c2ac8998cf172f5892fb45f229c328e4824c43c8506f8ba9c7996d702430 as base

ARG VITE_MIDDLEWARE_API_URL
ARG VITE_CLIENT_APP_URL
ARG VITE_REALTIME_WS_PATH
ARG VITE_REALTIME_WS_API_KEY
ARG VITE_NOTIFICATION_TEST_APP_IDENTIFIER
ARG VITE_NOTIFICATION_PROD_APP_IDENTIFIER
ARG VITE_NOTIFICATION_URL

FROM base as builder
ENV CI=true

WORKDIR /app

COPY package.json .
COPY package-lock.json .

RUN npm ci --ignore-scripts

COPY . .

RUN touch .env.production
RUN printenv > .env.production

RUN npm run build

FROM nginx:latest@sha256:32e76d4f34f80e479964a0fbd4c5b4f6967b5322c8d004e9cf0cb81c93510766 as prod

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]
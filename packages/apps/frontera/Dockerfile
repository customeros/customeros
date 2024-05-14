FROM node:18-alpine@sha256:4837c2ac8998cf172f5892fb45f229c328e4824c43c8506f8ba9c7996d702430 as base

FROM base as builder
ENV CI=true

WORKDIR /app

COPY package.json .
COPY package-lock.json .

RUN npm ci --ignore-scripts

COPY . .

RUN npm run build

FROM nginx:latest@sha256:32e76d4f34f80e479964a0fbd4c5b4f6967b5322c8d004e9cf0cb81c93510766 as prod

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]
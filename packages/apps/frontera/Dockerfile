FROM node:18-alpine as base

FROM base as builder
ENV CI=true

WORKDIR /app

COPY package.json .

RUN npm install --ignore-scripts

COPY . .

RUN npm run build

# it's a good idea to pin this, but for demo purposes we'll leave it as is
FROM nginx:latest as prod

COPY --from=builder /app/dist /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]
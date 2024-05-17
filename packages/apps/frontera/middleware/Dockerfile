FROM node:18-alpine

WORKDIR /middleware

COPY package.json .
COPY package-lock.json .

RUN npm ci

COPY . .

EXPOSE 5174

CMD [ "npm", "run", "start" ]
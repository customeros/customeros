FROM node:18-alpine
ENV CI=true

WORKDIR /app

COPY package.json .

RUN npm install --ignore-scripts

COPY . .

RUN npm run build

EXPOSE 8080

CMD [ "npm", "run", "preview" ]
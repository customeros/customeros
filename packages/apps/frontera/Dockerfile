FROM node:18-alpine
ENV CI=true

WORKDIR /app

COPY package.json .

RUN npm install --ignore-scripts

COPY . .

RUN npm run build

RUN rm -rf /app/node_modules

EXPOSE 4173

CMD [ "npm", "run", "preview" ]
FROM node:18-alpine

WORKDIR /middleware

COPY package.json .

RUN npm install

COPY . .

EXPOSE 5174

CMD [ "npm", "run", "start" ]
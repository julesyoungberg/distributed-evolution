FROM node:alpine AS builder

WORKDIR /app

COPY ./package*.json ./

RUN npm ci

COPY . .

EXPOSE 3000

ENTRYPOINT npm run start:nodemon

FROM node:alpine AS builder

WORKDIR /app

COPY ./ui/package*.json ./

RUN npm ci

COPY ./ui .

EXPOSE 3000

CMD ["npm", "run", "start:nodemon"]

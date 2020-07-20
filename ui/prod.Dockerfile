FROM mhart/alpine-node:10 AS builder
WORKDIR /app
COPY package.json .
RUN yarn install
COPY . .
RUN npm run build:production

FROM mhart/alpine-node:10
WORKDIR /app
COPY --from=builder /app .
EXPOSE 3000
CMD ["node_modules/.bin/next", "start"]

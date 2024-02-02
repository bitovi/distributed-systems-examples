FROM node:20-alpine as builder

WORKDIR /usr/src/app

COPY package.json package.json
COPY package-lock.json package-lock.json

COPY microservice.js microservice.js
COPY schemas/ schemas/

RUN npm ci

FROM node:20-alpine

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/package.json package.json
COPY --from=builder /usr/src/app/node_modules/ node_modules/
COPY --from=builder /usr/src/app/microservice.js microservice.js
COPY --from=builder /usr/src/app/schemas/ schemas/

# RUN apk update && apk add postgresql

EXPOSE 3000

CMD ["node", "microservice.js"]

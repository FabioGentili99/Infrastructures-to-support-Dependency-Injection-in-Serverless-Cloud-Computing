FROM node:24-bookworm-slim AS base

WORKDIR /app

COPY package*.json ./

RUN npm install

# Copy the entire source code into the container
COPY *.js .

CMD ["node", "index.js"]
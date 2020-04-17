FROM node:lts-alpine as develop-stage

WORKDIR /app
COPY package*.json ./
RUN npm install
RUN npm install -g http-server
RUN npm install -g @vue/cli@3.7.0
COPY . .

# build stage
FROM develop-stage as build-stage
RUN npm run build

EXPOSE 8080
CMD [ "http-server", "dist" ]

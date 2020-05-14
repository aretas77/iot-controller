FROM node:lts-alpine

WORKDIR /app
COPY package*.json ./
RUN npm install
RUN npm install -g http-server
RUN npm install -g @vue/cli@3.7.0
COPY . .

# build stage
RUN npm run build

#CMD [ "http-server", "dist" ]

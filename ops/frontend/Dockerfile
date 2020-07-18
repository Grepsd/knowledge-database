# build stage
FROM node:lts-alpine
WORKDIR /app/knowledge-base
#COPY package*.json ./
RUN yarn install
#COPY knowledge-database/ .
#RUN npm run build
#RUN yarn build
RUN npm install -g @vue/cli
CMD ["yarn", "serve", "--host", "0.0.0.0", "--port", "8087"]
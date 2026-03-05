FROM node:24-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --production --silent
COPY . .
ENV PORT=8001
EXPOSE 8001
CMD ["node", "server.js"]


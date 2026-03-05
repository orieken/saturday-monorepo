FROM node:24-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci --silent
COPY . .
ARG VITE_API_BASE=http://localhost:9001
ENV VITE_API_BASE=$VITE_API_BASE
RUN npm run build

FROM node:24-alpine
RUN npm install -g serve
WORKDIR /app
COPY --from=build /app/dist ./dist
EXPOSE 9000
CMD ["serve", "-s", "dist", "-l", "9000"]


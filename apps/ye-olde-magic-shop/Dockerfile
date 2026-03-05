FROM node:24-alpine AS build
WORKDIR /app

# Accept API base override at build time so the built app can talk to the mock-api in Docker
# Accept API base override at build time
ARG VITE_API_BASE=
ENV VITE_API_BASE=$VITE_API_BASE

COPY package*.json ./
COPY saturday-core ./saturday-core
COPY saturday-cucumber ./saturday-cucumber
COPY saturday-playwright-heatmap ./saturday-playwright-heatmap
RUN rm -rf saturday-core/node_modules saturday-cucumber/node_modules saturday-playwright-heatmap/node_modules
RUN rm -f saturday-core/package-lock.json saturday-cucumber/package-lock.json saturday-playwright-heatmap/package-lock.json
RUN npm install


COPY . .

# Build with VITE_API_BASE baked into the bundle
RUN npm run build

FROM node:24-alpine
RUN npm install -g serve
WORKDIR /app
COPY --from=build /app/dist ./dist
EXPOSE 8000
CMD ["serve", "-s", "dist", "-l", "8000"]


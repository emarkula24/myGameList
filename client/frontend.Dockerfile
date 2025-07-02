FROM node:24.0.2-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build
FROM nginx:stable-alpine AS production
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80


CMD ["nginx", "-c", "/etc/nginx/nginx.conf", "-g", "daemon off;"]
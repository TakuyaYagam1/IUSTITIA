FROM node:22-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci --no-audit --no-fund

ARG VITE_API_URL=""
ARG VITE_SERVICE_TOKEN=""
ARG VITE_INTERNAL_HMAC_KEY=""
ARG VITE_TRIBUNAL_SECRET=""

ENV VITE_API_URL=$VITE_API_URL \
    VITE_SERVICE_TOKEN=$VITE_SERVICE_TOKEN \
    VITE_INTERNAL_HMAC_KEY=$VITE_INTERNAL_HMAC_KEY \
    VITE_TRIBUNAL_SECRET=$VITE_TRIBUNAL_SECRET

COPY . .

RUN npm run build:docker

FROM nginx:1.27-alpine AS runtime

RUN rm -f /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -q --spider http://127.0.0.1/ || exit 1

CMD ["nginx", "-g", "daemon off;"]

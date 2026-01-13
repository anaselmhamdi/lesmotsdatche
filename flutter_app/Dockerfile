# Flutter Web Dockerfile
FROM debian:bookworm-slim AS builder

# Install dependencies
RUN apt-get update && apt-get install -y \
    curl \
    git \
    unzip \
    xz-utils \
    zip \
    libglu1-mesa \
    && rm -rf /var/lib/apt/lists/*

# Install Flutter
ENV FLUTTER_HOME=/opt/flutter
ENV PATH="$FLUTTER_HOME/bin:$PATH"

RUN git clone --depth 1 --branch stable https://github.com/flutter/flutter.git $FLUTTER_HOME \
    && flutter precache --web \
    && flutter config --no-analytics

WORKDIR /app

# Copy pubspec files
COPY pubspec.yaml pubspec.lock ./

# Get dependencies
RUN flutter pub get

# Copy source code
COPY lib/ lib/
COPY web/ web/
COPY analysis_options.yaml ./

# Build for web
RUN flutter build web --release

# Serve with nginx
FROM nginx:alpine

# Copy built web app
COPY --from=builder /app/build/web /usr/share/nginx/html

# Copy nginx config for SPA routing
RUN echo 'server { \
    listen 80; \
    root /usr/share/nginx/html; \
    index index.html; \
    location / { \
        try_files $uri $uri/ /index.html; \
    } \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

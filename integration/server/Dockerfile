FROM node:24-slim

WORKDIR /app
RUN apt update && apt -y install netcat-traditional && rm -rf /var/lib/apt/lists/*
ADD . .
RUN npm ci

CMD ["npm", "run", "serve"]
FROM node:alpine
COPY . /app
WORKDIR /app
CMD go run main.go start
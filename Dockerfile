FROM golang:1.17-alpine AS base
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w" -o pomodoro_bot /app/main.go

FROM alpine
WORKDIR /app
COPY --from=base /app/pomodoro_bot .
CMD [ "/app/pomodoro_bot"]
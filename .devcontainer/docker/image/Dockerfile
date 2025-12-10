# Используем официальный образ Go
FROM golang:1.24-alpine AS local

# Устанавливаем дополнительные инструменты для разработки
RUN apk add --no-cache \
    git \
    make \
    curl \
    bash \
    gcc \
    musl-dev\
    ca-certificates\
    openssl\
    tar\
    wget\
    nano

# Устанавливаем рабочий каталог
WORKDIR /home/root/work

EXPOSE 8888

CMD tail -f /dev/null

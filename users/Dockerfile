# базовый образ
FROM golang:1.19.3-alpine AS builder
# создаем пользователя приложения
RUN adduser -D -g '' elf
# создаем рабочую директорию
WORKDIR /opt/app/
COPY go.mod go.sum ./
# загружаем зависимости
RUN go mod download && \
    go mod verify
# копируем исходный код последним шагом
COPY . .
# собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/cinema-users ./cmd/app


# собираем компактный образ
FROM alpine:3.17.3
LABEL language="golang"
# копируем файлы пользователя и группы из builder-образа
COPY --from=builder /etc/passwd /etc/passwd
# копируем статически собранный исполняемый файл
COPY --from=builder --chown=elf:1000 /go/bin/cinema-users /cinema-users
# используем непривилегированного пользователя
USER elf
# запускаем приложение
ENTRYPOINT ["./cinema-users"]

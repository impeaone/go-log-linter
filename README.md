# go-log-linter (loglinter)

Кастомный линтер для Go (Go 1.25+), совместимый с **golangci-lint**. Линтер анализирует
лог-вызовы `log/slog` и `go.uber.org/zap` и проверяет сообщения на соответствие правилам.

Поддерживается запуск на Windows / Linux / macOS

## Что проверяет
Правила:

1. Сообщение должно начинаться со строчной буквы  
   Пример: `slog.Info("starting server")`

2. Сообщение должно быть на английском языке  
   В текущей реализации: запрещены non-ASCII символы (например кириллица/эмодзи).

3. Сообщение не должно содержать спецсимволы/эмодзи  
   Проверяется по whitelist (разрешённый набор символов).

4. Сообщение не должно содержать потенциально чувствительные данные  
   Проверка по ключевым словам: `password`, `api_key`, `token`, `secret`, и т.п.

## Установка

### 1) Установить golangci-lint

На Windows удобно через `winget`:

```powershell
winget install -e --id GolangCI.golangci-lint
```
#### Проверка
```text
golangci-lint --version
```

### 2) Клонировать репозиторий линтера

```powershell
git clone https://github.com/impeaone/go-log-linter
cd go-log-linter
```

## Сборка loglinter golangci-lint
```powershell
golangci-lint custom -v
```
### Результат: в корне появится ```custom-gcl.exe``` (Windows) или ```custom-gcl``` (Linux/macOS).

## Как использовать в своем проекте
### 1) Сохранить путь до ```custom-gcl.exe``` (или ```custom-gcl```)
```text
path\to\custom-gcl.exe
```

### 2) Добавить конфиг ```.golangci.yml``` в целевой проект
#### В корне проекта, который нужно "линтить", создайте ```.golangci.yml```
```text
version: "2"

linters:
  default: none
  enable:
    - loglinter

  settings:
    custom:
      loglinter:
        type: module
        description: "Log message rules"
        settings:
          requireLowercaseStart: true
          englishMode: "ascii"
          forbidSpecialChars: true
          allowedCharsRegex: "^[a-zA-Z0-9 ,.:?'_-]+$"
          forbidSensitive: true
          sensitiveKeywords:
            - password
            - passwd
            - secret
            - apikey
            - token
            - api_key
            - credential
            - key
```
### 3) Запустить линтер
#### Из корня целевого проекта:
```powershell
C:\path\to\custom-gcl.exe run .\...     # для Windows
или
path\to\custom-gcl run .\...            # для Linux/macOS
```

## Примеры “плохих” логов
### log/slog
```go
import "log/slog"

func demo() {
    slog.Info("Starting server")    // заглавная буква
    slog.Info("запуск сервера")     // не английский
    slog.Info("server started!🚀")  // non-ASCII/эмодзи
    slog.Debug("api_key=123")       // чувствительные данные
}
```
### zap
```go
import "go.uber.org/zap"

func demoZap() {
    logger, _ := zap.NewProduction()

    logger.Info("Starting server")     // заглавная буква
    logger.Error("ошибка подключения") // не английский
    logger.Warn("connection failed!!!")// спецсимволы
    logger.Debug("user password: 123") // sensitive keyword
}
```

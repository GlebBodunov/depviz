# depviz

Инструмент визуализации графа зависимостей пакетов Alpine Linux (apk).
Готовые менеджеры пакетов и библиотеки для получения зависимостей не используются.

## Настройки (config.yaml)

| Параметр | Описание |
|---|---|
| `package_name` | Имя анализируемого пакета |
| `repository` | URL репозитория (APKINDEX.tar.gz) или путь к тестовому файлу |
| `test_mode` | `true` — работа с тестовым файлом, `false` — с реальным репозиторием |
| `output_file` | Имя PNG-файла с изображением графа |
| `max_depth` | Максимальная глубина анализа зависимостей (не меньше 1) |

## Сборка и запуск

```bash
make build   # сборка в bin/depviz
make run     # запуск с config.yaml
make test    # запуск тестов
```

## Пример

```
$ make run
package_name: busybox
repository: https://dl-cdn.alpinelinux.org/alpine/v3.20/main/x86_64/APKINDEX.tar.gz
test_mode: false
output_file: graph.png
max_depth: 3
```
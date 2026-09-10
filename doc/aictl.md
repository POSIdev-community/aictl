# aictl

**aictl** (Application Inspector ConTroL) — командная строка для управления PT Application Inspector (AIE): проекты, ветки, исходники, сканы, отчёты и локальный context подключения.

## Зачем нужен

Вместо одноразового запуска (как **aisa**) или монолитного JAR (**ptai-cli-plugin**) aictl даёт **отдельные команды**. Их удобно собирать в скрипт или CI job: создать проект, загрузить код, запустить скан, дождаться результата, скачать отчёт.

## Поддерживаемые версии AIE

aictl работает с серверами Application Inspector в диапазоне:

**5.0.0 ≤ версия < 7.0.0**

Клиент сам выбирает API-слой под версию сервера (5.x, 6.0.x, 6.1–6.2 или 6.3+).

## Установка

1. Скачайте архив с [GitHub Releases](https://github.com/POSIdev-community/aictl/releases) под вашу ОС и архитектуру.
2. Распакуйте и положите бинарник `aictl` в каталог из `PATH` (на Unix: `chmod +x aictl`).

Имена артефактов (подставьте номер версии вместо `<ver>`):

| Платформа | Архив |
|-----------|--------|
| macOS Intel | `aictl_<ver>_darwin-amd64.tar.gz` |
| macOS Apple Silicon | `aictl_<ver>_darwin-arm64.tar.gz` |
| Linux x86_64 | `aictl_<ver>_linux-amd64.tar.gz` |
| Linux ARM64 | `aictl_<ver>_linux-arm64.tar.gz` |
| Windows x86_64 | `aictl_<ver>_windows-amd64.zip` |

Проверка:

```bash
aictl --version
```

### Токен доступа

Токен создаётся в AIE: **Администрирование → Токены доступа**. Обязательно отметьте галочку **токен для CI/CD**.

URI сервера и токен задаются через `aictl ctx set` или флаги `-u` / `-t` на connection-командах.

## Автодополнение

```bash
# текущая сессия (bash)
source <(aictl completion bash)

# Linux, постоянно
aictl completion bash > /etc/bash_completion.d/aictl

# zsh (пример)
eval "$(aictl completion zsh)"
```

Подробные флаги — в справочнике [`aictl completion`](#aictl-completion).

## Context

Локальные настройки хранятся в `~/.config/aictl/context.yaml`: URI, токен, TLS skip, путь к CA (`caCert`), id проекта и ветки по умолчанию.

```bash
aictl ctx set -u https://ai.example -t "$TOKEN"
aictl ctx set --cacert /path/to/corp-ca.pem   # при внутреннем CA
# либо осознанно: aictl ctx set --tls-skip
aictl ctx set -p "$project_id" -b "$branch_id"
aictl ctx show
aictl ctx unset --cacert
aictl ctx unset -p
aictl ctx clear -y
```

Флаги `-u`, `-t`, `--tls-skip`, `--cacert` на командах работы с сервером **переопределяют** значения из context на время одного вызова. `--cacert` и `--tls-skip` вместе нельзя. Многие команды также принимают `-p` / `-b`.

По умолчанию TLS-проверка включена (system trust). Если сервер с внутренним CA — передайте PEM через `--cacert` (сертификаты **добавляются** к system roots). Текст сертификата в context не хранится, только путь; сброс: `aictl ctx unset --cacert`.

Типичные общие флаги connection-команд: `-u`, `-t`, `--tls-skip`, `--cacert`, `-v` / `--verbose`, `-V` / `--debug`, `-l` / `--log-path`.

### Логирование (`-v`, `-V`, `-l`)

При ошибке API на stderr всегда одно **user-facing** сообщение (без Go-цепочки обёрток). Уровни:

| Флаги | Консоль | Файл (`-l`) |
|-------|---------|-------------|
| (нет) | только user-facing при ошибке | — |
| `-v` / `--verbose` | + operational-логи (`StdErr`) | — |
| `-V` / `--debug` | + operational-логи и debug-цепочка ошибки | — |
| `-v` и `-V` вместе | максимальный уровень (как `-V`) | — |
| `-l` / `--log-path` | без `-v`/`-V` консоль тихая (только user-facing при ошибке) | Error: ops и user-facing (с timestamp); **без** обычного stdout и **без** debug-цепочки |
| `-l` + `-V` | как `-V` | Error + Debug (в т.ч. цепочка) |

Обычный вывод команд (stdout / Info) в лог-файл не пишется.

## Типовой пайплайн

Укороченный сценарий (полный пример — [`examples/base-pipeline.sh`](../examples/base-pipeline.sh)):

```bash
aictl ctx set -u https://ai.example -t "$TOKEN" --tls-skip

project_id=$(aictl create project MyApp --safe)
aictl ctx set -p "$project_id"

# при наличии aiproj:
# aictl set project settings -f ./aiproj.json

branch_id=$(aictl create branch default --safe)
aictl ctx set -b "$branch_id"

aictl update sources ./src
scan_id=$(aictl scan branch "$branch_id")
aictl scan await "$scan_id" --fail-on-scan-failed
aictl scan check-policies "$scan_id" --fail-on-policies-rejected
aictl get scan report sarif "$scan_id" -o ./out/sarif.json --include-glossary --localization en
# с фильтрами уязвимостей (нужен ≥1 filter-флаг):
# aictl get scan report with-filters sarif "$scan_id" -o ./out/sarif-filtered.json \
#   --level-high --level-medium --status-confirmed --non-suppressed

aictl ctx clear -y
```

Полный пример с UI-подобным набором фильтров — [`examples/report-with-filters-pipeline.sh`](../examples/report-with-filters-pipeline.sh).

Укороченный сценарий (полный пример — [`examples/sbom-pipeline.sh`](../examples/sbom-pipeline.sh)):

```bash
project_id=$(aictl create sbom-project MySbom --file ./sbom.json --safe)
aictl ctx set -p "$project_id"

scan_id=$(aictl scan sbom)
aictl scan await "$scan_id"
aictl get scan report sarif "$scan_id" -o ./out/sarif.json --include-glossary --localization en

aictl ctx clear -y
```

`create sbom-project --file` сразу загружает SBOM; для повторной загрузки без пересоздания проекта — `aictl update sbom ./sbom.json`. `get projects` показывает колонку `TYPE` (`source` / `sbom`).

## Миграция с других CLI

### aisa → aictl

**Было** (один запуск):

```bash
aisa -u https://ai.example -t "$TOKEN" \
  --project-name MyApp \
  --scan-target ./src \
  --create-project \
  --branch-name default \
  --create-branch \
  --report Sarif \
  --reports-folder ./out
```

**Стало** (цепочка команд) — см. пайплайн выше.

Подробная таблица флагов: [doc/migration/aisa-to-aictl.md](migration/aisa-to-aictl.md).

### ptai-cli-plugin → aictl

**Было:**

```bash
java -jar ptai-cli-plugin.jar check-server --url … -t …
```

**Стало:**

```bash
aictl get healthcheck -u … -t …
aictl get version -u … -t …
```

Сценарий `ui-ast` (upload → scan → await → report) собирается теми же командами, что в типовом пайплайне.

Подробности: [doc/migration/ptai-cli-plugin-to-aictl.md](migration/ptai-cli-plugin-to-aictl.md), сводка пробелов: [doc/migration/gap-analysis.md](migration/gap-analysis.md).

## Коды выхода

| Код | Когда |
|-----|--------|
| **0** | Успех |
| **1** | Ошибка валидации входных данных, «fail»-условия (`--fail-on-scan-failed`, `--fail-on-policies-rejected` и т.п.), пустой ответ, где это считается ошибкой сценария |
| **2** | Ошибки API / сети / аутентификации / авторизации / not found / ответ сервера |
| **-1** | Неклассифицированная ошибка |

## Справочник команд

Имена команд и флагов — как в CLI. Описания — на русском. У каждой команды есть пример.

Общие inherited-флаги у connection-команд обычно включают `-u`, `-t`, `--tls-skip`, `--cacert`, `-v`, `-l`.

### `aictl`

Корневая команда CLI. Без подкоманды показывает справку; с `--version` — версию клиента.

**Usage:**

```
aictl [flags]
```

**Пример:**

```bash
aictl --version
```

**Флаги:**

```
  -h, --help      справка по aictl
      --version   показать версию aictl
```

### `aictl completion`

Генерация скриптов автодополнения для оболочки.

**Usage:**

```
aictl completion [command]
```

**Пример:**

```bash
aictl completion bash
```

**Флаги:**

```
  -h, --help   справка
```

### `aictl completion bash`

Скрипт автодополнения для bash.

**Usage:**

```
aictl completion bash [flags]
```

**Пример:**

```bash
aictl completion bash > /etc/bash_completion.d/aictl
```

**Флаги:**

```
  -h, --help              справка
      --no-descriptions   отключить описания в автодополнении
```

### `aictl completion zsh`

Скрипт автодополнения для zsh.

**Usage:**

```
aictl completion zsh [flags]
```

**Пример:**

```bash
eval "$(aictl completion zsh)"
```

**Флаги:**

```
  -h, --help              справка
      --no-descriptions   отключить описания в автодополнении
```

### `aictl completion fish`

Скрипт автодополнения для fish.

**Usage:**

```
aictl completion fish [flags]
```

**Пример:**

```bash
aictl completion fish > ~/.config/fish/completions/aictl.fish
```

**Флаги:**

```
  -h, --help              справка
      --no-descriptions   отключить описания в автодополнении
```

### `aictl completion powershell`

Скрипт автодополнения для PowerShell.

**Usage:**

```
aictl completion powershell [flags]
```

**Пример:**

```bash
aictl completion powershell | Out-String | Invoke-Expression
```

**Флаги:**

```
  -h, --help              справка
      --no-descriptions   отключить описания в автодополнении
```

### `aictl ctx`

Управление локальным context (подключение и id проекта/ветки по умолчанию).

**Usage:**

```
aictl ctx [flags]
```

**Пример:**

```bash
aictl ctx show
```

**Флаги:**

```
  -h, --help   справка
```

### `aictl ctx set`

Записать поля в `~/.config/aictl/context.yaml`. Нужен хотя бы один флаг; `--tls-skip` и `--no-tls-skip` вместе нельзя; `--cacert` и `--tls-skip` вместе нельзя. Сброс пути CA: `aictl ctx unset --cacert` (не через пустой `--cacert`).

**Usage:**

```
aictl ctx set [flags]
```

**Пример:**

```bash
aictl ctx set -u https://ai.example -t "$TOKEN"
```

**Флаги:**

```
  -b, --branch-id string    id ветки по умолчанию
      --cacert string       путь к PEM с CA (добавляется к system roots)
  -h, --help                справка
      --no-tls-skip         требовать проверку TLS-сертификата
  -p, --project-id string   id проекта по умолчанию
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI
  -u, --uri string          URI AI-сервера
```

### `aictl ctx show`

Показать текущий context (по умолчанию в удобном виде; опционально JSON/YAML).

**Usage:**

```
aictl ctx show [flags]
```

**Пример:**

```bash
aictl ctx show --yaml
```

**Флаги:**

```
  -h, --help   справка
      --json   вывести context в JSON
      --yaml   вывести context в YAML
```

### `aictl ctx unset`

Сбросить отдельные поля context (нужен хотя бы один флаг).

**Usage:**

```
aictl ctx unset [flags]
```

**Пример:**

```bash
aictl ctx unset -p -b
```

**Флаги:**

```
  -b, --branch-id    сбросить id ветки
      --cacert       сбросить путь к CA
  -h, --help         справка
  -p, --project-id   сбросить id проекта
      --tls-skip     сбросить настройку TLS skip
  -t, --token        сбросить токен
  -u, --uri          сбросить URI
```

### `aictl ctx clear`

Очистить весь локальный context (с подтверждением, если без `-y`).

**Usage:**

```
aictl ctx clear [flags]
```

**Пример:**

```bash
aictl ctx clear -y
```

**Флаги:**

```
  -h, --help   справка
  -y, --yes    не спрашивать подтверждение
```

### `aictl create`

Создание ресурсов на сервере (проект, ветка). Общий флаг `--safe` — не падать, если ресурс уже есть.

**Usage:**

```
aictl create [flags]
```

**Пример:**

```bash
aictl create project MyApp --safe
```

**Флаги:**

```
  -h, --help              справка
  -l, --log-path string   путь к файлу логов
      --safe              если ресурс уже есть — вернуть его id без ошибки
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl create project`

Создать проект по имени. Имя можно передать аргументом или через stdin. Печатает id проекта.

**Usage:**

```
aictl create project <project-name> [flags]
```

**Пример:**

```bash
aictl create project MyApp --safe
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --safe              если ресурс уже есть — вернуть его id без ошибки
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl create branch`

Создать ветку в проекте. Опционально сразу упаковать и загрузить исходники (`-s`).

**Usage:**

```
aictl create branch <branch-name> [flags]
```

**Пример:**

```bash
aictl create branch default --safe
```

**Флаги:**

```
  -e, --exclude stringArray        исключить путь (gitignore); можно повторять
      --exclude-from stringArray   файл с исключениями в стиле gitignore
  -h, --help                       справка
  -p, --project-id string          id проекта (переопределяет context)
  -s, --scan-target string         путь к исходникам для упаковки и загрузки
      --temp-dir string            каталог для временного zip при упаковке исходников
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --safe              если ресурс уже есть — вернуть его id без ошибки
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get`

Чтение данных с сервера: проекты, ветки, сканы, отчёты, здоровье и версия.

**Usage:**

```
aictl get [flags]
```

**Пример:**

```bash
aictl get healthcheck
```

**Флаги:**

```
  -h, --help              справка
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get healthcheck`

Проверить доступность AI-сервера.

**Usage:**

```
aictl get healthcheck [flags]
```

**Пример:**

```bash
aictl get healthcheck
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get version`

Показать версию Application Inspector на сервере.

**Usage:**

```
aictl get version [flags]
```

**Пример:**

```bash
aictl get version
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get agents`

Список агентов сканирования. `-q` — только id.

**Usage:**

```
aictl get agents [flags]
```

**Пример:**

```bash
aictl get agents -q
```

**Флаги:**

```
  -h, --help    справка
  -q, --quite   вывести только id
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get projects`

Список проектов, опционально отфильтрованных regexp. `-q` — только id.

**Usage:**

```
aictl get projects <regex> [flags]
```

**Пример:**

```bash
aictl get projects 'MyApp'
```

**Флаги:**

```
  -h, --help    справка
  -q, --quite   вывести только id
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get project`

Группа команд по одному проекту (id из context или `-p`).

**Usage:**

```
aictl get project [flags]
```

**Пример:**

```bash
aictl get project settings
```

**Флаги:**

```
  -h, --help                справка
  -p, --project-id string   id проекта (переопределяет context)
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get project settings`

Настройки сканирования проекта. `--json` — вывод в JSON.

**Usage:**

```
aictl get project settings [flags]
```

**Пример:**

```bash
aictl get project settings --json
```

**Флаги:**

```
  -h, --help   справка
      --json   вывод в JSON
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl get project policies`

Политики качества проекта.

**Usage:**

```
aictl get project policies [flags]
```

**Пример:**

```bash
aictl get project policies
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl get project exclusions`

Исключения проекта (пути/шаблоны).

**Usage:**

```
aictl get project exclusions [flags]
```

**Пример:**

```bash
aictl get project exclusions
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl get project aiproj`

Скачать текущие настройки проекта в формате aiproj (файл через `-o`).

**Usage:**

```
aictl get project aiproj [flags]
```

**Пример:**

```bash
aictl get project aiproj -o ./project.aiproj -f
```

**Флаги:**

```
  -f, --force           перезаписать существующий файл
  -h, --help            справка
  -o, --output string   путь к выходному файлу
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl get branches`

Список веток проекта. Опциональный regexp; `-q` — только id.

**Usage:**

```
aictl get branches [<regex>] [flags]
```

**Пример:**

```bash
aictl get branches -p <project-id>
```

**Флаги:**

```
  -h, --help                справка
  -p, --project-id string   id проекта (переопределяет context)
  -q, --quite               вывести только id
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get branch`

Информация о ветке по id.

**Usage:**

```
aictl get branch <branch-id> [flags]
```

**Пример:**

```bash
aictl get branch <branch-id>
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get scans`

Список сканов (опционально regexp по метке). `--latest` — только последний; `-q` — только id.

**Usage:**

```
aictl get scans [<regex>] [flags]
```

**Пример:**

```bash
aictl get scans --latest -q
```

**Флаги:**

```
  -b, --branch-id string   id ветки (переопределяет context)
  -h, --help               справка
      --latest             вернуть только последний результат скана
  -q, --quite              вывести только id
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get scanning`

Текущие выполняющиеся сканы. `-p` фильтрует по проекту (значение из ctx игнорируется).

**Usage:**

```
aictl get scanning [flags]
```

**Пример:**

```bash
aictl get scanning -p <project-id>
```

**Флаги:**

```
  -h, --help                справка
  -p, --project-id string   фильтр по id проекта (ctx -p игнорируется)
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get queue`

Очередь сканов. `-p` фильтрует по проекту (ctx `-p` игнорируется).

**Usage:**

```
aictl get queue [flags]
```

**Пример:**

```bash
aictl get queue
```

**Флаги:**

```
  -h, --help                справка
  -p, --project-id string   фильтр по id проекта (ctx -p игнорируется)
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get report-templates`

Шаблоны отчётов на сервере. Опциональный regexp; `-q` — только id.

**Usage:**

```
aictl get report-templates [<regex>] [flags]
```

**Пример:**

```bash
aictl get report-templates --localization ru
```

**Флаги:**

```
  -h, --help                  справка
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -q, --quite                 вывести только id
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get scan`

Группа команд по конкретному скану (нужен `<scan-id>` у подкоманд).

**Usage:**

```
aictl get scan <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan stage <scan-id>
```

**Флаги:**

```
  -h, --help                справка
  -p, --project-id string   id проекта (переопределяет context)
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl get scan stage`

Текущая стадия скана. С `--fail-on-scan-failed` — exit 1 при Failed/Aborted.

**Usage:**

```
aictl get scan stage <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan stage <scan-id> --fail-on-scan-failed
```

**Флаги:**

```
      --fail-on-scan-failed   код 1, если стадия скана Failed или Aborted
  -h, --help                  справка
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl get scan errors`

Ошибки скана.

**Usage:**

```
aictl get scan errors <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan errors <scan-id>
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl get scan logs`

Логи скана (в файл через `-o`).

**Usage:**

```
aictl get scan logs <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan logs <scan-id> -o ./scan.log -f
```

**Флаги:**

```
  -f, --force           перезаписать существующий файл
  -h, --help            справка
  -o, --output string   путь к выходному файлу
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl get scan statistic`

Статистика результатов скана. `--json` — JSON; `-o` — в файл.

**Usage:**

```
aictl get scan statistic <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan statistic <scan-id> --json
```

**Флаги:**

```
  -f, --force           перезаписать существующий файл
  -h, --help            справка
      --json            вывод в JSON
  -o, --output string   путь к выходному файлу
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl get scan sbom`

SBOM по результатам скана (файл через `-o`).

**Usage:**

```
aictl get scan sbom <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan sbom <scan-id> -o ./sbom.json -f
```

**Флаги:**

```
  -f, --force           перезаписать существующий файл
  -h, --help            справка
  -o, --output string   путь к выходному файлу
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl get scan aiproj`

Настройки проекта, зафиксированные на момент скана (aiproj).

**Usage:**

```
aictl get scan aiproj <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan aiproj <scan-id> -o ./scan.aiproj -f
```

**Флаги:**

```
  -f, --force           перезаписать существующий файл
  -h, --help            справка
  -o, --output string   путь к выходному файлу
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl get scan report`

Скачать отчёт по скану в выбранном формате. Форматы — подкоманды ниже; общие флаги наследуются.
Без фильтров уязвимостей (`useFilters=false`). Для фильтров — [`get scan report with-filters`](#aictl-get-scan-report-with-filters).

**Usage:**

```
aictl get scan report <report-name> <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report sarif <scan-id> -o ./out.sarif
```

**Флаги:**

```
  -f, --force                 перезаписать существующий файл
  -h, --help                  справка
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -o, --output string         путь к выходному файлу
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl get scan report with-filters`

То же, что `get scan report`, но с фильтрами уязвимостей (`useFilters=true`). Нужен **хотя бы один** filter-флаг.
Булевы флаги при наличии отправляют `true`, иначе поле omit; `--type` / `--language` / `--scan-module` без передачи уходят как пустые массивы `[]`.
`SecretDetection` / `MaliciousCodeDetection` / `OneC` требуют AIE ≥ 6.0; `Dart` — ≥ 6.1.

**Usage:**

```
aictl get scan report with-filters <report-name> <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report with-filters sarif <scan-id> -o ./out.sarif --level-high --level-medium --non-suppressed
```

**Filter-флаги:**

```
      --level-high|--level-medium|--level-low|--level-potential
      --status-undefined|--status-confirmed|--status-confirmed-auto|--status-rejected
      --mode-entry-point|--mode-public-methods|--mode-root-function|--mode-others
      --found-this-scan|--found-prev-scan
      --conditional|--non-conditional|--suppressed|--non-suppressed
      --suspected|--second-level|--only-favorite
      --type string                 (повторяемый)
      --scan-module string          whitelist API: StaticCodeAnalysis, PatternMatching, …
      --language string             без None; имена AIE (Java, Dart, …)
```

Остальные флаги (`-o`, `-f`, `--include-*`, `--localization`) — как у `get scan report`. Форматные подкоманды те же (`sarif`, `json`, …).

### `aictl get scan report autocheck`

Отчёт Autocheck по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report autocheck <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report autocheck <scan-id> -o ./out.autocheck -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report gitlab`

Отчёт для GitLab по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report gitlab <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report gitlab <scan-id> -o ./out.gitlab -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report json`

Отчёт JSON по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report json <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report json <scan-id> -o ./out.json -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report json-v2`

Отчёт JSON v2 по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report json-v2 <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report json-v2 <scan-id> -o ./out.json-v2 -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report markdown`

Отчёт Markdown по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report markdown <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report markdown <scan-id> -o ./out.markdown -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report nist`

Отчёт NIST по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report nist <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report nist <scan-id> -o ./out.nist -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report oud4`

Отчёт OUD4 по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report oud4 <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report oud4 <scan-id> -o ./out.oud4 -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report owasp`

Отчёт OWASP по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report owasp <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report owasp <scan-id> -o ./out.owasp -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report owaspm`

Отчёт OWASP Mobile по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report owaspm <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report owaspm <scan-id> -o ./out.owaspm -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report pcidss`

Отчёт PCI DSS по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report pcidss <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report pcidss <scan-id> -o ./out.pcidss -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report plain`

Текстовый отчёт по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report plain <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report plain <scan-id> -o ./out.plain -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report sans`

Отчёт SANS по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report sans <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report sans <scan-id> -o ./out.sans -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report sarif`

Отчёт SARIF по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report sarif <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report sarif <scan-id> -o ./out.sarif -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl get scan report xml`

Отчёт XML по указанному скану. Общие флаги — у родительской `get scan report`.

**Usage:**

```
aictl get scan report xml <scan-id> [flags]
```

**Пример:**

```bash
aictl get scan report xml <scan-id> -o ./out.xml -f
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -f, --force                 перезаписать существующий файл
      --include-comments      включить комментарии в отчёт
      --include-dfd           включить диаграммы потоков данных
      --include-glossary      включить глоссарий
      --localization string   локализация отчёта: en или ru (по умолчанию en)
  -l, --log-path string       путь к файлу логов
  -o, --output string         путь к выходному файлу
  -p, --project-id string     id проекта (переопределяет context)
      --tls-skip              не проверять TLS-сертификат сервера
  -t, --token string          токен доступа к AI (переопределяет context)
  -u, --uri string            URI AI-сервера (переопределяет context)
  -v, --verbose               подробный вывод
  -V, --debug               debug-вывод (цепочки ошибок)
```

### `aictl set`

Запись конфигурации ресурсов (политики, exclusions, aiproj).

**Usage:**

```
aictl set [flags]
```

**Пример:**

```bash
aictl set project settings -f ./aiproj.json
```

**Флаги:**

```
  -h, --help              справка
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl set project`

Группа set-команд для проекта (id из context или `-p`).

**Usage:**

```
aictl set project [flags]
```

**Пример:**

```bash
aictl set project policies -f ./policies.json
```

**Флаги:**

```
  -h, --help                справка
  -p, --project-id string   id проекта (переопределяет context)
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl set project settings`

Загрузить настройки проекта из aiproj/JSON файла (`-f`).

**Usage:**

```
aictl set project settings [flags]
```

**Пример:**

```bash
aictl set project settings -f ./aiproj.json
```

**Флаги:**

```
  -f, --file string   путь к файлу aiproj.json
  -h, --help          справка
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl set project policies`

Задать политики качества (JSON файл, `-` или аргумент).

**Usage:**

```
aictl set project policies [json|-] [flags]
```

**Пример:**

```bash
aictl set project policies -f ./policies.json
```

**Флаги:**

```
  -f, --file string   JSON файл политик или - для stdin
  -h, --help          справка
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl set project exclusions`

Задать exclusions (файл, `-` или текст аргументом).

**Usage:**

```
aictl set project exclusions [text|-] [flags]
```

**Пример:**

```bash
aictl set project exclusions -f ./exclusions.txt
```

**Флаги:**

```
  -f, --file string   файл исключений или - для stdin
  -h, --help          справка
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl update`

Обновление ресурсов: исходники, языки, настройки скана проекта.

**Usage:**

```
aictl update [flags]
```

**Пример:**

```bash
aictl update sources ./src
```

**Флаги:**

```
  -h, --help              справка
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl update sources`

Упаковать каталог исходников и загрузить на сервер в текущую (или указанную) ветку.

**Usage:**

```
aictl update sources <path> [flags]
```

**Пример:**

```bash
aictl update sources ./src -e '**/test/**'
```

**Флаги:**

```
  -b, --branch-id string           id ветки (переопределяет context)
  -e, --exclude stringArray        исключить путь (gitignore); можно повторять
      --exclude-from stringArray   файл с исключениями в стиле gitignore
  -h, --help                       справка
  -p, --project-id string          id проекта (переопределяет context)
      --temp-dir string            каталог для временного zip при упаковке исходников
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl update project`

Группа update-команд для проекта.

**Usage:**

```
aictl update project [flags]
```

**Пример:**

```bash
aictl update project settings --priority High
```

**Флаги:**

```
  -h, --help                справка
  -p, --project-id string   id проекта (переопределяет context)
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl update project settings`

Частично обновить настройки скана проекта (приоритет, агенты). Нужен хотя бы один флаг.

**Usage:**

```
aictl update project settings [flags]
```

**Пример:**

```bash
aictl update project settings --priority High --preferred-agents-only
```

**Флаги:**

```
      --agents string              предпочтительные агенты сканирования (id через запятую)
  -h, --help                       справка
      --no-preferred-agents-only   разрешить любые агенты, не только выбранные
      --preferred-agents-only      использовать только выбранные агенты
      --priority string            приоритет скана: None, Low, Medium, High или Critical
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl update project languages`

Обновить список языков проекта на сервере (по загруженным исходникам).

**Usage:**

```
aictl update project languages [flags]
```

**Пример:**

```bash
aictl update project languages
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -l, --log-path string     путь к файлу логов
  -p, --project-id string   id проекта (переопределяет context)
      --tls-skip            не проверять TLS-сертификат сервера
  -t, --token string        токен доступа к AI (переопределяет context)
  -u, --uri string          URI AI-сервера (переопределяет context)
  -v, --verbose             подробный вывод
  -V, --debug             debug-вывод (цепочки ошибок)
```

### `aictl scan`

Управление сканами: старт, ожидание, остановка, проверка политик.

**Usage:**

```
aictl scan [flags]
```

**Пример:**

```bash
aictl scan await <scan-id>
```

**Флаги:**

```
  -h, --help              справка
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl scan branch`

Запустить скан ветки (source-проекты). Печатает id скана.

**Usage:**

```
aictl scan branch <branch-id> [flags]
```

**Пример:**

```bash
aictl scan branch <branch-id> --scan-label ci-1
```

**Флаги:**

```
      --full-scan           полный скан вместо инкрементального
  -h, --help                справка
  -p, --project-id string   id проекта (переопределяет context)
      --scan-label string   метка скана (до 40 символов)
```

### `aictl scan project`

Запустить скан source-проекта. Печатает id скана. Для SBOM-проектов используйте `aictl scan sbom`.

**Usage:**

```
aictl scan project <project-id> [flags]
```

**Пример:**

```bash
aictl scan project <project-id> --full-scan
```

**Флаги:**

```
      --full-scan           полный скан вместо инкрементального
  -h, --help                справка
  -p, --project-id string   id проекта (переопределяет context)
      --scan-label string   метка скана (до 40 символов)
```

### `aictl scan sbom`

Запустить скан SBOM-проекта (AIE ≥ 6.3). Печатает id скана. Без `--full-scan` (всегда incremental).

**Usage:**

```
aictl scan sbom <project-id> [flags]
```

**Пример:**

```bash
aictl scan sbom -p <project-id> --scan-label nightly
```

**Флаги:**

```
  -h, --help                справка
  -p, --project-id string   id проекта (переопределяет context)
      --scan-label string   метка скана (до 40 символов)
```

### `aictl scan start` (устарело)

Устарело: используйте `aictl scan branch` / `aictl scan project`. Команды `scan start *` сохранены для совместимости.

### `aictl scan await`

Ждать завершения скана. С `--fail-on-scan-failed` — exit 1 при Failed/Aborted.

**Usage:**

```
aictl scan await <scan-id> [flags]
```

**Пример:**

```bash
aictl scan await <scan-id> --fail-on-scan-failed
```

**Флаги:**

```
      --fail-on-scan-failed   код 1, если стадия скана Failed или Aborted
  -h, --help                  справка
  -p, --project-id string     id проекта (переопределяет context)
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl scan check-policies`

Проверить результат политик по скану. С `--fail-on-policies-rejected` — exit 1 при Rejected.

**Usage:**

```
aictl scan check-policies <scan-id> [flags]
```

**Пример:**

```bash
aictl scan check-policies <scan-id> --fail-on-policies-rejected
```

**Флаги:**

```
      --fail-on-policies-rejected   код 1, если PolicyState = Rejected
  -h, --help                        справка
  -p, --project-id string           id проекта (переопределяет context)
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl scan stop`

Остановить выполняющийся скан.

**Usage:**

```
aictl scan stop <scan-id> [flags]
```

**Пример:**

```bash
aictl scan stop <scan-id>
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl delete`

Удаление ресурсов на сервере.

**Usage:**

```
aictl delete [flags]
```

**Пример:**

```bash
aictl delete projects <project-id>
```

**Флаги:**

```
  -h, --help              справка
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```

### `aictl delete projects`

Удалить один или несколько проектов по id (аргументы или stdin).

**Usage:**

```
aictl delete projects <project-id>... [flags]
```

**Пример:**

```bash
aictl delete projects <project-id>
```

**Флаги:**

```
  -h, --help   справка
```

**Унаследованные флаги:**

```
  -l, --log-path string   путь к файлу логов
      --tls-skip          не проверять TLS-сертификат сервера
  -t, --token string      токен доступа к AI (переопределяет context)
  -u, --uri string        URI AI-сервера (переопределяет context)
  -v, --verbose           подробный вывод
  -V, --debug           debug-вывод (цепочки ошибок)
```


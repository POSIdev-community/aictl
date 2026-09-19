# Миграция: aisa → aictl

**aisa** (AI.ScanAgent.Cli) — one-shot CLI: один запуск = connect + (опционально create) + upload + scan + wait + reports.

**aictl** — набор отдельных команд. Тот же сценарий собирается скриптом или CI job.

См. также: [gap-analysis](gap-analysis.md), [гайд aictl](../aictl.md), [автоген справочник](../gen/aictl.md), пример [`examples/base-pipeline.sh`](../../examples/base-pipeline.sh).

Команды и флаги со статусом **планируется** ещё не реализованы; имена — целевые предложения.

---

## Быстрый старт

### Было (aisa)

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

### Стало (aictl)

```bash
aictl ctx set -u https://ai.example -t "$TOKEN" --tls-skip   # --tls-skip при необходимости

project_id=$(aictl create project MyApp --safe)
aictl ctx set -p "$project_id"

# если есть .aiproj:
# aictl set project settings -f ./project.aiproj

branch_id=$(aictl create branch default --safe)
aictl ctx set -b "$branch_id"

aictl update sources ./src --update-languages
# aictl update sources ./src --exclude-from .aisaignore --update-languages

scan_id=$(aictl scan branch "$branch_id")
aictl scan await "$scan_id" --fail-on-scan-failed
aictl scan check-policies "$scan_id" --fail-on-policies-rejected

aictl get scan report sarif "$scan_id" -o ./out/sarif.json --include-glossary --localization en

aictl ctx clear -y
```

Без ожидания (аналог `--no-wait`): не вызывайте `scan await`; сохраните `scan_id` и позже вызовите `scan await` / `get scan stage`.

---

## Таблица флагов aisa → aictl

### Подключение и настройки клиента

| aisa | aictl | Примечание |
|------|-------|------------|
| `-u` | `ctx set -u` или `-u` на connection-командах | |
| `-t` | `ctx set -t` или `-t` | |
| `--set-settings` | `ctx set` | Persist в `~/.config/aictl/context.yaml` |
| `--log-level` | `-v` / `--verbose`, `-l` / `--log-path` | Нет полного паритета уровней NLog |
| `-v` / `--version` | `aictl --version` | `-v` в aictl — verbose |
| (нет TLS-skip) | `--tls-skip` | Только в aictl |

### Проект, ветка, источники

| aisa | aictl | Примечание |
|------|-------|------------|
| `--project-name` | `create project <name>`; `get projects <regex>` | Дальше работа по UUID в `ctx` / `-p` |
| `--project-id` | `-p` / `--project-id`, `ctx set -p` | |
| `--create-project` | `create project` (`--safe` — не падать, если есть) | |
| `--branch-name` | `create branch <name>`; `ctx set -b` | |
| `--create-branch` | `create branch` | |
| `--scan-target` | `update sources <path> [--update-languages]` | По умолчанию aisa: CWD; в пайплайнах — с `--update-languages` |
| `--file-exclusions` | `update sources -e` / `--exclude-from` | Синтаксис gitignore; при полном пайплайне добавьте `--update-languages` |
| `--project-settings-file` | `set project settings -f` | `.aiproj` / JSON |
| `--scan-off` | Не вызывать `scan branch` / `scan project` | Upload/settings без скана |
| `--policy-settings-file` | `set project policies -f` | |
| `--policies-path` | `set project policies -f` | |

### Скан

| aisa | aictl | Примечание |
|------|-------|------------|
| (старт в составе one-shot) | `scan branch` / `scan project` | |
| `--full-scan` | `scan branch|project --full-scan` | По умолчанию incremental |
| `--scan-label` | `scan branch|project --scan-label` | |
| `--no-wait` | Не вызывать `scan await` | |
| `--status` + `--project-id` + `--scan-result-id` | `get scan stage` / `scan await` | |
| `--scan-result-id` | аргумент `<scan-id>` | |
| `--retry` / `--retry-time` | Скрипт retry | См. [`pipeline-with-retry.sh`](../../examples/pipeline-with-retry.sh); нативный флаг не планируется как P0 |
| (ожидание до Done) | `scan await` | Без `--fail-on-scan-failed` exit 0 и на Failed/Aborted |
| (fail при ошибке скана) | `scan await --fail-on-scan-failed` | Failed/Aborted → exit 1; также на `get scan stage` |
| (policy violated → exit 10) | `scan check-policies --fail-on-policies-rejected` | exit 1 при Rejected |

### Отчёты и артефакты

| aisa | aictl | Примечание |
|------|-------|------------|
| `--reports-folder` | `-o` у `get scan report` | Каталог создавайте сами |
| `--report PlainReport` / `HTML` | `get scan report plain` | HTML в aisa ≈ PlainReport |
| `--report JSON` | `get scan report json` | Также `json-v2` |
| `--report Sarif` | `get scan report sarif` | |
| `--report Gitlab` | `get scan report gitlab` | |
| `--report MD` | `get scan report markdown` | |
| `--report AutoCheck` | `get scan report autocheck` | |
| `--report Owasp` / `Owaspm` / `Pcidss` / `Sans` / `Nist` / `Oud4` | соответствующие subcommands | |
| `--report <CustomName>` | `get scan report <name> <scan-id>` | Пользовательский шаблон |
| `--waf-patch` | `get scan waf-patch <scan-id> -o …` | **Планируется** |
| (список шаблонов) | `get report-templates` | **Планируется** (из ptai) |

### Не переносим

| aisa | Причина |
|------|---------|
| `--list-results` | Hidden stub |
| `--restore-sources` | Hidden / не реализовано |
| `--nolocalize`, `--noclearcache`, `--debug`, `--ddpython`, `--service` | Не используются в runtime |
| `--incremental` | Hidden; в aictl incremental по умолчанию, полный — `--full-scan` |
| Single-instance mutex | Не переносим |
| Пути `~/aisa/appSettings.user.json` | Замена: `ctx` / `context.yaml` |

---

## Exit codes

aisa использует детальные коды. aictl сохраняет схему **0 / 1 / 2** (и **-1** для неизвестных ошибок). Ниже — ориентир для миграции скриптов, не бинарная совместимость.

| Код aisa (примеры) | Смысл | В aictl |
|--------------------|-------|---------|
| 0 | OK | 0 |
| 10 | Policy violated | `scan check-policies --fail-on-policies-rejected` → **1** |
| 60 | Policy + non-critical | То же / смотреть stderr |
| 2 | Scan target not found | **1** (validation) |
| 4 | Project not found | **2** (API) / **1** |
| 5, 8, 11, … | Невалидные settings/reports | **1** |
| 25 | Invalid URI | **1** |
| 29 | Invalid token | **2** |
| 31 | Branch scan already running | **2** / retry-скрипт |
| 32 | Connection failed | **2** |
| 100 | Cancelled | зависит от сигнала / ctx cancel |
| 1000 | Unknown (aisa) | **-1** |
| -1 | Another instance | Не применимо |

Полный список кодов aisa живёт в исходниках ScanAgent (`ExitCodes`); в aictl ориентируйтесь на текст ошибки в stderr и общую схему 0/1/2.

---

## Типичные сценарии

### Только загрузить settings и sources без скана

```bash
aictl set project settings -f ./project.aiproj
aictl update sources ./src --update-languages
# не вызывать scan branch  → аналог --scan-off
```

### Полный скан с меткой

```bash
aictl scan branch "$branch_id" --full-scan --scan-label "ci-$CI_COMMIT_SHA"
```

### Проверка статуса позже

```bash
aictl get scan stage "$scan_id"
# позже:
aictl scan await "$scan_id" --fail-on-scan-failed
```

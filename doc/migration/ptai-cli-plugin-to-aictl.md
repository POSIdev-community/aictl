# Миграция: ptai-cli-plugin → aictl

**ptai-cli-plugin** (`java -jar ptai-cli-plugin.jar` / Docker) — CLI с подкомандами `ui-ast`, `json-ast`, `check-server`, `list-report-templates`, `generate-report`, `delete-project`.

**aictl** заменяет этот CLI отдельными командами. Плагины **Jenkins** (`ptaiAst`) и **TeamCity** (`ptsecurity`) из репозитория `ptai-ee-tools` **не** входят в этот гайд (отдельный трек).

См. также: [gap-analysis](gap-analysis.md), [миграция aisa](aisa-to-aictl.md), [гайд aictl](../aictl.md), [автоген справочник](../gen/aictl.md), [`examples/base-pipeline.sh`](../../examples/base-pipeline.sh).

Статус **планируется** — предлагаемое API, ещё не в CLI.

---

## Общие флаги подключения

| ptai-cli-plugin | aictl | Примечание |
|-----------------|-------|------------|
| `--url` | `-u` / `--uri`, `ctx set -u` | |
| `-t` / `--token` | `-t` / `--token`, `ctx set -t` | |
| `--user` / `--password` | — | **Не планируется** (только token) |
| `--insecure` | `--tls-skip` | |
| `--truststore` | — | **Не планируется** (системный trust / `--tls-skip`) |
| `-v` / `--verbose` | `-v` / `--verbose` | |
| `--version` | `aictl --version` | |

---

## `check-server` → aictl

```bash
# было
java -jar ptai-cli-plugin.jar check-server --url … -t …

# стало
aictl get healthcheck -u … -t …
aictl get version -u … -t …
```

---

## `ui-ast` → aictl

Настройки проекта уже на сервере (UI). Нужно: найти/задать project+branch → upload → start → await → reports.

```bash
# было (схема)
java -jar ptai-cli-plugin.jar ui-ast \
  --url … -t … \
  -p MyApp --input ./src \
  -b default --scan-label ci-1 \
  --full-scan --async=false \
  --fail-if-failed \
  --sarif-report-file ./out/sarif.json

# стало
aictl ctx set -u … -t …

# project id: из ctx или через список
# aictl get projects 'MyApp'   # regexp; возьмите UUID
aictl ctx set -p "$project_id"
# ветка:
branch_id=$(aictl create branch default --safe)   # или get branches
aictl ctx set -b "$branch_id"

aictl update sources ./src
# excludes: -e / --exclude-from (gitignore; не Ant)

# priority на проекте (если нужно):
# aictl update project settings --priority High

scan_id=$(aictl scan start branch "$branch_id" --full-scan --scan-label ci-1)
aictl scan await "$scan_id" --fail-on-scan-failed
aictl scan check-policies "$scan_id" --fail-on-policies-rejected
aictl get scan report sarif "$scan_id" -o ./out/sarif.json
```

### Флаги `ui-ast`

| ptai | aictl | Примечание |
|------|-------|------------|
| `-p` / `--project` | `ctx -p` / UUID после `get projects` | Точный `get project --name` — **планируется** |
| `--input` | `update sources <path>` | |
| `-b` / `--branch-name` | `create branch` / `ctx -b` | Ветка создаётся, если нет |
| `--scan-label` | `scan start --scan-label` | |
| `--output` | каталог для `-o` | Дефолт ptai: `.ptai` |
| `-i` / `--includes` | staging dir или `update sources --include` | `--include` — **планируется**; сейчас — подготовить дерево файлов |
| `-e` / `--excludes` | `-e` / `--exclude-from` | **gitignore**, не Ant |
| `--use-default-excludes` | — | Эмулировать списком excludes |
| `--fail-if-failed` | `scan check-policies --fail-on-policies-rejected` | |
| `--fail-if-unstable` | — | Нет прямого аналога; смотреть statistic / policy |
| `--async` | не вызывать `scan await` | |
| `--full-scan` | `scan start --full-scan` | |
| `--priority` | `update project settings --priority` или `scan start --priority` | На start — **планируется** |
| Reporting flags | см. ниже | |

---

## `json-ast` → aictl

Как `ui-ast`, но settings (и опционально policy) из JSON:

```bash
aictl set project settings -f ./settings.json     # --settings-json
aictl set project policies -f ./policy.json       # --policy-json
aictl update sources ./src
scan_id=$(aictl scan start branch "$branch_id")
aictl scan await "$scan_id" --fail-on-scan-failed
aictl scan check-policies "$scan_id" --fail-on-policies-rejected
```

| ptai | aictl |
|------|-------|
| `--settings-json` | `set project settings -f` |
| `--policy-json` | `set project policies -f`; `[]` для очистки — уточнится при использовании |
| остальные флаги | как у `ui-ast` |

Создание проекта при отсутствии: `create project` + `set project settings -f` (в ptai `json-ast` создаёт/обновляет из AIPROJ).

---

## `list-report-templates` → aictl

```bash
# было
java -jar ptai-cli-plugin.jar list-report-templates -l EN --url … -t …

# стало
aictl get report-templates --localization en
```

---

## `generate-report` → aictl

Отчёты по уже существующему результату скана:

```bash
# было: --project-name / --project-id, --scan-result-id, reporting flags

aictl ctx set -p "$project_id"
# последний скан:
# aictl get scans --latest
aictl get scan report sarif "$scan_id" -o ./out/sarif.json
aictl get scan report plain "$scan_id" -o ./out/report.html --include-dfd --include-glossary --localization en
```

| ptai | aictl |
|------|-------|
| `--project-name` / `--project-id` | UUID + `ctx` / `-p` |
| `-b` / `--branch-name` | `ctx -b` / фильтр `get scans` |
| `--scan-result-id` | `<scan-id>`; иначе `get scans --latest` |
| `--output` | пути в `-o` |
| reporting | см. ниже |

---

## `delete-project` → aictl

```bash
# было
java -jar ptai-cli-plugin.jar delete-project --project-name MyApp -y --url … -t …
# или --project-id / --project-name-regexp

# сейчас
aictl get projects 'MyApp'          # найти UUID
aictl delete projects "$project_id"

# станет
aictl delete projects --name MyApp -y              # планируется
aictl delete projects --regexp 'e2e-.*' -y         # планируется
```

| ptai | aictl |
|------|-------|
| `--project-id` | `delete projects <uuid>…` |
| `--project-name` | `--name` (**планируется**) |
| `--project-name-regexp` | `--regexp` (**планируется**) |
| `-y` / `--yes` | для name/regexp — **планируется**; для UUID подтверждение по поведению текущей команды |

---

## Reporting

| ptai | aictl | Примечание |
|------|-------|------------|
| `--report-template` + `--report-file` | `get scan report <template-or-type> -o` | |
| `--report-include-dfd` | `--include-dfd` | |
| `--report-include-glossary` | `--include-glossary` | |
| `-l` / `--locale` | `--localization en\|ru` | |
| `--raw-data-file` | `get scan report json` / `json-v2` или `raw` | `raw` — **планируется**, если json недостаточен |
| `--sarif-report-file` | `get scan report sarif -o` | |
| `--giif-report-file` | `get scan report giif -o` | **Планируется** |
| `--report-json` | `get scan reports -f reports.json` | **Планируется**; workaround — несколько вызовов |

---

## Exit codes

| ptai | aictl |
|------|-------|
| 0 | 0 |
| 1 (в т.ч. fail-if-*) | 1 при `--fail-on-scan-failed` / `--fail-on-policies-rejected`; иначе смотреть stage вручную |
| 1000 (bad args) | 1 (validation) |

Схема aictl: **0 / 1 / 2** (без кода 1000).

---

## Что не переносится в aictl

- Jenkins / TeamCity UI и Pipeline DSL
- Advanced JVM `-Dptai.*` (таймауты, diagnostic filenames) — при необходимости отдельные флаги/env позже
- Ant-семантика includes/excludes 1:1 (gitignore + staging)
- Auth login/password и PEM truststore

---

## Чеклист миграции CI (GitLab / generic)

1. Заменить `java -jar ptai-cli-plugin.jar …` на `aictl` (бинарник или свой image).
2. Вынести `-u`/`-t` в `aictl ctx set` в начале job.
3. `ui-ast` / `json-ast` разбить на: settings (если json) → sources → start → await → policy check → reports.
4. `--async` → не вызывать await; сохранить `scan_id` в артефакт.
5. `--fail-if-failed` → `scan check-policies --fail-on-policies-rejected`.
6. Удаление тестовых проектов: UUID сейчас; `--name`/`--regexp` — после P0.

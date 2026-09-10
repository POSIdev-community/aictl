# Gap-анализ: aisa / ptai-cli-plugin → aictl

Документ фиксирует, что уже покрыто **aictl**, чего не хватает для замены **aisa** и **ptai-cli-plugin** в CI, и какие имена команд предлагаются для доработок.

Связанные гайды:

- [Миграция aisa → aictl](aisa-to-aictl.md)
- [Миграция ptai-cli-plugin → aictl](ptai-cli-plugin-to-aictl.md)
- [Backlog доработок](backlog.md)
- [Гайд aictl](../aictl.md)
- [Автоген справочник команд](../gen/aictl.md)

**Модель aictl:** набор отдельных команд (не one-shot). Типичный пайплайн собирается скриптом — см. [`examples/base-pipeline.sh`](../../examples/base-pipeline.sh).

**Вне scope этого анализа:** плагины Jenkins / TeamCity из `ptai-ee-tools` (отдельный трек).

Статусы в таблицах:

| Статус | Значение |
|--------|----------|
| Есть | Реализовано в aictl |
| Планируется | Предлагаемое API; ещё не в CLI |
| Workaround | Можно закрыть существующими командами / скриптом |
| Не планируется | Намеренно не переносим |

---

## Уже покрыто

| Возможность | aisa / ptai-cli-plugin | aictl |
|-------------|------------------------|-------|
| Подключение и сохранение | `-u`/`-t`, `--set-settings` | `ctx set` / `show` / `clear` / `unset`; флаги `-u`/`-t`/`--tls-skip`/`--cacert` |
| Health / версия сервера | `check-server` | `get healthcheck`, `get version` |
| Создание проекта | `--create-project` | `create project` |
| Создание ветки | `--create-branch`, `--branch-name` | `create branch` |
| Настройки из `.aiproj` / JSON | `--project-settings-file`, `--settings-json` | `set project settings -f` |
| Загрузка исходников | `--scan-target` / `--input` | `update sources <path>` |
| Исключения при upload | `--file-exclusions` (gitignore) | `-e` / `--exclude-from` (gitignore) |
| Старт скана | (часть one-shot / `ui-ast`) | `scan branch` / `scan project` |
| Метка скана | `--scan-label` | `--scan-label` на `scan branch` / `scan project` |
| Полный скан | `--full-scan` | `--full-scan` на `scan branch` / `scan project` |
| Async / без ожидания | `--no-wait`, `--async` | `scan branch` / `scan project` без `scan await` |
| Ожидание скана | по умолчанию / `--status` | `scan await` |
| Остановка скана | — | `scan stop` |
| Отчёты (основные типы) | `--report`, SARIF/HTML/… | `get scan report <type>` (+ шаблон по имени) |
| Priority проекта | `--priority` (на AST в ptai) | `update project settings --priority` |
| Списки проектов / веток / сканов | частично | `get projects` / `branches` / `scans` (`--latest`) |
| Удаление проекта | `delete-project` | `delete projects <uuid>…`; по имени: `get projects <regex>` → UUID |
| Поиск проекта по имени | `-p` / `--project-name` | `get projects <regex>` |
| Агенты, SBOM, логи, статистика | — / частично | `get agents`, `get scan sbom` / `logs` / `statistic` / `stage` |
| Fail при Failed/Aborted | ненулевой exit / `--fail-if-*` | `scan await --fail-on-scan-failed`; `get scan stage --fail-on-scan-failed` |
| Проверка политики скана | aisa exit 10 / `--fail-if-failed` | `scan check-policies [--fail-on-policies-rejected]` |
| Правила политик проекта | `--policies-path` / `--policy-json` | `get` / `set project policies` |

---

## Gaps P0 (блокеры типичного CI)

| Gap | Источник | Предлагаемое API aictl | Workaround сейчас | Статус |
|-----|----------|------------------------|-------------------|--------|
| Список шаблонов отчётов | ptai `list-report-templates` | `get report-templates [<regex>] [-q]` (id+имя) | Знать имя заранее | Есть |
| Ошибки скана | — | `get scan errors` (построчный текст) | Нет | Есть |
| Актуализация языков | — | `update project languages` (detect по загруженным sources) | Settings вручную | Есть |
| Исключения проекта | — | `get`/`set project exclusions` (get/set: gitignore-текст; set — полная замена) | Upload `-e` на клиенте | Есть |
| Очередь сканов | — | `get queue` [`-p`]; без `-p` ctx не фильтрует; таблица | Нет | Есть |
| Активные сканы | — | `get scanning` [`-p`]; таблица | Нет | Есть |
| Каталог temp-архива | — | `--temp-dir` на `update sources` / `create branch`; удалить zip после upload | Системный temp | Есть |

**Поведение fail-флагов:** по умолчанию выключены. `--fail-on-scan-failed` и `--fail-on-policies-rejected` — разные команды; policy gate **не** вешается на `scan await`.

Пример целевого CI:

```bash
sid=$(aictl scan branch "$bid")
aictl scan await "$sid" --fail-on-scan-failed
aictl scan check-policies "$sid" --fail-on-policies-rejected
aictl get scan report sarif "$sid" -o out/sarif.json
```

---

## Gaps P1

| Gap | Источник | Предлагаемое API aictl | Workaround сейчас | Статус |
|-----|----------|------------------------|-------------------|--------|
| Баги `nist` / `oud4` | aictl | Исправить cobra `Use` у существующих команд | Вызывать другие типы отчётов | **готово** |

---

## Gaps P2

| Gap | Источник | Предлагаемое API / решение | Статус |
|-----|----------|----------------------------|--------|
| Обновление sources из Git | stub / SourceControl | Агент сканирования обновляет VCS; команду aictl не добавлять | Не планируется |
| Auth user/password | ptai `--user`/`--password` | — | Не планируется (token) |
| Truststore PEM | ptai `--truststore` | `--cacert <path>` (+ `ctx set` / `ctx unset --cacert`) | Сделано: путь к PEM, append к system roots |
| Retry при занятой ветке | aisa `--retry` / `--retry-time` | Скрипт [`pipeline-with-retry.sh`](../../examples/pipeline-with-retry.sh); в CLI не добавлять, пока AIE не поддержит множественные сканы на разных ветках | Не планируется (скрипт) |
| Raw / GIIF отчёты | ptai | Форматов нет в AIE | Не планируется |
| Priority на `scan branch` / `scan project` | ptai | Задать нельзя при старте; `update project settings --priority` | Не планируется |
| Include при upload | ptai `-i` | Достаточно excludes (`-e` / `--exclude-from`) | Не планируется |
| Пакетная генерация отчётов | ptai `--report-json` | Несколько вызовов `get scan report` | Не планируется |
| `get scan result` | удалено | — | Не для миграции |
| Rich exit codes aisa (2…43) | aisa | Справочник в [aisa-to-aictl.md](aisa-to-aictl.md); схема aictl остаётся **0 / 1 / 2** | Не планируется |
| One-shot meta-команда | aisa / `ui-ast` | — | Не планируется |
| Скрытые/мёртвые флаги aisa | `--list-results`, `--restore-sources`, … | — | Не планируется |

---

## Exit codes

| Инструмент | Схема |
|------------|-------|
| aictl (сейчас и план) | **0** успех, **1** validation / бизнес-gate (в т.ч. fail-флаги), **2** API/сеть, **-1** неизвестная ошибка |
| aisa | Детальные коды (10 — policy, 4 — project not found, 29 — token, …) |
| ptai-cli-plugin | **0** / **1** / **1000** (невалидный ввод) |

Коды aisa **не** копируются 1:1. Соответствие «бывший код aisa → сообщение / exit aictl» — в [aisa-to-aictl.md](aisa-to-aictl.md).

---

## Вне scope

- Порт Jenkins `ptaiAst` и TeamCity `ptsecurity`
- One-shot `aictl run` как клон aisa
- Паритет скрытых и неиспользуемых флагов aisa
- Совместимость с устаревшими PDF/XML export ptai до 4.0

---

## Порядок доработок

1. ~~Fail-флаги (`apperror`); `scan check-policies`; get/set policies~~ **готово**
2. ~~get/set exclusions (gitignore-текст); `get scan errors`; `get queue`; `get scanning`~~ **готово**
3. ~~`update project languages`; `--temp-dir`; `get report-templates` (id+имя, regex, `-q`)~~ **готово**
4. ~~Багфикс nist/oud4~~ **готово**
5. ~~E2E под новые команды~~ **готово**

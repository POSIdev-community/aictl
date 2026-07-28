# Backlog доработок aictl

Список фич для реализации. Без сравнения с другими CLI — только целевое API aictl.

Связанные документы: [gap-analysis](gap-analysis.md), [aisa → aictl](aisa-to-aictl.md), [ptai-cli-plugin → aictl](ptai-cli-plugin-to-aictl.md).

## Соглашения (зафиксировано)

- Exit codes: **0** / **1** / **2**
- Fail-флаги — **opt-in**; при срабатывании — exit **1** через **отдельный `apperror`** (не `validation.Error`)
- `--fail-on-scan-failed`: стадии **Failed** и **Aborted**
- `--fail-on-policies-rejected`: exit **1**, если `PolicyState` = **`Rejected`** (enum: `None` \| `Rejected` \| `Confirmed`)
- `scan check-policies` без флага: в stdout — **сырое** значение `PolicyState`
- Get policies: stdout — **сырой** ответ API (JSON)
- Get exclusions: stdout — **тело** поля `exclusions` (gitignore-like текст), не JSON-обёртка
- Set policies: JSON; ввод: `-f <file>`, `-f -` / аргумент `-` = stdin, позиционный аргумент = текст
- Set exclusions: **gitignore-like текст** (не JSON); тот же контракт поля `exclusions` API; ввод: `-f` / stdin `-` / позиционный аргумент; **полная замена**
- Табличный вывод (где указано «таблица»): только таблица, без JSON-режима

---

## P0

1. ~~**`scan await --fail-on-scan-failed`**~~ **готово**  
   Failed/Aborted → отдельный `apperror`, exit 1.

2. ~~**`get scan stage --fail-on-scan-failed`**~~ **готово**  
   То же для однократной проверки стадии.

3. ~~**`scan check-policies <scan-id> [--fail-on-policies-rejected]`**~~ **готово**  
   Печатает сырой `PolicyState`. С `--fail-on-policies-rejected`: exit 1 при `Rejected`.

4. ~~**`get project policies`**~~ **готово**  
   Сырой ответ API (правила политик проекта).

5. ~~**`set project policies`**~~ **готово**  
   JSON правил политик. Ввод: `-f <file>` | `-f -` / `-` (stdin) | позиционный аргумент (текст JSON). Как `set project settings`.

6. ~~**`get report-templates [<regex>] [-q|--quite] [--localization en|ru]`**~~ **готово**  
   Вывод: **id** и **имя**; фильтр по имени (regex, как `get projects`); `-q` — только id.

7. ~~**`get scan errors <scan-id>`**~~ **готово**  
   Ошибки скана; stdout — **построчный текст**. Project id: `-p` / ctx.

8. ~~**`update project languages`**~~ **готово**  
   Актуализировать языки по **уже загруженным** sources на сервере. Без аргументов; project: ctx или `-p`.

9. ~~**`get project exclusions`**~~ **готово**  
   В stdout — **тело** поля `exclusions` (gitignore-like текст), не JSON-модель API.

10. ~~**`set project exclusions`**~~ **готово**  
    Полная замена exclusions. Тело — **gitignore-like текст** (не JSON). Ввод: `-f <file>` | stdin (`-`) | позиционный аргумент.

11. ~~**`get queue`**~~ **готово**  
    Очередь сканирований (API scan queue). Вывод — **таблица**: project id, branch id, scan id, стадия.
  - без `-p` — вся очередь; **ctx `-p` не учитывать**
  - `-p <project-id>` — фильтр по проекту

12. ~~**`get scanning`**~~ **готово**  
    Список сканов, которые **сейчас сканируются** (API active scans). Вывод — **таблица** с теми же колонками (project id, branch id, scan id, стадия).
  - без `-p` — все; **ctx `-p` не учитывать**
  - `-p <project-id>` — фильтр по проекту

13. ~~**`--temp-dir <path>`** на **`update sources`** и **`create branch`**~~ **готово**  
    Каталог, где создать временный zip при упаковке sources. После успешной загрузки — **удалить** архив.  
    Не путать с `--scan-target` (путь к исходникам на `create branch`).

---

## P1

14. ~~**Багфикс `get scan report nist` / `oud4`**~~ **готово**  
    Cobra `Use` больше не дублирует markdown/json.

---

## Не делать

- One-shot meta-команда (`aictl run` и аналоги)
- Auth user/password
- Truststore PEM (достаточно `--tls-skip` / системный trust)
- Расширение exit codes под legacy-схемы (10, 60, 1000, …)
- `get scan result` — удалено (stub без потребителей)
- `delete projects` по имени / regexp — достаточно `get projects <regex>` → `delete projects <uuid>…`
- `get project --name` — достаточно `get projects <regex>`
- `update sources git` — удалено; обновление sources из VCS выполняет агент сканирования
- Retry при занятой ветке в CLI — скрипт [`pipeline-with-retry.sh`](../../examples/pipeline-with-retry.sh), пока AIE не поддержит множественные сканы на разных ветках
- `get scan report raw` / `giif` — форматов нет в AIE
- `--priority` на `scan start` — нельзя; есть `update project settings --priority`
- `update sources --include` — достаточно `-e` / `--exclude-from`
- `get scan reports -f` — отдельные `get scan report …`

---

## Порядок реализации

1. ~~Fail-флаги + `apperror`; `scan check-policies`; get/set project policies~~ **готово**
2. ~~get/set project exclusions; `get scan errors`; `get queue`; `get scanning`~~ **готово**
3. ~~`update project languages`; `--temp-dir`; `get report-templates`~~ **готово**
4. ~~Багфикс nist/oud4~~ **готово**
5. ~~E2E под новые команды~~ **готово** (leaf вне smoke: `TestLeafCommandsOutsideSmoke`)

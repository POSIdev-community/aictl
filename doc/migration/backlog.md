# Backlog доработок aictl

Список фич для реализации. Без сравнения с другими CLI — только целевое API aictl.

Связанные документы: [gap-analysis](gap-analysis.md), [aisa → aictl](aisa-to-aictl.md), [ptai-cli-plugin → aictl](ptai-cli-plugin-to-aictl.md).

## Соглашения (зафиксировано)

- Exit codes: **0** / **1** / **2**
- Fail-флаги — **opt-in**; при срабатывании — exit **1** через **отдельный `apperror`** (не `validation.Error`)
- `--fail-on-scan-failed`: стадии **Failed** и **Aborted**
- `--fail-on-policy`: exit **1**, если `PolicyState` ∈ {**`Confirmed`**, **`None`**} (enum: `None` \| `Rejected` \| `Confirmed`)
- `get scan policy` без флага: в stdout — **сырое** значение `PolicyState`
- Get policies / exclusions: stdout — **сырой** ответ API
- Set policies: JSON; ввод: `-f <file>`, `-f -` / аргумент `-` = stdin, позиционный аргумент = текст
- Set exclusions: **gitignore-like текст** (не JSON); тот же контракт поля `exclusions` API; ввод: `-f` / stdin `-` / позиционный аргумент; **полная замена**
- Табличный вывод (где указано «таблица»): только таблица, без JSON-режима
- Реализацию кода **не** начинать, пока backlog не согласован отдельно; этот документ — спецификация

---

## P0

1. **`scan await --fail-on-scan-failed`**  
   Failed/Aborted → отдельный `apperror`, exit 1.

2. **`get scan stage --fail-on-scan-failed`**  
   То же для однократной проверки стадии.

3. **`get scan policy <scan-id> [--fail-on-policy]`**  
   Печатает сырой `PolicyState`. С `--fail-on-policy`: exit 1 при `Confirmed` или `None`.

4. **`get project policies`**  
   Сырой ответ API (правила политик проекта).

5. **`set project policies`**  
   JSON правил политик. Ввод: `-f <file>` | `-f -` / `-` (stdin) | позиционный аргумент (текст JSON). Как `set project settings`.

6. **`get report-templates [<regex>] [-q|--quite] [--localization en|ru]`**  
   Вывод: **id** и **имя**; фильтр по имени (regex, как `get projects`); `-q` — только id.

7. **`get scan errors <scan-id>`**  
   Ошибки скана; stdout — **построчный текст**. Project id: `-p` / ctx.

8. **`update project languages`**  
   Актуализировать языки по **уже загруженным** sources на сервере. Без аргументов; project: ctx или `-p`.

9. **`get project exclusions`**  
   Сырой ответ API (модель с полем exclusions).

10. **`set project exclusions`**  
    Полная замена exclusions. Тело — **gitignore-like текст** (не JSON). Ввод: `-f <file>` | stdin (`-`) | позиционный аргумент.

11. **`get queue`**  
    Очередь сканирований (API scan queue). Вывод — **таблица**: project id, branch id, scan id, стадия.
  - без `-p` — вся очередь; **ctx `-p` не учитывать**
  - `-p <project-id>` — фильтр по проекту

12. **`get scanning`**  
    Список сканов, которые **сейчас сканируются** (API active scans). Вывод — **таблица** с теми же колонками (project id, branch id, scan id, стадия).
  - без `-p` — все; **ctx `-p` не учитывать**
  - `-p <project-id>` — фильтр по проекту

13. **`--temp-dir <path>`** на **`update sources`** и **`create branch`**  
    Каталог, где создать временный zip при упаковке sources. После успешной загрузки — **удалить** архив.  
    Не путать с `--scan-target` (путь к исходникам на `create branch`).

---

## P1

14. **Багфикс `get scan report nist` / `oud4`**  
    Исправить cobra `Use` (сейчас дублируют markdown/json).

---

## Не делать

- One-shot meta-команда (`aictl run` и аналоги)
- Auth user/password
- Truststore PEM (достаточно `--tls-skip` / системный trust)
- Расширение exit codes под legacy-схемы (10, 60, 1000, …)
- `get scan result` (stub без потребителей)
- `delete projects` по имени / regexp — достаточно `get projects <regex>` → `delete projects <uuid>…`
- `get project --name` — достаточно `get projects <regex>`
- `update sources git` — обновление sources из VCS выполняет агент сканирования
- Retry при занятой ветке в CLI — скрипт [`pipeline-with-retry.sh`](../../examples/pipeline-with-retry.sh), пока AIE не поддержит множественные сканы на разных ветках
- `get scan report raw` / `giif` — форматов нет в AIE
- `--priority` на `scan start` — нельзя; есть `update project settings --priority`
- `update sources --include` — достаточно `-e` / `--exclude-from`
- `get scan reports -f` — отдельные `get scan report …`

---

## Порядок реализации (когда начнётся код)

1. Fail-флаги + `apperror`; `get scan policy`; get/set project policies
2. get/set project exclusions; `get scan errors`; `get queue`; `get scanning`
3. `update project languages`; `--temp-dir`; `get report-templates`
4. Багфикс nist/oud4
5. E2E под новые команды

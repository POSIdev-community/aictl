# ADR-0001: Изолированные локальные contexts

- **Задача:**
- **Спецификация:**
- **Статус:** <span style="color:#757575;">ЧЕРНОВИК</span>
- **Дата:** **2026-08-26**

## Содержание

1. [Цель](#цель)
2. [Проблема](#проблема)
3. [Требования/Ограничения](#требованияограничения)
4. [Решение](#решение)
5. [Варианты/Возможности](#вариантывозможности)
    - [Отдельный context-файл](#отдельный-context-файл)
    - [Environment variables для подключения](#environment-variables-для-подключения)
    - [Только per-command flags](#только-per-command-flags)
    - [Общий current-context как механизм изоляции CI](#общий-current-context-как-механизм-изоляции-ci)
    - [Lock глобального context.yaml](#lock-глобального-contextyaml)
6. [Вывод](#вывод)
7. [Вопросы](#вопросы)

---

## Цель

- Изолировать настройки независимых job/скриптов/терминалов.
- Сохранить существующий UX без дополнительных флагов для последовательной работы.
- Дать интерактивному пользователю именованные contexts.
- Дать CI-job прямой способ управлять жизненным циклом контекста.
- Сделать операции записи устойчивыми к одновременным процессам и аварийному завершению.

---

## Проблема

Сейчас aictl хранит единственную конфигурацию пользователя в:

```text
~/.config/aictl/context.yaml
```

При заданном `XDG_CONFIG_HOME` используется `$XDG_CONFIG_HOME/aictl/context.yaml`. Файл содержит:

```yaml
uri: https://example.com
token: token
tls-skip: false
projectId: 00000000-0000-0000-0000-000000000000
branchId: 00000000-0000-0000-0000-000000000000
```

`aictl ctx set`, `show`, `unset` и `clear` работают только с этим файлом.
connection-команды затем накладывают на него `-u`, `-t`, `--tls-skip`, `-p` и `-b` для текущего процесса.

Такая модель создаёт гонку между независимыми процессами. 
Две CI-job/два скрипта могут по очереди выполнить `ctx set -p ... -b ...`; 
последний writer меняет scope всех последующих команд, использующих context.

---

## Требования/Ограничения

### Опыт kubectl

Kubernetes разделяет выбор файла и выбор logical context:

- default — `$HOME/.kube/config`;
- `--kubeconfig` выбирает ровно один файл;
- `KUBECONFIG` задаёт список файлов для merge;
- `--context` имеет приоритет над `current-context` из конфигурации;
- context группирует cluster, user и namespace под именем.

Это описано в официальной документации [Organizing Cluster Access Using kubeconfig Files](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/) 
и [kubectl config](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_config/).

[`kubectl config use-context`](https://kubernetes.io/docs/reference/kubectl/generated/kubectl_config/kubectl_config_use-context/) 
записывает `current-context` обратно в kubeconfig. 
Следовательно, общий `current-context` сам по себе не решает исходную проблему для параллельных job: 
общий selector всё равно изменяется несколькими процессами. 
kubectl закрывает это отдельным файлом, выбранным через `KUBECONFIG`; aictl отказывается от второго механизма выбора и закрывает тот же случай 
изолированным конфиг-каталогом job (`XDG_CONFIG_HOME`) поверх единственного механизма — named contexts.

Из kubectl принимаются следующие принципы:

1. Существующий default-файл продолжает работать без настройки.
2. Явный CLI selector имеет наивысший приоритет.
3. Environment selector позволяет один раз настроить всю job или shell.
4. context имеет имя.
5. Конфликтующие явные selectors должны завершаться ошибкой.

### Единственный механизм выбора

Выбор context выполняется только по имени.
Произвольный файл участвует только как 
**источник данных при создании** named context (`ctx create --from-file`) и никогда не выбирается для записи.

### Обратная совместимость
1) Нет новых flags/env и не выполнялся `ctx use` - используется прежний `context.yaml`
2) `ctx set/show/unset` без selector до `ctx use` - работают как раньше
3) Синтаксис `ctx clear` сохраняется, но storage-семантика меняется:
раньше файл удалялся, теперь context остаётся выбранным, а все его поля очищаются
4) `ctx use default` - удаляет сохранённый selector; при отсутствии CLI/env selector команды используют файл `context.yaml`


Named contexts и default используют одну и ту же YAML schema, поэтому context можно при необходимости скопировать вручную.

### Риски

1. **Token в аргументах процесса.** `ctx set -t "$TOKEN"` делает token видимым в `ps`, `/proc/<pid>/cmdline` и в логах при `set -x`.
Для CI рекомендуется передавать весь context через `ctx create --from-file -` (stdin).
2. **Token на диске.** Named context хранит token в `contexts/<name>.yaml`.
Файл переживает завершение job, если каталог не изолирован; 
основную очистку должны обеспечивать per-job `XDG_CONFIG_HOME` и `ctx delete` в trap.
3. **Гонка общего selector.** `ctx use` и `ctx create --use` пишут один общий файл `current-context`.
Параллельные процессы с общим конфиг-каталогом перезапишут его друг у друга;
изоляция обеспечивается `AICTL_CONTEXT` или отдельным `XDG_CONFIG_HOME`.
4. **Остаточные секреты.** Job может завершиться аварийно до `ctx delete`; страховкой служит временный конфиг-каталог runner’а.
5. **Совместимость permissions.** Изменение legacy mode с `0644` на `0600`
может затронуть нестандартные сценарии совместного чтения несколькими OS users.

## Решение

Поддержать два источника context:

1. **Default context** — существующий `~/.config/aictl/context.yaml`.
2. **Named context** — отдельный файл `~/.config/aictl/contexts/<name>.yaml`.

### Selector

Добавить один root-флаг:

```text
--context <name>
```

И одну переменную окружения:

```text
AICTL_CONTEXT - выбирает именованный context из доступных.
Для терминала, постоянных окружений и CI: dev, stage, prod, ci-<job id>

export AICTL_CONTEXT=dev
```

Каждая команда работает ровно с одним context-файлом.

Флаг `--context` и команда `ctx` — разные вещи, несмотря на созвучные имена:
`--context <name>` **выбирает**, над каким context выполняется вызов, и действует только на текущий процесс;
`ctx` **управляет содержимым** выбранного context (`set`, `show`, `unset`, `clear`) и его жизненным циклом
(`create`, `list`, `current`, `delete`, `use`). Команды `context` нет — селектор существует только в форме флага,
переменной окружения и сохранённого выбора `ctx use`.

Флаг персистентный, поэтому Cobra принимает его в любой позиции: `aictl ctx set --context dev -u ...`
и `aictl --context dev ctx set -u ...` эквивалентны. В документации и примерах selector ставится **после** имени команды:
в позиции перед командой слово `--context` читается как имя команды `ctx` и сбивает читателя.

### Параметры подключения

Отдельных environment variables для подключения нет. Значения берутся из выбранного context-файла 
и могут быть переопределены только флагами текущего вызова.

Приоритет параметров подключения, от самого высокого к самому низкому:

1. CLI-флаги `-u`, `-t`, `--tls-skip`, `-p` и `-b`.
2. Значения из выбранного context-файла.

Флаги действуют один вызов и не записываются обратно в context-файл.
Единственная команда, которая пишет значения на диск, — `ctx set`.

Следствие: `ctx show` показывает ровно те значения, которые будут использованы командой

### Приоритет выбора context

context определяется один раз на процесс:

1. CLI `--context <name>`.
2. `AICTL_CONTEXT`.
3. Имя, сохраненное командой `ctx use` или `ctx create --use` в `$XDG_CONFIG_HOME/aictl/current-context`.
4. default `~/.config/aictl/context.yaml`.

Правила конфликтов:

- если CLI selector задан, environment selector и сохранённый selector(файл) -  игнорируются;
- если CLI selector отсутствует, `AICTL_CONTEXT` перекрывает сохраненный selector;
- пустая environment variable считается незаданной;
- имя валидируется до любого обращения к файловой системе.

#### Зарезервированные имена

Имена `default` и `current` зарезервированы и недопустимы в качестве имен named contexts.

- `default` — явный alias изначального `context.yaml`; `--context default`, `AICTL_CONTEXT=default`, `ctx use default` и `--from default` выбирают именно его. 
Файл `contexts/default.yaml` не создаётся и не читается никогда.
- `current` — обозначение context, выбранного для текущего процесса по правилам приоритета; собственного файла у него нет.


#### Формат имен
Имена named contexts ограничиваются безопасным набором `[A-Za-z0-9][A-Za-z0-9._-]{0,63}`; path separators и `..` запрещены. 
Правило применяется во всех командах и selectors, принимающих имя: `--context`, `AICTL_CONTEXT`, `ctx create`, `ctx use`, `ctx delete`. 
Путь, собранный из имени, после нормализации обязан лежать внутри `contexts/`; иначе команда завершается validation error.


### Поведение ctx

Перед выполнением connection-команды aictl:

1. Определяет текущий context по правилам в пункте `Приоритет выбора context`.
2. Читает настройки из выбранного context-файла: `context.yaml` или named `contexts/<name>.yaml`.
3. Накладывает переданные для одного вызова флаги `-u`, `-t`, `--tls-skip`, `-p` и `-b` поверх полученной конфигурации.
4. Валидирует необходимые поля.
5. Выполняет запрос к <СИСТЕМЕ>; сам context при этом не меняется.

Существующие команды управления меняют только выбранный context:

1) `ctx set` - Создаёт/изменяет текущий context-файл
2) `ctx show`- Показывает текущий context
3) `ctx unset` - Сбрасывает поля только текущего context
4) `ctx clear` - Очищает все поля текущего context, но сохраняет сам файл и его выбор

Добавить команды:

```text
aictl ctx create <name> [--from <name|default|current>] [--from-file <path|->] [--use]
aictl ctx list [--json]
aictl ctx current [--json]
aictl ctx delete <name>... [-y] [--ignore-missing]
aictl ctx use <name|default>
```

#### Создание

Создание пустого named context:

```bash
aictl ctx create dev
aictl ctx set --context dev -u https://dev.example -t "$TOKEN"
```
Создание из существующего context выполняет копирование по шаблону:

```bash
aictl ctx create stage --from default
aictl ctx create feature-123 --from dev
```

После создания source и target полностью независимы.
Последующие `ctx set` в source не меняют target.

По умолчанию `--from` копирует текущие значения всех пяти полей source context:

```yaml
uri: https://example.com
token: token
tls-skip: false
projectId: 123e4567-e89b-12d3-a456-426614174011
branchId: 123e4567-e89b-12d3-a456-426614174000
```

Для частого сценария:

```bash
aictl ctx create feature-123 --from dev
```

`--from current` копирует context, выбранный для текущего процесса:

```bash
AICTL_CONTEXT=dev aictl ctx create debug-copy --from current
```

Создание завершается ошибкой, если target уже существует.
Для намеренной замены используются `ctx delete` и повторный `ctx create`.

Создание также завершается ошибкой, если target-имя — `default` или `current`: оба зарезервированы (см. [Зарезервированные имена](#зарезервированные-имена)).
Имя проверяется до чтения source, поэтому `ctx create current --from dev` не создаёт и не изменяет ни одного файла.

`ctx set` сохраняет create-on-write поведение: если выбранный named context ещё не существует, он создаётся. 
`ctx create` нужен для явного управления жизненным циклом контекста.

##### Создание из файла

`--from-file` создаёт named context из произвольного YAML-файла:

```bash
aictl ctx create ci-667 --from-file ./ci-context.yaml
```

Файл только читается и никогда не становится целью записи: aictl пишет исключительно в `contexts/<name>.yaml`. 
Поэтому произвольный путь здесь не создаёт рисков подмены или очистки чужого файла.

`--from-file` выполняет жёсткую валидацию источника:

1. Документ должен быть одиночным валидным YAML mapping; пустой файл, multi-document YAML и не-YAML содержимое — ошибка.
2. Разрешён ровно набор ключей `uri`, `token`, `tls-skip`, `projectId`, `branchId`. Любой неизвестный ключ — ошибка с указанием его имени.
3. Отсутствующий ключ допустим: поле остаётся `<unset>`.
4. Типы проверяются строго: `tls-skip` — bool, остальные — строки.
5. `uri` — абсолютный URL со схемой `http` или `https`.
6. `projectId` и `branchId` — UUID.
7. Размер входа ограничен; превышение лимита — ошибка.

`--from` использует тот же валидатор: разница между источниками только в происхождении данных, но не в строгости проверки.

`--from` и `--from-file` взаимоисключающие.

Значение `-` читает описание context из stdin. Это рекомендуемый способ для CI: 
token не попадает ни в аргументы процесса, ни в файл репозитория.

```bash
printf 'uri: %s\ntoken: %s\ntls-skip: false\n' "$CI_AICTL_URI" "$CI_AICTL_TOKEN" \
  | aictl ctx create ci-667 --from-file -
```

##### Немедленное переключение

`--use` переключает сохранённый selector на созданный context сразу после успешного создания:

```bash
aictl ctx create dev --from default --use
```

Флаг эквивалентен последующему `ctx use <name>` и подчиняется тем же правилам:

- selector записывается в общий файл `$XDG_CONFIG_HOME/aictl/current-context`;
- переключение выполняется только после успешного создания context; при ошибке создания selector не меняется;
- CLI `--context` и `AICTL_CONTEXT` остаются приоритетнее сохранённого selector, 
поэтому при заданном `AICTL_CONTEXT` эффект `--use` на последующие команды не виден;
- `current-context` — общий mutable файл, 
поэтому параллельные процессы с одним конфиг-каталогом не должны рассчитывать на `--use` как на средство изоляции (см. [CI-сценарий](#ci-сценарий)).

#### Просмотр

`ctx list` показывает `default` и named contexts из `contexts/`.

```bash
aictl ctx list

CURRENT  NAME                  KIND      CONNECTION  PROJECT  BRANCH
         default               default   incomplete  unset    unset
*        dev2                  named     ready       set      set
```
Token и его фрагменты не выводятся.
Поле `CURRENT` означает selection текущего процесса.

`ctx current` печатает, какой context выбран, как именно он выбран и где лежит его файл.

```bash
aictl ctx current

Name:   default
Kind:   default
Source: default
Exists: true
Path:   /home/user/.config/aictl/context.yaml
```

```bash
AICTL_CONTEXT=ci-667 
aictl ctx current --json

{
    "name": "ci-667",
    "kind": "named",
    "source": "env",
    "exists": true,
    "path": "/home/user/.config/aictl/contexts/ci-667.yaml"
}
```

Поле `source` принимает значения `flag`, `env`, `use` и `default` и показывает, какое правило приоритета сработало. 
Это основной инструмент диагностики окружения: без него неочевидно, почему `ctx use` или `--use` не влияют на команду.

При `exists: false` путь всё равно печатается — иначе опечатку в имени не отладить. 
Значения полей и token команда не выводит никогда; для содержимого есть `ctx show`.

```bash
aictl ctx show --context dev

uri: https://example.com
token: token_check
tls-skip: false
projectId: 123e4567-e89b-12d3-a456-426614174011
branchId: 123e4567-e89b-12d3-a456-426614174000


aictl ctx show --context dev --json

{
    "uri": "https://dev.example.com",
    "token": "token_check",
    "tlsSkip": false,
    "projectId": "123e4567-e89b-12d3-a456-426614174011",
    "branchId": "123e4567-e89b-12d3-a456-426614174000"
}
```

#### Редактирование

```bash
aictl ctx set --context dev -p "$PROJECT_ID" -b "$BRANCH_ID"
aictl ctx unset --context dev -p -b
```

#### Очистка и удаление

`ctx clear` и `ctx delete` разделены по уровню воздействия: `clear` работает с полями, `delete` - с файлом.

`ctx clear` сохраняет прежнее назначение — очистить выбранный context, но меняет способ хранения: файл больше не удаляется.

```bash
aictl ctx clear --context dev -y
```

`ctx clear` сбрасывает uri, token, tls-skip, projectID и branchID, но не удаляет context-файл. 
Выбранный через `ctx use` context остается текущим, а сохранённый selector не изменяется.

Если context выбран через `--context` или `AICTL_CONTEXT`, очищается только указанный context. 
CLI-флаг действует один вызов, а environment variable продолжает выбирать тот же, теперь пустой context.

`ctx clear` применим и к `default`: это единственный способ очистить `context.yaml`.

`ctx delete` **удаляет сам файл** named context и может принимать несколько имён:

```bash
aictl ctx delete dev feature-123
aictl ctx delete dev -y
```

Правила:

1. Удаляется только `contexts/<name>.yaml`. Файл `context.yaml` не удаляется никогда; `ctx delete default` 
запрещён — для default существует `ctx clear`. `ctx delete current` запрещён по той же причине: `current` — зарезервированное имя, 
а не файл; чтобы удалить выбранный context, его имя указывается явно.
2. Каталог `contexts/` не удаляется, даже если стал пустым.
3. Имена валидируются до обращения к файловой системе; при невалидном имени команда не выполняет ни одного удаления.
4. Если `contexts/<name>.yaml` оказался symbolic link, aictl не идёт по ссылке: удаляется сама ссылка.
5. Без `-y` требуется подтверждение.
6. По умолчанию удаление отсутствующего context — ошибка. `--ignore-missing` делает операцию идемпотентной; 
это необходимо для `trap ... EXIT` в CI, иначе cleanup подменяет код возврата job.
7. Удаление context, выбранного сохранённым selector, разрешено: 
selector сбрасывается, и без CLI/env selector команды снова используют `default`. 
Запрет ломал бы CI-job, удаляющую собственный context.
8. При нескольких именах ошибка файловой системы на одном имени не прерывает обработку остальных;
команда завершается ненулевым кодом и списком неудач.

Удаление выполняется обычным `unlink` и не затирает содержимое: token может остаться в освободившихся блоках файловой системы или в снапшотах.

#### Переключение

Для последовательной интерактивной работы добавляется команда:

```bash
aictl ctx use dev
aictl ctx use default
```

Она принимает только существующий named context либо `default`; `ctx use current` — validation error, 
поскольку `current` не имя context, а обозначение уже сделанного выбора. 
`ctx use default` удаляет только сохранённый selector `$XDG_CONFIG_HOME/aictl/current-context` и возвращает выбор к `context.yaml`,
если он не переопределён CLI-флагом или environment variable.
Команда не изменяет `AICTL_CONTEXT` в родительском shell; если переменная задана, environment selector продолжает иметь приоритет.
Для разовой команды используется `--context dev`.

`ctx use` и `ctx create --use` намеренно не считаются гарантией параллельной изоляции:
это общий mutable selector для удобного однопользовательского сценария. 
Параллельные процессы, делящие конфиг-каталог, не должны на него опираться.

### Default и сохранённый выбор

На чистой установке и после `ctx use default`, если отсутствуют CLI/env selectors, поведение полностью прежнее:

```bash
aictl ctx set -u https://example -t "$TOKEN"
aictl get projects
```

Обе команды используют `~/.config/aictl/context.yaml`. 
После `ctx use dev` команды без явного selector используют named context `dev`. 
Автоматической миграции или переименования legacy-файла нет.

### CI-сценарий

Каждая независимая CI-job создаёт собственный named context в начале работы и удаляет его в конце. 
Уникальность имени обеспечивает идентификатор job. URI и token задаются в настройках CI как masked/protected secret variables 
и передаются в context через stdin, чтобы не появляться в аргументах процесса. ProjectID и branchID принадлежат только этой job:

```bash
set -euo pipefail

# собственный конфиг-каталог job: current-context и contexts/ становятся job-local
export XDG_CONFIG_HOME="${RUNNER_TEMP}/aictl-${CI_JOB_ID}"
CTX="ci-${CI_JOB_ID}"

cleanup() {
    aictl ctx delete "$CTX" -y --ignore-missing
}
trap cleanup EXIT

# CI_AICTL_URI и CI_AICTL_TOKEN заданы в настройках CI как masked variables
printf 'uri: %s\ntoken: %s\ntls-skip: false\n' "$CI_AICTL_URI" "$CI_AICTL_TOKEN" \
  | aictl ctx create "$CTX" --from-file - --use

aictl ctx set -p "$PROJECT_ID" -b "$BRANCH_ID"

aictl get scans
aictl scan start
```

```text
aictl ctx current

Name:   ci-667
Kind:   named
Source: use
Exists: true
Path:   /builds/tmp/aictl-667/aictl/contexts/ci-667.yaml


aictl ctx list

CURRENT  NAME              KIND      CONNECTION  PROJECT  BRANCH
         default           default   incomplete  unset    unset
*        ci-667            named     ready       set      set
```

Изоляцию здесь обеспечивает **не** `--use`, а собственный `XDG_CONFIG_HOME` job: 
файл `current-context` в нём принадлежит только этой job, поэтому переключение не видно соседним job на том же runner. 
Каталог удаляется runner’ом вместе с `RUNNER_TEMP`, а `ctx delete` в trap отрабатывает штатное завершение.

Если изолировать конфиг-каталог невозможно (общий `$HOME`, shell executor), `--use` использовать нельзя — 
общий `current-context` будет перезаписан соседней job. В этом случае выбор задаётся переменной окружения, наследуемой всеми процессами job:

```bash
export AICTL_CONTEXT="ci-${CI_JOB_ID}"
trap 'aictl ctx delete "$AICTL_CONTEXT" -y --ignore-missing' EXIT

printf 'uri: %s\ntoken: %s\ntls-skip: false\n' "$CI_AICTL_URI" "$CI_AICTL_TOKEN" \
  | aictl ctx create "$AICTL_CONTEXT" --from-file -

aictl ctx set -p "$PROJECT_ID" -b "$BRANCH_ID"
aictl get scans
```

Разные job получают разные имена и не конкурируют ни при чтении, ни при записи данных context. 
Единственное разделяемое состояние в этом варианте — `current-context`, и его не трогает ни одна команда job.

Передача uri и token через `--from-file -` предпочтительнее флагов `-u`/`-t`: 
значение не попадает в `ps`, `/proc/<pid>/cmdline` и в логи при `set -x`.

Named context остаётся тем же механизмом и для постоянных локальных окружений:

```bash
aictl ctx create dev
aictl ctx set --context dev -u https://dev.example -t "$DEV_TOKEN"

aictl ctx create prod
aictl ctx set --context prod -u https://prod.example -t "$PROD_TOKEN"

aictl ctx use dev
```

### Конкурентность и сохранность файла

Разные файлы контекста являются основной границей изоляции.
Каждая независимая CI-job или локальная session, которой требуется собственный scope, должна использовать отдельный named context.
Несколько процессов `aictl` внутри одной job/session используют один и тот же выделенный ей context.

Для защиты от частично записанного YAML и потери успешно завершённой записи используется полный crash-safe цикл:

1. Создать родительский каталог с правами `0700`, если он отсутствует.
2. Создать уникальный временный файл в том же каталоге, что и target. Временный файл создаётся с правами `0600`.
3. Полностью сериализовать и записать YAML, проверив ошибки `Write`.
4. Выполнить `fsync` временного файла, чтобы его содержимое было передано из буферов ОС в хранилище.
5. Закрыть временный файл и проверить ошибку `Close`.
6. Атомарно заменить target через `rename`.
7. Открыть родительский каталог и выполнить его `fsync`, чтобы зафиксировать изменение directory entry после `rename`.
8. При любой ошибке удалить временный файл; успешно установленный target не удалять.

Тот же цикл применяется к файлу `current-context`: он такое же persistent состояние, и его повреждение оставляет CLI без валидного selector.

Чтение выполняется без блокировки: atomic replacement возвращает либо старую, либо новую целую версию файла.
На платформах, где синхронизация каталога недоступна, adapter должен использовать документированный platform-specific fallback.

### Изменения в архитектуре

Сейчас `application.NewApplication()` читает фиксированный default-файл до разбора Cobra arguments.
Для root selector это слишком рано.

Новый флоу должен выглядеть так:

```text
1) Context CLI flag + selector environment
        |
2) ContextSelector.Resolve
        |
3) ContextStore(selected path)
        |
4) domain Config
        |
5) connection flag overlay
        |
6) use case
----------------------------------------------------------------
1) Собирает источники выбора context: --context, AICTL_CONTEXT,
сохранённый ctx use / ctx create --use или default.
2) Применяет правила приоритета, валидирует имя и проверяет конфликты.
3) Определяет путь к YAML-файлу и загружает его. Например: ~/.config/aictl/contexts/dev.yaml.
4) Преобразует YAML в Config в памяти процесса.
5) Накладывает -u, -t, --tls-skip, -p и -b. Значения действуют только для текущего вызова и не записываются обратно в context-файл.
6) Валидирует эффективную конфигурацию и выполняет бизнес-сценарий команды: get, scan, create, update и т. д.
```

## Варианты/Возможности

### Отдельный context-файл

Второй selector — `--context-file` / `AICTL_CONTEXT_FILE`, выбирающий произвольный путь. 
Технически достаточен для CI, но удваивает матрицу приоритетов и правил конфликта, 
требует валидации произвольных путей на запись (symlink, `..`, чужой файл) и не даёт ничего, 
чего не даёт named context в изолированном `XDG_CONFIG_HOME`. Отклонён в пользу единственного механизма выбора по имени. 
Произвольный файл сохранён только как источник данных для `ctx create --from-file`, то есть в режиме чтения.

### Environment variables для подключения

Отдельные `AICTL_URI` и `AICTL_TOKEN`, накладываемые поверх context. 
Дают три источника значений (файл → environment → флаги), из-за чего `ctx show` перестаёт отвечать на вопрос «что реально уйдёт в запрос», 
а переопределение становится невидимым в выводе диагностических команд. Отклонены: один источник значений — context-файл, разовое переопределение — только флаги.

### Только per-command flags

Уже возможно, но заставляет повторять URI/token/project/branch в каждой команде 
и не сохраняет удобство context. При этом флаги остаются легитимным способом разового переопределения, 
а для передачи секрета без появления в аргументах процесса предусмотрен `ctx create --from-file -`.

### Общий current-context как механизм изоляции CI

`ctx create <name> --use` в начале каждой job выглядит как замена per-job файлу, но не является ею: 
`--use` пишет общий `current-context`, и при общем конфиг-каталоге последняя стартовавшая job переключает scope остальных — 
ровно исходная проблема этого ADR. Вариант принят только в связке с job-local `XDG_CONFIG_HOME`; 
в остальных случаях выбор задаётся `AICTL_CONTEXT`.

### Lock глобального context.yaml

Защитит файл от повреждения, но не решит логическую гонку: 
последняя job всё равно поменяет scope для остальных.

## Вывод

1. Источников context два: default `context.yaml` и named `contexts/<name>.yaml`.
2. Механизм выбора один — по имени: флаг `--context`, переменная `AICTL_CONTEXT`, сохранённый selector `ctx use` / `ctx create --use`, затем default.
3. Имена `default` и `current` зарезервированы как keywords селектора и недопустимы в качестве имён named contexts: файлы `contexts/default.yaml` и `contexts/current.yaml` не создаются и не читаются.
4. Параметры подключения хранятся только в context-файле; отдельных environment variables нет, разовое переопределение — только флаги команды.
5. Жизненным циклом управляют `ctx create` (в том числе `--from`, `--from-file`, `--use`), `ctx clear` (очистка полей) и `ctx delete` (удаление файла named context; `default` удалить нельзя).
6. Диагностика окружения — `ctx current`, который печатает имя, вид, источник выбора, факт существования и путь.
7. Изоляция параллельных job обеспечивается уникальным именем context, а при общем `current-context` — job-local `XDG_CONFIG_HOME` или `AICTL_CONTEXT`; `ctx use` и `--use` изоляции не гарантируют.
8. Все записи на диск (context-файлы и `current-context`) выполняются crash-safe циклом с правами `0600`.

## Вопросы

1. Нужно ли выдавать warning или error, если существующий context-файл имеет права шире `0600`? Предпочтительна ошибка: в named context гарантированно лежит token.
2. Нужна ли автоматическая очистка неиспользуемых named contexts по TTL для случая общего `$HOME` в CI, или lifecycle всегда должен оставаться явным?
3. Какой platform-specific fallback использовать для crash-safe записи на системах, где `fsync` родительского каталога не поддерживается?
4. Достаточно ли `ctx create --from-file -` для безопасной передачи секрета, или флагам `-t`/`-u` нужна форма `--token-file` / `--token-stdin`?
5. Сохранять ли create-on-write в `ctx set` теперь, когда есть явный `ctx create`, или требовать предварительного создания named context?

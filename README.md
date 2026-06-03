# HarnessBeaver (`bvr`)

Launcher para abrir o **Claude Code** (ou PowerShell/cmd) em vários diretórios de
uma vez — em **abas** do Windows Terminal ou em **janelas separadas**. É um sistema
de atalhos dinâmico para diretórios: você registra projetos, agrupa em pacotes e
abre tudo com um comando.

Tem duas frentes sobre o mesmo núcleo:

- **TUI rica** (`bvr` sem argumentos) — menu interativo navegável.
- **CLI clássica** (`bvr <comando>`) — scriptável, para automação.

## Requisitos

- **Para rodar:** [Windows Terminal](https://aka.ms/terminal) (`wt`) e, para o shell
  `claude`, o Claude Code CLI no PATH. **Nenhum runtime** — o `bvr.exe` é estático.
- **Para compilar:** [Go](https://go.dev/dl/) 1.22+ (`winget install GoLang.Go`).

## Build

```powershell
go mod tidy
.\build.ps1            # gera bvr.exe
.\build.ps1 -All       # cross-compila para win/linux/mac em ./dist
```

Distribua o `bvr.exe`: o destinatário só coloca em uma pasta do `PATH` e roda.

## Uso — TUI

```powershell
bvr        # ou: bvr ui
```

Menu: Abrir projetos · Abrir pacote · Adicionar projeto · Remover projeto ·
Escanear & importar · Criar pacote · Remover pacote · Executar comando ·
Analisar aprendizados · Configurações.

Em qualquer lista, tecle **`/`** para **filtrar** por nome/descrição (útil em listas
longas); `enter` aplica o filtro, `esc` limpa. Navegação: `↑/↓` item, `←/→` (ou PgUp/PgDn)
páginas, `Home/End` início/fim.
Teclas: `↑/↓` navega, `espaço` marca (multi-seleção), `enter` confirma, `esc` volta,
`q`/`ctrl+c` sai.

## Uso — CLI

```powershell
bvr list [projects|packages]            # lista o registro
bvr add "Meu App" F:\proj\meu-app --shell claude
bvr edit meu-app --name "Outro" --path F:\novo --shell pwsh
bvr remove meu-app
bvr scan --root F:\02-company-os\produtos --import   # importa projetos de código
bvr pkg add "SMB Ativo" --projects os247,beauty --mode tabs --shell claude
bvr pkg remove smb-ativo
bvr open os247 beauty --mode windows --shell pwsh    # abre projetos/pacotes
bvr open smb-ativo --dry-run                         # mostra os comandos wt
bvr status                                           # status git de todos os projetos
bvr status smb-ativo                                 # status git de um pacote/ids
bvr each smb-ativo -- git pull                       # roda um comando em cada projeto
bvr each repair beauty --parallel -- git fetch       # idem, em paralelo
```

`bvr each <ids...> -- <comando>` executa o comando no diretório de cada projeto alvo
(pacote ou ids), registrando cada execução no diário. Sequencial mostra a saída ao vivo;
`--parallel` captura e imprime por projeto. Sai com código ≠ 0 se algum projeto falhar.

### Layouts (abas/painéis por projeto)

Um **layout** abre um projeto numa janela com várias abas e/ou painéis (split-panes),
cada um com seu shell e comando — ex.: um painel com `claude`, outro rodando o dev server,
outra aba com `git`.

```powershell
bvr layout add dev --preset dev          # cria layout (presets: claude|dev|triple)
bvr layout tab dev --title git           # adiciona uma aba
bvr layout pane dev --tab 2 --shell pwsh --command "git status" --split V
bvr layout show dev --project repair     # mostra o comando wt resultante
bvr layout assign dev repair             # define o layout padrão do projeto
bvr open repair --layout dev             # abre com o layout (ou só 'bvr open repair' se for o padrão)
bvr open repair --no-layout              # ignora o layout padrão
bvr layout list                          # lista; bvr layout rm <id> remove
```

Na TUI, **Abrir com layout** escolhe o layout e o projeto. Cada painel usa o diretório do
projeto por padrão; `--split H` (horizontal) ou `V` (vertical) define como o painel divide
o anterior. Os painéis dividem sequencialmente o último painel criado.

`bvr status` mostra, por projeto, a branch, se há alterações pendentes (`✗`/`✓`) e
ahead/behind do upstream (`↑`/`↓`). Na **TUI**, esse status aparece (carregado em
segundo plano) na descrição de cada projeto ao Abrir/Remover/montar pacote.

## Diário de bordo e aprendizados

O `bvr` também executa comandos em seu nome, **registra** tudo num diário diário e, a
cada novo dia, **oferece analisar os aprendizados** do dia anterior usando IA.

```powershell
bvr run echo ola                       # executa no shell default e registra
bvr run --shell cmd "echo oi & dir"    # cmd; também: --shell pwsh|bash|zsh
bvr shell                              # REPL: digita comandos seguidos, 'exit' p/ sair
bvr review                             # analisa o dia pendente (ou mostra o cache de hoje)
bvr review 2026-06-02                  # se já houver análise, mostra o cache
bvr review 2026-06-02 --force          # refaz a análise mesmo com cache
bvr review --dry-run                   # mostra o prompt sem chamar a IA (sem custo)

bvr history                            # comandos recentes (mais novos primeiro)
bvr history --project repair --failed  # filtra por projeto e só os que falharam
bvr history --grep "git" --days 30     # filtra por texto, últimos 30 dias

bvr alias set deploy "npm run build && npm run deploy"   # cria um atalho
bvr alias                              # lista os atalhos
bvr do deploy                          # executa o atalho
bvr alias set commit "git commit -m {1}"
bvr do commit "fix: ajuste"            # {1}..{N} e {*} recebem os argumentos
```

Cada comando executado pelo `bvr` é **associado ao projeto** (quando o diretório cai
dentro de um projeto registrado), permitindo `bvr history --project <id>` e contexto de
projeto na análise. Na **TUI**, há as telas **Histórico** (re-executa o comando escolhido)
e **Atalhos** (roda um alias) — ambas com o filtro `/`.

- **Cache:** uma vez analisado, o dia fica salvo. `bvr review <dia>` mostra o cache;
  use `--force` para reanalisar. Na TUI, escolher um dia já analisado oferece **ver o
  cache** ou **reanalisar**. O **dia atual** sempre pode ser reanalisado (o log ainda cresce).
- A análise na TUI roda **em segundo plano** com um indicador de progresso (não trava a tela).

- **Logs:** `~/.harnessbeaver/logs/AAAA-MM-DD.jsonl` (comando, shell, cwd, exit code,
  duração, stdout/stderr — cada stream truncado em 10KB).
- **Análises:** `~/.harnessbeaver/learnings/AAAA-MM-DD.md`.
- **Motor de análise** (`settings.insightsEngine`): `auto` (default — usa o `claude` CLI
  se presente, senão a API), `claude`, ou `api` (requer `ANTHROPIC_API_KEY`;
  modelo em `settings.insightsModel`).
- **Auto-oferta:** ao abrir o `bvr` (TUI) ou `bvr shell` num novo dia, se houver log do
  dia anterior sem análise, ele pergunta se quer analisar (uma vez por dia). O `bvr run`
  só exibe uma dica não-bloqueante.
- **Escopo:** o `bvr` registra apenas o que **ele** executa (`run`/`shell`), não o que é
  digitado dentro das abas abertas pelo launcher.
- O shell default do runner é `settings.defaultRunShell` (`pwsh`).

### Modos e shells

- **Modo** `tabs`: uma janela do Windows Terminal com uma aba por projeto.
- **Modo** `windows`: uma janela separada por projeto.
- **Shell** por aba: `claude` (roda Claude Code), `pwsh`, `cmd`.

Precedência de modo: default global → pacote → flag `--mode`.
Precedência de shell: default global → shell do projeto → pacote/flag.

## Configuração

Persistida em `%USERPROFILE%\.harnessbeaver\config.json` (criada no primeiro uso).
Contém `settings`, `projects` e `packages`. Edite as settings pela CLI ou pela tela
**Configurações** da TUI (sem mexer no JSON à mão):

```powershell
bvr config                                   # lista todas as settings
bvr config get insightsEngine
bvr config set insightsModel claude-opus-4-8
bvr config set learningsExtraDir F:\02-company-os\_content\learnings  # cópia extra das análises
bvr config set learningsExtraDir ""          # desliga a cópia extra
```

Chaves: `scanRoot`, `defaultMode` (tabs|windows), `defaultShell` (claude|pwsh|cmd),
`defaultRunShell` (pwsh|cmd|bash|zsh; default OS-aware), `insightsEngine` (auto|claude|api),
`insightsModel`, `learningsExtraDir`, `upgradeSource`, `tabColors` (true|false — cor diferente
por aba/janela aberta).

## Manutenção

```powershell
bvr version                       # versão do binário
bvr completion powershell|bash|zsh   # script de autocompletar
bvr upgrade --source <url|caminho>   # baixa e troca o binário (ou settings.upgradeSource)
```

## Detecção de projetos (scan)

O scanner varre `scanRoot` até 3 níveis e marca como projeto de código diretórios
com: `.git`, `package.json`, `go.mod`, `*.sln`, `*.csproj`, `pyproject.toml`,
`Cargo.toml`, `requirements.txt`, `src/`, `CLAUDE.md` ou `AGENTS.md`. Ao achar um
marcador, poda o ramo (não desce em subpastas do mesmo repositório).

**Grupos (diretórios-pai):** o `bvr scan --groups` também oferece os diretórios-pai
que agrupam ≥2 projetos (ex.: `produtos/smb`, `produtos/engenharia`, e a própria raiz),
marcados como `grupo: N projetos`. Importá-los registra o diretório-pai como um projeto,
permitindo abrir uma ferramenta (claude/pwsh/…) no nível que agrupa vários projetos. Na
TUI, a tela **Escanear & importar** já lista os grupos junto com os projetos.

# HarnessBeaver (`bvr`)

A launcher to open **Claude Code** (or PowerShell/cmd) across many directories at
once — in **tabs** of Windows Terminal or in **separate windows**. It's a dynamic
shortcut system for directories: you register projects, group them into packages, and
open everything with a single command.

It has two faces over the same core:

- **Rich TUI** (`bvr` with no arguments) — an interactive, navigable menu.
- **Classic CLI** (`bvr <command>`) — scriptable, for automation.

## Requirements

- **To run:** [Windows Terminal](https://aka.ms/terminal) (`wt`) and, for the `claude`
  shell, the Claude Code CLI on your PATH. **No runtime needed** — `bvr.exe` is static.
- **To build:** [Go](https://go.dev/dl/) 1.22+ (`winget install GoLang.Go`).

## Build

```powershell
go mod tidy
.\build.ps1            # produces bvr.exe
.\build.ps1 -All       # cross-compiles for win/linux/mac into ./dist
```

Distribute `bvr.exe`: the recipient just drops it in a folder on `PATH` and runs it.

## Usage — TUI

```powershell
bvr        # or: bvr ui
```

Menu: Open projects · Open package · Add project · Remove project ·
Scan & import · Create package · Remove package · Run command ·
Analyze learnings · Settings.

In any list, press **`/`** to **filter** by name/description (handy for long lists);
`enter` applies the filter, `esc` clears it. Navigation: `↑/↓` item, `←/→` (or PgUp/PgDn)
pages, `Home/End` start/end.
Keys: `↑/↓` navigate, `space` toggle (multi-select), `enter` confirm, `esc` back,
`q`/`ctrl+c` quit.

## Usage — CLI

```powershell
bvr list [projects|packages]            # list the registry
bvr add "My App" F:\proj\my-app --shell claude
bvr edit my-app --name "Other" --path F:\new --shell pwsh
bvr remove my-app
bvr scan --root F:\02-company-os\produtos --import   # import code projects
bvr pkg add "SMB Active" --projects os247,beauty --mode tabs --shell claude
bvr pkg remove smb-active
bvr open os247 beauty --mode windows --shell pwsh    # open projects/packages
bvr open smb-active --dry-run                         # show the wt commands
bvr status                                           # git status of all projects
bvr status smb-active                                 # git status of a package/ids
bvr each smb-active -- git pull                       # run a command in each project
bvr each repair beauty --parallel -- git fetch       # same, in parallel
```

`bvr each <ids...> -- <command>` runs the command in each target project's directory
(package or ids), logging every execution to the journal. Sequential mode shows live output;
`--parallel` captures and prints per project. Exits with a non-zero code if any project fails.

### Layouts (tabs/panes per project)

A **layout** opens a project in one window with several tabs and/or panes (split-panes),
each with its own shell and command — e.g. one pane running `claude`, another running the
dev server, another tab with `git`.

```powershell
bvr layout add dev --preset dev          # create layout (presets: claude|dev|triple)
bvr layout tab dev --title git           # add a tab
bvr layout pane dev --tab 2 --shell pwsh --command "git status" --split V
bvr layout show dev --project repair     # show the resulting wt command
bvr layout assign dev repair             # set the project's default layout
bvr open repair --layout dev             # open with the layout (or just 'bvr open repair' if it's the default)
bvr open repair --no-layout              # ignore the default layout
bvr layout list                          # list; bvr layout rm <id> removes
```

In the TUI, **Open with layout** picks the layout and the project. Each pane uses the
project's directory by default; `--split H` (horizontal) or `V` (vertical) defines how the
pane splits the previous one. Panes split the most recently created pane sequentially.

`bvr status` shows, per project, the branch, whether there are pending changes (`✗`/`✓`)
and ahead/behind of upstream (`↑`/`↓`). In the **TUI**, this status appears (loaded in the
background) in each project's description when Opening/Removing/building a package.

## Journal and learnings

`bvr` also runs commands on your behalf, **logs** everything to a daily journal and, at
the start of each new day, **offers to analyze the previous day's learnings** using AI.

```powershell
bvr run echo hello                     # runs in the default shell and logs it
bvr run --shell cmd "echo hi & dir"    # cmd; also: --shell pwsh|bash|zsh
bvr shell                              # REPL: type commands one after another, 'exit' to quit
bvr review                             # analyze the pending day (or show today's cache)
bvr review 2026-06-02                  # if already analyzed, show the cache
bvr review 2026-06-02 --force          # re-run the analysis even with a cache
bvr review --dry-run                   # show the prompt without calling the AI (no cost)

bvr history                            # recent commands (newest first)
bvr history --project repair --failed  # filter by project and only the failed ones
bvr history --grep "git" --days 30     # filter by text, last 30 days

bvr alias set deploy "npm run build && npm run deploy"   # create a shortcut
bvr alias                              # list shortcuts
bvr do deploy                          # run the shortcut
bvr alias set commit "git commit -m {1}"
bvr do commit "fix: tweak"             # {1}..{N} and {*} receive the arguments
```

Each command run by `bvr` is **associated with a project** (when the directory falls
within a registered project), enabling `bvr history --project <id>` and project context in
the analysis. In the **TUI**, there are **History** (re-runs the chosen command) and
**Shortcuts** (runs an alias) screens — both with the `/` filter.

- **Cache:** once analyzed, the day is saved. `bvr review <day>` shows the cache;
  use `--force` to re-analyze. In the TUI, picking an already-analyzed day offers **view the
  cache** or **re-analyze**. The **current day** can always be re-analyzed (the log is still growing).
- The TUI analysis runs **in the background** with a progress indicator (it doesn't block the screen).

- **Logs:** `~/.harnessbeaver/logs/YYYY-MM-DD.jsonl` (command, shell, cwd, exit code,
  duration, stdout/stderr — each stream truncated at 10KB).
- **Analyses:** `~/.harnessbeaver/learnings/YYYY-MM-DD.md`.
- **Analysis engine** (`settings.insightsEngine`): `auto` (default — uses the `claude` CLI
  if present, otherwise the API), `claude`, or `api` (requires `ANTHROPIC_API_KEY`;
  model in `settings.insightsModel`).
- **Auto-offer:** when opening `bvr` (TUI) or `bvr shell` on a new day, if there's a log
  from the previous day without an analysis, it asks whether to analyze (once per day). `bvr run`
  only shows a non-blocking hint.
- **Scope:** `bvr` logs only what **it** runs (`run`/`shell`), not what is typed
  inside the tabs opened by the launcher.
- The runner's default shell is `settings.defaultRunShell` (`pwsh`).

### Modes and shells

- **Mode** `tabs`: one Windows Terminal window with one tab per project.
- **Mode** `windows`: a separate window per project.
- **Shell** per tab: `claude` (runs Claude Code), `pwsh`, `cmd`.

Mode precedence: global default → package → `--mode` flag.
Shell precedence: global default → project shell → package/flag.

## Configuration

Persisted in `%USERPROFILE%\.harnessbeaver\config.json` (created on first use).
Contains `settings`, `projects` and `packages`. Edit settings via the CLI or the
**Settings** screen in the TUI (without touching the JSON by hand):

```powershell
bvr config                                   # list all settings
bvr config get insightsEngine
bvr config set insightsModel claude-opus-4-8
bvr config set learningsExtraDir F:\02-company-os\_content\learnings  # extra copy of the analyses
bvr config set learningsExtraDir ""          # turn off the extra copy
```

Keys: `scanRoot`, `defaultMode` (tabs|windows), `defaultShell` (claude|pwsh|cmd),
`defaultRunShell` (pwsh|cmd|bash|zsh; OS-aware default), `insightsEngine` (auto|claude|api),
`insightsModel`, `learningsExtraDir`, `upgradeSource`, `tabColors` (true|false — a different
color per opened tab/window).

## Maintenance

```powershell
bvr version                       # binary version
bvr completion powershell|bash|zsh   # autocompletion script
bvr upgrade --source <url|path>      # downloads and swaps the binary (or settings.upgradeSource)
```

## Project detection (scan)

The scanner walks `scanRoot` up to 3 levels deep and marks as a code project any directory
with: `.git`, `package.json`, `go.mod`, `*.sln`, `*.csproj`, `pyproject.toml`,
`Cargo.toml`, `requirements.txt`, `src/`, `CLAUDE.md` or `AGENTS.md`. When it finds a
marker, it prunes the branch (it doesn't descend into subfolders of the same repository).

**Groups (parent directories):** `bvr scan --groups` also offers the parent directories
that group ≥2 projects (e.g. `produtos/smb`, `produtos/engenharia`, and the root itself),
marked as `group: N projects`. Importing them registers the parent directory as a project,
letting you open a tool (claude/pwsh/…) at the level that groups several projects. In the
TUI, the **Scan & import** screen already lists the groups alongside the projects.

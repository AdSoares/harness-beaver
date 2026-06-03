// Package runner executa comandos em PowerShell ou cmd em nome do usuário,
// espelhando a saída ao vivo (opcional) e capturando-a para o diário.
package runner

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"time"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/journal"
)

// Result resume a execução de um comando.
type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
}

// ResolveProjectID, se definido, mapeia um cwd para um id de projeto, gravado
// no diário para permitir histórico/estatísticas por projeto. Configurado pela
// camada cmd/tui após carregar a config.
var ResolveProjectID func(cwd string) string

// shellArgs mapeia o shell para o executável e seus argumentos de invocação.
func shellArgs(shell config.Shell, command string) (string, []string) {
	switch shell {
	case config.ShellCmd:
		return "cmd", []string{"/C", command}
	case config.ShellBash:
		return "bash", []string{"-lc", command}
	case config.ShellZsh:
		return "zsh", []string{"-lc", command}
	case config.ShellPwsh:
		fallthrough
	default:
		return "pwsh", []string{"-NoProfile", "-Command", command}
	}
}

// normalizeShell garante um shell de execução válido (pwsh/cmd/bash/zsh); shells
// vazios ou de launcher (claude) caem no default do sistema operacional.
func normalizeShell(shell config.Shell) config.Shell {
	switch shell {
	case config.ShellCmd, config.ShellPwsh, config.ShellBash, config.ShellZsh:
		return shell
	default:
		return config.DefaultRunShellForOS()
	}
}

// Run executa command no shell indicado a partir de cwd. Se stream for true,
// a saída também é espelhada em os.Stdout/os.Stderr ao vivo. Em qualquer caso
// a execução é registrada no diário.
func Run(shell config.Shell, command, cwd string, stream bool) (Result, error) {
	shell = normalizeShell(shell)
	name, args := shellArgs(shell, command)

	if cwd == "" {
		if wd, err := os.Getwd(); err == nil {
			cwd = wd
		}
	}

	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Dir = cwd
	if stream {
		// Terminal real (run/shell): herda stdin e espelha a saída ao vivo.
		cmd.Stdin = os.Stdin
		cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	} else {
		// Captura (TUI): sem stdin para não disputar o terminal com a TUI.
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
	}

	start := time.Now()
	runErr := cmd.Run()
	dur := time.Since(start)

	exitCode := 0
	if runErr != nil {
		if ee, ok := runErr.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			// Falha ao iniciar o processo (ex.: shell ausente).
			exitCode = -1
			errBuf.WriteString(runErr.Error())
		}
	}

	res := Result{
		ExitCode: exitCode,
		Stdout:   outBuf.String(),
		Stderr:   errBuf.String(),
		Duration: dur,
	}

	projectID := ""
	if ResolveProjectID != nil {
		projectID = ResolveProjectID(cwd)
	}

	_ = journal.Append(journal.Entry{
		Time:       start,
		Shell:      string(shell),
		Cwd:        cwd,
		ProjectID:  projectID,
		Command:    command,
		ExitCode:   exitCode,
		DurationMs: dur.Milliseconds(),
		Stdout:     res.Stdout,
		Stderr:     res.Stderr,
	})

	return res, nil
}

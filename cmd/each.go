package cmd

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/runner"
)

var (
	eachShell    string
	eachParallel bool
)

var eachCmd = &cobra.Command{
	Use:   "each <pkgId|projId...> -- <comando...>",
	Short: "Roda um comando em cada projeto (de um pacote ou ids)",
	Long: `Executa o mesmo comando no diretório de cada projeto alvo.
Os alvos (ids de projeto e/ou pacote) vêm antes de '--' e o comando vem depois.
Cada execução é registrada no diário (com o cwd do projeto).

Exemplos:
  bvr each smb-ativo -- git pull
  bvr each repair beauty --shell bash -- git status -s
  bvr each smb-ativo --parallel -- git fetch`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dash := cmd.ArgsLenAtDash()
		if dash < 0 {
			return fmt.Errorf("use '--' para separar alvos do comando: bvr each <ids...> -- <comando...>")
		}
		targetArgs := args[:dash]
		cmdArgs := args[dash:]
		if len(targetArgs) == 0 {
			return fmt.Errorf("informe ao menos um projeto/pacote antes de '--'")
		}
		if len(cmdArgs) == 0 {
			return fmt.Errorf("informe o comando após '--'")
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		projects, err := collectProjects(cfg, targetArgs)
		if err != nil {
			return err
		}
		if len(projects) == 0 {
			fmt.Println("Nenhum projeto alvo.")
			return nil
		}

		shell := cfg.Settings.DefaultRunShell
		if eachShell != "" {
			s, err := validateRunShell(eachShell)
			if err != nil {
				return err
			}
			shell = s
		}
		command := strings.Join(cmdArgs, " ")

		var failures int
		if eachParallel {
			failures = runEachParallel(projects, shell, command)
		} else {
			failures = runEachSequential(projects, shell, command)
		}

		fmt.Printf("\nResumo: %d projeto(s), %d com erro (exit != 0).\n", len(projects), failures)
		if failures > 0 {
			os.Exit(1)
		}
		return nil
	},
}

// runEachSequential roda projeto a projeto com saída ao vivo.
func runEachSequential(projects []config.Project, shell config.Shell, command string) int {
	failures := 0
	for _, p := range projects {
		fmt.Printf("\n=== %s · %s ===\n", p.Name, p.Path)
		res, _ := runner.Run(shell, command, p.Path, true)
		fmt.Printf("(exit %d, %dms)\n", res.ExitCode, res.Duration.Milliseconds())
		if res.ExitCode != 0 {
			failures++
		}
	}
	return failures
}

// runEachParallel roda em paralelo (captura) e imprime os blocos na ordem dos projetos.
func runEachParallel(projects []config.Project, shell config.Shell, command string) int {
	results := make([]runner.Result, len(projects))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 4)
	for i, p := range projects {
		wg.Add(1)
		go func(i int, p config.Project) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i], _ = runner.Run(shell, command, p.Path, false)
		}(i, p)
	}
	wg.Wait()

	failures := 0
	for i, p := range projects {
		r := results[i]
		fmt.Printf("\n=== %s · %s === (exit %d, %dms)\n", p.Name, p.Path, r.ExitCode, r.Duration.Milliseconds())
		if out := strings.TrimRight(r.Stdout, "\n"); out != "" {
			fmt.Println(out)
		}
		if errOut := strings.TrimRight(r.Stderr, "\n"); errOut != "" {
			fmt.Println(errOut)
		}
		if r.ExitCode != 0 {
			failures++
		}
	}
	return failures
}

func init() {
	eachCmd.Flags().StringVar(&eachShell, "shell", "", "shell de execução: pwsh|cmd|bash|zsh")
	eachCmd.Flags().BoolVar(&eachParallel, "parallel", false, "roda em paralelo (saída capturada por projeto)")
	rootCmd.AddCommand(eachCmd)
}

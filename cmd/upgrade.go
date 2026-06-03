package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
)

var upgradeSourceFlag string

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Atualiza o binário do bvr a partir de uma fonte (URL ou caminho)",
	Long: `Baixa o binário novo de --source (ou settings.upgradeSource) e substitui o
executável atual. A fonte pode ser uma URL http(s) ou um caminho de arquivo
local/compartilhado. O binário anterior é mantido como backup (.old).`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		source := upgradeSourceFlag
		if source == "" {
			source = cfg.Settings.UpgradeSource
		}
		if source == "" {
			return fmt.Errorf("nenhuma fonte: use --source <url|caminho> ou 'bvr config set upgradeSource <...>'")
		}

		exe, err := os.Executable()
		if err != nil {
			return err
		}
		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			exe = resolved
		}

		newPath := exe + ".new"
		if err := fetchBinary(source, newPath); err != nil {
			return fmt.Errorf("falha ao obter binário: %w", err)
		}
		_ = os.Chmod(newPath, 0o755)

		// Swap Windows-safe: renomear o exe em uso é permitido; apagá-lo/sobrescrevê-lo não.
		oldPath := exe + ".old"
		_ = os.Remove(oldPath)
		if err := os.Rename(exe, oldPath); err != nil {
			_ = os.Remove(newPath)
			return fmt.Errorf("não foi possível mover o binário atual: %w", err)
		}
		if err := os.Rename(newPath, exe); err != nil {
			_ = os.Rename(oldPath, exe) // tenta restaurar
			_ = os.Remove(newPath)
			return fmt.Errorf("não foi possível instalar o novo binário: %w", err)
		}
		fmt.Printf("Atualizado: %s\n(backup do anterior em %s — pode apagar)\n", exe, oldPath)
		return nil
	},
}

// fetchBinary baixa (http/https) ou copia (caminho local) a fonte para dest.
func fetchBinary(source, dest string) error {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		client := &http.Client{Timeout: 5 * time.Minute}
		resp, err := client.Get(source)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			return fmt.Errorf("status %d", resp.StatusCode)
		}
		return copyToFile(dest, resp.Body)
	}
	f, err := os.Open(source)
	if err != nil {
		return err
	}
	defer f.Close()
	return copyToFile(dest, f)
}

func copyToFile(dest string, r io.Reader) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, r)
	return err
}

func init() {
	upgradeCmd.Flags().StringVar(&upgradeSourceFlag, "source", "", "URL ou caminho do binário novo")
	rootCmd.AddCommand(upgradeCmd)
}

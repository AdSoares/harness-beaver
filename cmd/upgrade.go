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
	"harnessbeaver/internal/i18n"
)

var upgradeSourceFlag string

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: i18n.T("Update the bvr binary from a source (URL or path)"),
	Long:  i18n.T("Download the new binary from --source (or settings.upgradeSource) and replace\nthe current executable. The source can be an http(s) URL or a local/shared file\npath. The previous binary is kept as a backup (.old)."),
	Args:  cobra.NoArgs,
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
			return i18n.Errorf("no source: use --source <url|path> or 'bvr config set upgradeSource <...>'")
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
			return fmt.Errorf(i18n.T("failed to fetch binary: %w"), err)
		}
		_ = os.Chmod(newPath, 0o755)

		// Swap Windows-safe: renomear o exe em uso é permitido; apagá-lo/sobrescrevê-lo não.
		oldPath := exe + ".old"
		_ = os.Remove(oldPath)
		if err := os.Rename(exe, oldPath); err != nil {
			_ = os.Remove(newPath)
			return fmt.Errorf(i18n.T("could not move the current binary: %w"), err)
		}
		if err := os.Rename(newPath, exe); err != nil {
			_ = os.Rename(oldPath, exe) // tenta restaurar
			_ = os.Remove(newPath)
			return fmt.Errorf(i18n.T("could not install the new binary: %w"), err)
		}
		fmt.Printf(i18n.T("Updated: %s\n(backup of the previous binary at %s — safe to delete)\n"), exe, oldPath)
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
			return fmt.Errorf(i18n.T("status %d"), resp.StatusCode)
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
	upgradeCmd.Flags().StringVar(&upgradeSourceFlag, "source", "", i18n.T("URL or path of the new binary"))
	rootCmd.AddCommand(upgradeCmd)
}

package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
)

func keyList() string {
	keys := make([]string, 0, len(config.SettingDefs))
	for _, d := range config.SettingDefs {
		keys = append(keys, d.Key)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

var configCmd = &cobra.Command{
	Use:   "config [get <chave> | set <chave> <valor>]",
	Short: "Mostra ou altera as configurações do bvr",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		switch {
		case len(args) == 0:
			for _, d := range config.SettingDefs {
				fmt.Printf("  %-18s %s\n", d.Key, dash(d.Get(cfg)))
			}
			return nil

		case args[0] == "get":
			if len(args) != 2 {
				return fmt.Errorf("uso: bvr config get <chave>")
			}
			d, ok := config.FindSetting(args[1])
			if !ok {
				return fmt.Errorf("chave desconhecida %q (válidas: %s)", args[1], keyList())
			}
			fmt.Println(d.Get(cfg))
			return nil

		case args[0] == "set":
			if len(args) < 2 {
				return fmt.Errorf("uso: bvr config set <chave> <valor>")
			}
			d, ok := config.FindSetting(args[1])
			if !ok {
				return fmt.Errorf("chave desconhecida %q (válidas: %s)", args[1], keyList())
			}
			value := strings.Join(args[2:], " ") // permite valor vazio p/ desligar
			if err := d.Apply(cfg, value); err != nil {
				return err
			}
			if err := cfg.Save(); err != nil {
				return err
			}
			fmt.Printf("%s = %s\n", d.Key, dash(d.Get(cfg)))
			return nil
		}
		return fmt.Errorf("subcomando inválido %q (use: get, set, ou nada para listar)", args[0])
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

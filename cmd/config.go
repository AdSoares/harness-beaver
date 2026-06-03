package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
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
	Use:   i18n.T("config [get <key> | set <key> <value>]"),
	Short: i18n.T("Show or change bvr settings"),
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
				return i18n.Errorf("usage: bvr config get <key>")
			}
			d, ok := config.FindSetting(args[1])
			if !ok {
				return fmt.Errorf(i18n.T("unknown key %q (valid: %s)"), args[1], keyList())
			}
			fmt.Println(d.Get(cfg))
			return nil

		case args[0] == "set":
			if len(args) < 2 {
				return i18n.Errorf("usage: bvr config set <key> <value>")
			}
			d, ok := config.FindSetting(args[1])
			if !ok {
				return fmt.Errorf(i18n.T("unknown key %q (valid: %s)"), args[1], keyList())
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
		return fmt.Errorf(i18n.T("invalid subcommand %q (use: get, set, or nothing to list)"), args[0])
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

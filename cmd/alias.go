package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"harnessbeaver/internal/config"
	"harnessbeaver/internal/i18n"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: i18n.T("Manage command shortcuts (use 'bvr do <name>' to run)"),
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if len(cfg.Aliases) == 0 {
			fmt.Println(i18n.T("No shortcuts. Create one with: bvr alias set <name> <command...>"))
			return nil
		}
		names := make([]string, 0, len(cfg.Aliases))
		for n := range cfg.Aliases {
			names = append(names, n)
		}
		sort.Strings(names)
		for _, n := range names {
			fmt.Printf("  %-16s %s\n", n, cfg.Aliases[n])
		}
		return nil
	},
}

var aliasSetCmd = &cobra.Command{
	Use:   i18n.T("set <name> <command...>"),
	Short: i18n.T("Create or update a shortcut"),
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		name := args[0]
		command := strings.Join(args[1:], " ")
		cfg.SetAlias(name, command)
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("%s = %s\n", name, command)
		return nil
	},
}

var aliasRmCmd = &cobra.Command{
	Use:     i18n.T("rm <name>"),
	Aliases: []string{"remove"},
	Short:   i18n.T("Remove a shortcut"),
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if !cfg.RemoveAlias(args[0]) {
			return fmt.Errorf(i18n.T("shortcut not found: %s"), args[0])
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf(i18n.T("Shortcut removed: %s\n"), args[0])
		return nil
	},
}

func init() {
	aliasCmd.AddCommand(aliasSetCmd, aliasRmCmd)
	rootCmd.AddCommand(aliasCmd)
}

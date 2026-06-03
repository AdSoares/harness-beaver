package config

import (
	"fmt"
	"strconv"
	"strings"
)

// SettingDef descreve uma chave editável de Settings, com leitura/escrita por
// string e, opcionalmente, um conjunto de valores válidos (enum). É a fonte
// única usada pela CLI (`bvr config`) e pela TUI (tela Configurações).
type SettingDef struct {
	Key  string
	Help string
	Enum []string // vazio = texto livre
	Get  func(*Config) string
	Set  func(*Config, string)
}

// Apply valida (contra Enum, quando houver) e grava o valor no Config.
func (d SettingDef) Apply(c *Config, v string) error {
	if len(d.Enum) > 0 {
		ok := false
		for _, e := range d.Enum {
			if e == v {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("valor inválido %q (use: %s)", v, strings.Join(d.Enum, ", "))
		}
	}
	d.Set(c, v)
	return nil
}

func shellStrings(ss []Shell) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = string(s)
	}
	return out
}

// SettingDefs lista todas as chaves de configuração editáveis.
var SettingDefs = []SettingDef{
	{
		Key:  "scanRoot",
		Help: "diretório raiz do scan de projetos",
		Get:  func(c *Config) string { return c.Settings.ScanRoot },
		Set:  func(c *Config, v string) { c.Settings.ScanRoot = v },
	},
	{
		Key:  "defaultMode",
		Help: "modo padrão de abertura",
		Enum: []string{string(ModeTabs), string(ModeWindows)},
		Get:  func(c *Config) string { return string(c.Settings.DefaultMode) },
		Set:  func(c *Config, v string) { c.Settings.DefaultMode = Mode(v) },
	},
	{
		Key:  "defaultShell",
		Help: "shell padrão do launcher",
		Enum: shellStrings(ValidShells),
		Get:  func(c *Config) string { return string(c.Settings.DefaultShell) },
		Set:  func(c *Config, v string) { c.Settings.DefaultShell = Shell(v) },
	},
	{
		Key:  "defaultRunShell",
		Help: "shell padrão do run/shell",
		Enum: shellStrings(ValidRunShells),
		Get:  func(c *Config) string { return string(c.Settings.DefaultRunShell) },
		Set:  func(c *Config, v string) { c.Settings.DefaultRunShell = Shell(v) },
	},
	{
		Key:  "insightsEngine",
		Help: "motor da análise de aprendizados",
		Enum: []string{"auto", "claude", "api"},
		Get:  func(c *Config) string { return c.Settings.InsightsEngine },
		Set:  func(c *Config, v string) { c.Settings.InsightsEngine = v },
	},
	{
		Key:  "insightsModel",
		Help: "modelo usado no caminho da API",
		Get:  func(c *Config) string { return c.Settings.InsightsModel },
		Set:  func(c *Config, v string) { c.Settings.InsightsModel = v },
	},
	{
		Key:  "learningsExtraDir",
		Help: "pasta extra p/ copiar as análises (vazio = desligado)",
		Get:  func(c *Config) string { return c.Settings.LearningsExtraDir },
		Set:  func(c *Config, v string) { c.Settings.LearningsExtraDir = v },
	},
	{
		Key:  "upgradeSource",
		Help: "URL ou caminho do binário p/ 'bvr upgrade'",
		Get:  func(c *Config) string { return c.Settings.UpgradeSource },
		Set:  func(c *Config, v string) { c.Settings.UpgradeSource = v },
	},
	{
		Key:  "tabColors",
		Help: "cor diferente por aba/janela aberta",
		Enum: []string{"true", "false"},
		Get:  func(c *Config) string { return strconv.FormatBool(c.Settings.TabColors) },
		Set:  func(c *Config, v string) { c.Settings.TabColors = v == "true" },
	},
}

// FindSetting devolve a definição de uma chave.
func FindSetting(key string) (SettingDef, bool) {
	for _, d := range SettingDefs {
		if d.Key == key {
			return d, true
		}
	}
	return SettingDef{}, false
}

package cmd

import (
	"fmt"
	"github.com/mdreem/json2line/configuration"
	"github.com/mdreem/json2line/configuration/load"
	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
	"path/filepath"
	"testing"
)

func TestConfigureCommandBufferSize(t *testing.T) {
	defer func() { viper.Reset() }()
	const bufferSize = "1000"

	tempDir := t.TempDir()
	initRealUserConfigDirFunc := configuration.UserConfigDir
	loadRealUserConfigDirFunc := load.UserConfigDir

	mockUserConfigDirFunc := func() (string, error) {
		return tempDir, nil
	}
	configuration.UserConfigDir = mockUserConfigDirFunc
	load.UserConfigDir = mockUserConfigDirFunc

	defer func() {
		configuration.UserConfigDir = initRealUserConfigDirFunc
		load.UserConfigDir = loadRealUserConfigDirFunc
	}()

	RootCmd.SetArgs([]string{"configure", "init"})
	_, _ = executeCommandWithInput(t, "")

	RootCmd.SetArgs([]string{"configure", "buffer_size", "-S", bufferSize})
	_, _ = executeCommandWithInput(t, "")

	configToml := filepath.Join(tempDir, "json2line", "json2line.toml")

	fmt.Printf("Loading: %s\n", configToml)
	tree, err := toml.LoadFile(configToml)
	if err != nil {
		t.Fatal(err)
	}
	configurationSection := tree.Get("configuration").(*toml.Tree)
	if configurationSection == nil {
		t.Fatal("configuration section has not been created.")
	}

	bufferSizeFromConfig := configurationSection.Get("buffer_size").(string)

	if bufferSizeFromConfig != bufferSize {
		t.Fatalf("buffer size is '%s'. Expected was '%s'", bufferSizeFromConfig, bufferSize)
	}
}

func TestConfigureCommandFormatter(t *testing.T) {
	defer func() { viper.Reset() }()
	const formatterKey = "some_key"
	const formatterValue = "some_value"

	tempDir := t.TempDir()
	initRealUserConfigDirFunc := configuration.UserConfigDir
	loadRealUserConfigDirFunc := load.UserConfigDir

	mockUserConfigDirFunc := func() (string, error) {
		return tempDir, nil
	}
	configuration.UserConfigDir = mockUserConfigDirFunc
	load.UserConfigDir = mockUserConfigDirFunc

	defer func() {
		configuration.UserConfigDir = initRealUserConfigDirFunc
		load.UserConfigDir = loadRealUserConfigDirFunc
	}()

	RootCmd.SetArgs([]string{"configure", "init"})
	_, _ = executeCommandWithInput(t, "")

	RootCmd.SetArgs([]string{"configure", "formatter", "-k", formatterKey, "-v", formatterValue})
	_, _ = executeCommandWithInput(t, "")

	configToml := filepath.Join(tempDir, "json2line", "json2line.toml")

	fmt.Printf("Loading: %s\n", configToml)
	tree, err := toml.LoadFile(configToml)
	if err != nil {
		t.Fatal(err)
	}
	configurationSection := tree.Get("templates").(*toml.Tree)
	if configurationSection == nil {
		t.Fatal("templates section has not been created.")
	}

	formatterValueFromConfig := configurationSection.Get(formatterKey).(string)

	if formatterValueFromConfig != formatterValue {
		t.Fatalf("formatter value is '%s'. Expected was '%s'", formatterValue, formatterValueFromConfig)
	}
}

func TestConfigureCommandReplacement(t *testing.T) {
	defer func() { viper.Reset() }()
	const replacementKey = "some_key"
	const replacementValue = "some_value"

	tempDir := t.TempDir()
	initRealUserConfigDirFunc := configuration.UserConfigDir
	loadRealUserConfigDirFunc := load.UserConfigDir

	mockUserConfigDirFunc := func() (string, error) {
		return tempDir, nil
	}
	configuration.UserConfigDir = mockUserConfigDirFunc
	load.UserConfigDir = mockUserConfigDirFunc

	defer func() {
		configuration.UserConfigDir = initRealUserConfigDirFunc
		load.UserConfigDir = loadRealUserConfigDirFunc
	}()

	RootCmd.SetArgs([]string{"configure", "init"})
	_, _ = executeCommandWithInput(t, "")

	RootCmd.SetArgs([]string{"configure", "replacement", "-k", replacementKey, "-v", replacementValue})
	_, _ = executeCommandWithInput(t, "")

	configToml := filepath.Join(tempDir, "json2line", "json2line.toml")

	fmt.Printf("Loading: %s\n", configToml)
	tree, err := toml.LoadFile(configToml)
	if err != nil {
		t.Fatal(err)
	}
	configurationSection := tree.Get("replacements").(*toml.Tree)
	if configurationSection == nil {
		t.Fatal("replacements section has not been created.")
	}

	replacementValueFromConfig := configurationSection.Get(replacementKey).(string)

	if replacementValueFromConfig != replacementValue {
		t.Fatalf("replacement value is '%s'. Expected was '%s'", replacementValue, replacementValueFromConfig)
	}
}

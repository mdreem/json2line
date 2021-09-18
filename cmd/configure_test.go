package cmd

import (
	"fmt"
	"github.com/mdreem/json2line/configuration"
	"github.com/mdreem/json2line/configuration/load"
	"github.com/pelletier/go-toml"
	"path/filepath"
	"testing"
)

const bufferSize = "1000"

func TestConfigureCommandBufferSize(t *testing.T) {
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

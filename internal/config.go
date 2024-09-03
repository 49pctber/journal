package journal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const JournalDir string = "journal"
const ConfigFilename string = "config.json"

var configFilename string

func init() {
	ConfigDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(filepath.Join(ConfigDir, JournalDir), 0766)
	if err != nil {
		panic(err)
	}

	// check for configuration file
	configFilename = filepath.Join(ConfigDir, JournalDir, ConfigFilename)
	_, err = os.Stat(configFilename)
	if os.IsNotExist(err) {
		// if the file doesn't exist, create it
		file, err := os.Create(configFilename)
		if err != nil {
			fmt.Println(configFilename)
			panic(err)
		}
		defer file.Close()
	} else if err != nil {
		fmt.Printf("Error checking file: %v\n", err)
		panic(err)
	}
}

type Config struct {
	EntryDirectory  string
	PreferredEditor string
}

func (config Config) SaveToFile() error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFilename, data, 0766)
}

func (config *Config) LoadFromFile() error {

	jsondata, err := os.ReadFile(configFilename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsondata, config)
	if err != nil {
		return err
	}

	return nil
}

func (config *Config) Configure() error {

	var user_input string

	// get directory for storing journal entries
	for {
		fmt.Println("Which directory should be used for storing your journal entries? Relative paths will be converted to absolute paths.")
		fmt.Scanln(&user_input)

		// convert to absolute path
		abs_path, err := filepath.Abs(user_input)
		if err != nil {
			fmt.Printf("There was an error converting %s to an absolute path.\n", abs_path)
			continue
		}

		// make directory if necessary
		info, err := os.Stat(abs_path)
		if os.IsNotExist(err) {
			fmt.Printf("Directory %s does not exist. Create? (Y/n) ", abs_path)
			fmt.Scanln(&user_input)
			confirmation := strings.ToLower(user_input)
			if confirmation != "y" && confirmation != "yes" {
				continue
			}

			err := os.MkdirAll(abs_path, 0766)
			if err != nil {
				fmt.Printf("There was an error creating %v: %v\n", abs_path, err)
				continue
			}
			config.EntryDirectory = abs_path
			break
		} else if err != nil {
			fmt.Printf("Error checking directory: %v\n", err)
			continue
		} else if info.IsDir() {
			fmt.Printf("Using %s\n", abs_path)
			config.EntryDirectory = abs_path
			break
		} else {
			fmt.Println("Path already exists but is not a directory.")
			continue
		}
	}

	// get text editor
	editor := os.Getenv("EDITOR")

	if len(editor) != 0 {
		fmt.Printf("Your EDITOR environment variable is set to %s. Use that for your text editor? (Y/n) ", editor)
		fmt.Scanln(&user_input)
		confirmation := strings.ToLower(user_input)
		if confirmation != "y" && confirmation != "yes" {
			editor = ""
		}
	}

	if len(editor) == 0 {
		fmt.Printf("Please enter the command you would like to use to open your text editor (e.g. notepad.exe, vim, nano, /path/to/editor): ")
		fmt.Scanln(&editor)
	}
	config.PreferredEditor = editor

	// save configuration
	return config.SaveToFile()
}

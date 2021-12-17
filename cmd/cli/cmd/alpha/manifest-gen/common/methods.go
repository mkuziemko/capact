package common

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"os"
	"path/filepath"
)

// AskForDirectory asks for a directory using suggestion options for suggesting the list of dirs
func AskForDirectory(msg string, defaultDir string) (string, error){
	chosenDir := ""
	directoryPrompt := &survey.Input{
		Message: msg,
		Suggest: func (toComplete string) []string {
			files, err := filepath.Glob(toComplete + "*")
			if err != nil {
				fmt.Println("Cannot getting the names of files")
				return nil
			}
			var dirs []string
			for _, match := range files {
				file, err := os.Stat(match)
				if err != nil {
					fmt.Println("Cannot getting the information about the file")
					return nil
				}
				if file.IsDir() {
					dirs = append(dirs, match)
				}
			}
			return dirs
		},
	}
	if defaultDir != "" {
		directoryPrompt.Default = defaultDir
	}

	err := survey.AskOne(directoryPrompt, &chosenDir)
	return chosenDir, err
}

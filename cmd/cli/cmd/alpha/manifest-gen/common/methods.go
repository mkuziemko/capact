package common

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
)

// ValidateFun is a function that validates the user's answers. It is used in the survey library.
type ValidateFun func(ans interface{}) error

// AskForDirectory asks for a directory. It suggests to a user the list of dirs that can be used.
func AskForDirectory(msg string, defaultDir string) (string, error) {
	chosenDir := ""
	directoryPrompt := &survey.Input{
		Message: msg,
		Suggest: func(toComplete string) []string {
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

// CreateManifestPath create a manifest path based on a manifest type and suffix.
func CreateManifestPath(manifestType string, suffix string) string {
	suffixes := map[string]string{
		AttributeManifest:      "attribute",
		TypeManifest:           "type",
		InterfaceManifest:      "interface",
		InterfaceGroupManifest: "interfaceGroup",
		ImplementationManifest: "implementation",
	}
	return "cap." + suffixes[manifestType] + "." + suffix
}

// GetDefaultMetadata creates a new Metadata object and sets default values.
func GetDefaultMetadata() Metadata {
	var metadata Metadata
	metadata.DocumentationURL = "https://example.com"
	metadata.SupportURL = "https://example.com"
	metadata.IconURL = "https://example.com/icon.png"
	metadata.Maintainers = []Maintainers{
		{
			Email: "dev@example.com",
			Name:  "Example Dev",
			URL:   "https://example.com",
		},
	}
	return metadata
}

//ManyValidators allow using many validators function in the Survey validator.
func ManyValidators(validateFuns []ValidateFun) func(ans interface{}) error {
	return func(ans interface{}) error {
		for _, fun := range validateFuns {
			err := fun(ans)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// ValidateURL validates a URL.
func ValidateURL(ans interface{}) error {
	if str, ok := ans.(string); !ok || (ans != "" && !isUrl(str)) {
		return errors.New("URL is not valid")
	}
	return nil
}

// ValidateEmail validates an email.
func ValidateEmail(ans interface{}) error {
	if str, ok := ans.(string); !ok || (ans != "" && !isEmail(str)) {
		return errors.New("email is not valid")
	}
	return nil
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Package common - Defines common constants and variables and utility functions
package common

import (
	"os"
	"path"
	"strconv"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// ===== [ Constants and Variables ] =====

var (
	bundle    *i18n.Bundle    = nil
	localizer *i18n.Localizer = nil
)

// ===== [ Types ] =====

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

func getPath(langPath string) string {
	if !path.IsAbs(langPath) {
		currPath, _ := os.Getwd()
		return path.Join(currPath, langPath)
	}
	return langPath
}

// ===== [ Public Functions ] =====

// LoadMessages - Load and configure the messages for specified languages
func LoadMessages(langPath string, langs []string) {
	bundle = i18n.NewBundle(language.Korean)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	logger.Debugf("language Folder: %s to %s", langPath, getPath(langPath))

	for _, lang := range langs {
		bundle.MustLoadMessageFile(getPath(langPath) + "/message." + lang + ".yaml")
	}
	ChangeLocalizer(langs[0])
}

// ChangeLocalizer - Creates the new localizer for specific language
func ChangeLocalizer(lang string) {
	localizer = i18n.NewLocalizer(bundle, lang)
}

// GetMessageByCode - Returns the message for specific code
func GetMessageByCode(code int) string {
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: strconv.Itoa(code),
	})
}

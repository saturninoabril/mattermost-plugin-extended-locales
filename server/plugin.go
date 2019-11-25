package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/plugin"
	"github.com/saturninoabril/mattermost-plugin-extended-locales/server/locale"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

func readJSONFromUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.Bytes(), nil
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/get_languages":
		p.handleGetLanguages(w, r)
	case "/get_translation":
		p.handleGetTranslation(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) handleGetLanguages(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	b, err := p.getLanguages()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(b)
}

func (p *Plugin) getLanguages() ([]byte, error) {
	config := p.getConfiguration()
	if config.EnableTranslationService && config.TranslationServiceURL != "" {
		url := fmt.Sprintf("%s/locales.json", config.TranslationServiceURL)
		return readJSONFromUrl(url)

	}

	// This ExtendedLocales is added to show maintainability issue when not using external translation service
	return json.Marshal(locale.ExtendedLocales)
}

func (p *Plugin) handleGetTranslation(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	lang := r.URL.Query().Get("lang")
	client := r.URL.Query().Get("client")

	config := p.getConfiguration()
	if config.EnableTranslationService && config.TranslationServiceURL != "" {
		b, err := getTranslationFromTranslationService(config.TranslationServiceURL, lang, client, p.API.GetServerVersion(), r.URL.Query().Get("app_version"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Write(b)
		return
	}

	// This block is added to show maintainability issue when not using external translation service
	b, err := getTranslationFromPlugin(lang, client)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(b)
}

func getTranslationFromTranslationService(baseUrl, lang, client, serverVersion, appVersion string) ([]byte, error) {
	var url string
	if client == "rn" {
		majorMobileAppVersion, minorMobileAppVersion, _ := SplitVersion(appVersion)
		url = fmt.Sprintf("%s/%s.%s/%s.json", baseUrl, majorMobileAppVersion, minorMobileAppVersion, lang)
	} else {
		majorServerVersion, minorServerVersion, _ := SplitVersion(serverVersion)
		url = fmt.Sprintf("%s/%s.%s/%s.json", baseUrl, majorServerVersion, minorServerVersion, lang)
	}

	b, err := readJSONFromUrl(url)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func getTranslationFromPlugin(lang, client string) ([]byte, error) {
	var b []byte
	var err error

	switch lang {
	case "tl":
		if client == "rn" {
			b, err = json.Marshal(locale.TagalogRN)
		} else {
			b, err = json.Marshal(locale.Tagalog)
		}
	case "no":
		if client == "rn" {
			b, err = json.Marshal(locale.NorwegianRN)
		} else {
			b, err = json.Marshal(locale.Norwegian)
		}
	default:
	}

	if err != nil {
		return nil, err
	}

	return b, nil
}

func SplitVersion(version string) (string, string, string) {
	parts := strings.Split(version, ".")

	return parts[0], parts[1], parts[2]
}

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

	fmt.Println(lang)
	fmt.Println(client)

	config := p.getConfiguration()
	if config.EnableTranslationService && config.TranslationServiceURL != "" {
		var url string
		if client == "rn" {
			majorMobileAppVersion, minorMobileAppVersion, _ := SplitVersion(r.URL.Query().Get("app_version"))
			url = fmt.Sprintf("%s/%s.%s/%s.json", config.TranslationServiceURL, majorMobileAppVersion, minorMobileAppVersion, lang)
		} else {
			majorServerVersion, minorServerVersion, _ := SplitVersion(p.API.GetServerVersion())
			url = fmt.Sprintf("%s/%s.%s/%s.json", config.TranslationServiceURL, majorServerVersion, minorServerVersion, lang)
		}

		b, err := readJSONFromUrl(url)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Write(b)
	} else {
		var b []byte
		var jsonErr error

		switch lang {
		case "tl":
			if client == "rn" {
				b, jsonErr = json.Marshal(locale.TagalogRN)
			} else {
				b, jsonErr = json.Marshal(locale.Tagalog)
			}
		case "no":
			if client == "rn" {
				b, jsonErr = json.Marshal(locale.NorwegianRN)
			} else {
				b, jsonErr = json.Marshal(locale.Norwegian)
			}
		default:
		}

		if jsonErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Write(b)
	}
}

func SplitVersion(version string) (string, string, string) {
	parts := strings.Split(version, ".")

	return parts[0], parts[1], parts[2]
}

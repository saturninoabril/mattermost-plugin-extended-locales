package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-server/plugin"
	"github.com/saturninoabril/mattermost-plugin-extended-locales/locale"
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
	fmt.Println(ExtendedLocales)

	b, jsonErr := json.Marshal(ExtendedLocales)
	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(b)
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

	var b []byte
	var jsonErr error

	switch lang {
	case "tl":
		if client == "rn" {
			b, jsonErr = json.Marshal(TagalogRN)
		} else {
			b, jsonErr = json.Marshal(Tagalog)
		}
	case "no":
		if client == "rn" {
			b, jsonErr = json.Marshal(NorwegianRN)
		} else {
			b, jsonErr = json.Marshal(Norwegian)
		}
	default:
	}

	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(b)
}

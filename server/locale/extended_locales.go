package locale

type Locale struct {
	Value string `json:"value"`
	Name  string `json:"name"`
	Order int    `json:"order"`
	URL   string `json:"url"`
}

var ExtendedLocales = []Locale{
	{
		Value: "no",
		Name:  "Norwegian",
		Order: 101,
		URL:   "/plugins/com.mattermost.plugin-extended-locales/get_translation?lang=no",
	},
	{
		Value: "tl",
		Name:  "Tagalog",
		Order: 102,
		URL:   "/plugins/com.mattermost.plugin-extended-locales/get_translation?lang=tl",
	},
}

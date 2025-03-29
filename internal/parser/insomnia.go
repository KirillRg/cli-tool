package parser

type InsomniaCollection struct {
	Type       string        `yaml:"type"`
	Name       string        `yaml:"name"`
	Collection []RequestItem `yaml:"collection"`
}

type RequestItem struct {
	URL        string          `yaml:"url"`
	Name       string          `yaml:"name"`
	Method     string          `yaml:"method"`
	Body       RequestBody     `yaml:"body"`
	Parameters []RequestParam  `yaml:"parameters"`
	Headers    []RequestHeader `yaml:"headers"`
	Settings   RequestSettings `yaml:"settings"`
}

type RequestBody struct {
	MimeType string `yaml:"mimeType"`
	Text     string `yaml:"text"`
}

type RequestParam struct {
	Name     string `yaml:"name"`
	Value    string `yaml:"value"`
	Disabled bool   `yaml:"disabled"`
}

type RequestHeader struct {
	Name     string `yaml:"name"`
	Value    string `yaml:"value"`
	Disabled bool   `yaml:"disabled"`
}

type RequestSettings struct {
	RenderRequestBody bool `yaml:"renderRequestBody"`
	EncodeURL         bool `yaml:"encodeUrl"`
}

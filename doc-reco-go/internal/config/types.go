package config

type config struct {
	ServerPort    string        `yaml:"server_port"`
	Elasticsearch elasticsearch `yaml:"elasticsearch"`
	Mathpix       mathpix       `yaml:"mathpix"`
	GoogleVision  googleVision  `yaml:"google_vision"`
	Sentry        sentry        `yaml:"sentry"`
	Datadog       datadog       `yaml:"datadog"`
}

func (c *config) SetDefault() {
	if c.ServerPort == "" {
		c.ServerPort = "4000"
	}
}

type elasticsearch struct {
	AutoSuggestion string `yaml:"auto_suggestions"`
	Search         string `yaml:"search"`
}

type mathpix struct {
	ApiUrl string `yaml:"api_url"`
	AppId  string `yaml:"app_id"`
	AppKey string `yaml:"app_key"`
}

type googleVision struct {
	ApiUrl string `yaml:"api_url"`
}

type sentry struct {
	Dsn   string `yaml:"dsn"`
	Env   string `yaml:"env"`
	Debug bool   `yaml:"debug"`
}

type datadog struct {
	Version     string `yaml:"version"`
	Env         string `yaml:"env"`
	ServiceName string `yaml:"service_name"`
	AgentHost   string `yaml:"agent_host"`
	AgentPort   string `yaml:"agent_port"`
}

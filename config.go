package traefik_hydrate_headers

type RemoteConfig struct {
	Url    string `yaml:"url"`
	Method string `yaml:"method"`
}

type ConditionConfig struct {
	StatusCodes []int `yaml:"statusCode"`
}

type Config struct {
	Remote         *RemoteConfig     `yaml:"remote"`
	AppendOn       *ConditionConfig  `yaml:"appendOn"`
	NextOn         *ConditionConfig  `yaml:"nextOn"`
	ForwardHeaders []string          `yaml:"forwardHeaders"`
	Headers        map[string]string `yaml:"headers"`
}

func CreateConfig() *Config {
	return &Config{
		Remote: &RemoteConfig{
			Url:    "",
			Method: "GET",
		},
		AppendOn: &ConditionConfig{
			StatusCodes: []int{},
		},
		NextOn: &ConditionConfig{
			StatusCodes: []int{},
		},
		ForwardHeaders: []string{},
		Headers:        make(map[string]string),
	}
}

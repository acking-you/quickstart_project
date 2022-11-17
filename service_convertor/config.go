package service_convertor

type Config struct {
	enableVOFileSingle bool
	enableTOFileSingle bool
	enableAutoFormTag  bool
	enableAutoJsonTag  bool
	enableDebug        bool
	savePath           string
	defaultMethodName  string
}

func DefaultConfig() *Config {
	return &Config{
		enableVOFileSingle: false,
		enableTOFileSingle: false,
		enableAutoFormTag:  true,
		enableAutoJsonTag:  true,
		enableDebug:        false,
		savePath:           "./service/",
		defaultMethodName:  "XXX",
	}
}

func (s *Config) EnableVOFileSingle(t bool) *Config {
	s.enableVOFileSingle = t
	return s
}

func (s *Config) EnableTOFileSingle(t bool) *Config {
	s.enableTOFileSingle = t
	return s
}

func (s *Config) EnableAutoFormTag(t bool) *Config {
	s.enableAutoFormTag = t
	return s
}

func (s *Config) EnableAutoJsonTag(t bool) *Config {
	s.enableAutoJsonTag = t
	return s
}

func (s *Config) EnableDebug(t bool) *Config {
	s.enableDebug = t
	return s
}

func (s *Config) SavePath(t string) *Config {
	s.savePath = t
	return s
}

func (s *Config) DefaultMethodName(t string) *Config {
	s.defaultMethodName = t
	return s
}

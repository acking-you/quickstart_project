package controller_convertor

const KPackageName = "controller"

type Config struct {
	savePath         string
	ginPackageName   string
	enableDebug      bool
	enableVOBind     bool //是否启用vo绑定
	enableResponse   bool //是否自动生成response
	enableGinExample bool //gin官方提供的起手例子
	enableCreate     bool //以下为CRUD生成
	enableQuery      bool
	enableUpdate     bool
	enableDelete     bool
}

func DefaultCConfig() *Config {
	return &Config{
		savePath:         "./controller/",
		ginPackageName:   "github.com/gin-gonic/gin",
		enableDebug:      false,
		enableVOBind:     true,
		enableResponse:   true,
		enableGinExample: true,
		enableCreate:     true,
		enableQuery:      true,
		enableUpdate:     true,
		enableDelete:     true,
	}
}

func (s *Config) SavePath(t string) *Config {
	s.savePath = t
	return s
}

func (s *Config) GinPackageName(t string) *Config {
	s.ginPackageName = t
	return s
}

func (s *Config) EnableDebug(t bool) *Config {
	s.enableDebug = t
	return s
}

func (s *Config) EnableVOBind(t bool) *Config {
	s.enableVOBind = t
	return s
}

func (s *Config) EnableResponse(t bool) *Config {
	s.enableResponse = t
	return s
}

func (s *Config) EnableGinExample(t bool) *Config {
	s.enableGinExample = t
	return s
}

func (s *Config) EnableCreate(t bool) *Config {
	s.enableCreate = t
	return s
}

func (s *Config) EnableQuery(t bool) *Config {
	s.enableQuery = t
	return s
}

func (s *Config) EnableUpdate(t bool) *Config {
	s.enableUpdate = t
	return s
}

func (s *Config) EnableDelete(t bool) *Config {
	s.enableDelete = t
	return s
}

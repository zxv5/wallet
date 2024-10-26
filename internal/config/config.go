package config

type Config struct {
	Server Server `json:"Server"`
	DBCfg  DBCfg  `json:"DBCfg"`
	Jwt    Jwt    `json:"Jwt"`
}

type Server struct {
	RunMode      string `json:"RunMode"`
	Host         string `json:"Host"`
	Port         int    `json:"Port"`
	ReadTimeout  int    `json:"ReadTimeout"`
	WriteTimeout int    `json:"WriteTimeout"`
}

type DBCfg struct {
	Host           string `json:"Host"`
	Port           int    `json:"Port"`
	User           string `json:"User"`
	Password       string `json:"Password"`
	DBName         string `json:"DBName"`
	MaxIdleConns   int    `json:"MaxIdleConns"`
	MaxOpenConns   int    `json:"MaxOpenConns"`
	MigrationsPath string `json:"MigrationsPath"`
}

type Jwt struct {
	Secret string `json:"Secret"`
	Exp    int    `json:"Exp"`
}

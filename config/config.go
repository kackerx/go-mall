package config

var (
	Conf *Config
)

type Config struct {
	App *App
	DB  *DB
}

type App struct {
	Name       string `mapstructure:"name"`
	Env        string `mapstructure:"env"`
	Log        *Log
	Pagination *Pagination
}

type Pagination struct {
	DefaultSize int `mapstructure:"default_size"`
	MaxSize     int `mapstructure:"max_size"`
}

type Log struct {
	Path    string `mapstructure:"path"`
	MaxSize int    `mapstructure:"max_size"`
	MaxAge  int    `mapstructure:"max_age"`
}

type DB struct {
	Type        string `mapstructure:"type"`
	Dsn         string `mapstructure:"dsn"`
	MaxOpen     int    `mapstructure:"max_open"`
	MaxIdle     int    `mapstructure:"max_idle"`
	MaxLiftTime int    `mapstructure:"max_lift_time"`
}

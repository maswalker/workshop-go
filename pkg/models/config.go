package models

type LogConfig struct {
	Mode     string `json:",default=console,options=console|file"`
	Level    int    `json:",default=0,range=[0:5]"`
	Path     string `json:",default=./"`
	Compress bool   `json:",default=true"`
}

type Database struct {
}

type WebConfig struct {
	Port uint32 `json:",default=8080,range=[3000:65535]"`
}

type ServiceConfig struct {
	Mode string `json:",default=pro,options=dev|pro|test"`
}

type Config struct {
	Service ServiceConfig
	Log     LogConfig
	Web     WebConfig
}

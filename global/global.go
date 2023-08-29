package global

import (
	"file-push/config"
	"github.com/spf13/viper"
)

var (
	GVA_VP     *viper.Viper
	GVA_CONFIG config.Server
)

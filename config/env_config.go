package config

import(
	"os"
	"log"
	"github.com/spf13/viper"

)
func GetConfig()( *viper.Viper,error){
	path,err := os.Getwd()
	
	if err != nil {
		log.Println("read path error")
	}
	confi_path := path[:len(path)-4] +"/config"; 
	
    config := viper.New()
	config.AddConfigPath(confi_path) 
	config.SetConfigType("yaml")
    config.SetConfigName ("dev-env")
	if err := config.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            // Config file not found; ignore error if desired
            log.Println("no such config file")
        } else {
            // Config file was found but another error was produced
            log.Println("read config error")
        }
        //log.Fatal(err) // failed to read configuration file. Fatal error
	}
	return config, err
}

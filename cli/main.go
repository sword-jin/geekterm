package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/rivo/tview"
	geekhub "github.com/rrylee/geekterm"
	"github.com/spf13/viper"
)

var (
	geekhubDir string
	app        *tview.Application // The tview application.
)

var (
	configFile = flag.String("config-file", "", "yaml config file")
	cookie     = flag.String("cookie", "", "geekhub cookie")
)

// Main entry point.
func main() {
	flag.Parse()

	v := viper.New()
	if *configFile != "" {
		v.SetConfigFile(*configFile)
		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}
	}

	cfg := initConfig(v)

	//Start the application.
	app = tview.NewApplication()

	geekhub.Setup(cfg)
	geekhub.Draw(app)
	geekhub.Keybinds(app)

	if err := app.Run(); err != nil {
		fmt.Printf("Error running application: %s\n", err)
	}
}

func initConfig(v *viper.Viper) *geekhub.Config {
	geekhubDir = getUserHomeDir() + "/.geekhub/"
	if _, err := os.Stat(geekhubDir); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(geekhubDir, 0644)
			if err != nil {
				panic(err)
			}
		}
	}

	cfg := &geekhub.Config{}
	if v.GetString("log-file") != "" {
		cfg.LogFile = v.GetString("log-file")
	} else {
		cfg.LogFile = geekhubDir + "log.txt"
	}

	cfg.LogLevel = v.GetInt("log-level")

	if *cookie != "" {
		cfg.Cookie = *cookie
	} else if v.GetString("cookie") != "" {
		cfg.Cookie = v.GetString("cookie")
	}

	return cfg
}

func getUserHomeDir() string {
	var home string
	switch runtime.GOOS {
	case "windows":
		home, _ = os.LookupEnv("LOCALAPPDATA")
	case "linux":
		home, _ = os.LookupEnv("HOME")
		break
	case "darwin":
		home, _ = os.LookupEnv("HOME")
		break
	}
	return home
}

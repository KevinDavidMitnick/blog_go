package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/ulricqin/blog/config"
	"github.com/ulricqin/blog/http"
	"github.com/ulricqin/blog/model"
)

var (
	cfg     *string
	version *bool
	help    *bool
)

// set log flag LstdFlags | Lshortfile
func prepare() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// parse command line parms
func handle_cmd() {
	cfg = flag.String("c", "cfg.json", "configuration file")
	version = flag.Bool("v", false, "show version")
	help = flag.Bool("h", false, "help")
	flag.Parse()
}

// read config files
func read_configs() {
	handle_version(*version)
	handle_help(*help)
	handle_config(*cfg)
}

// handle version
func handle_version(displayVersion bool) {
	if displayVersion {
		fmt.Println(config.AUTHOR, config.EMAIL, config.VERSION)
		os.Exit(0)
	}
}

// handle help
func handle_help(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}

// handle config
func handle_config(configFile string) {
	err := config.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	// prepare log output
	prepare()

	// handle cmd params
	handle_cmd()

	// read config files
	read_configs()

}

func main() {
	model.InitMysql()
	model.InitCache()
	http.Start()
}

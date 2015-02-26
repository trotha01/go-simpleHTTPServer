package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var chttp = http.NewServeMux()

type config struct {
	host     string
	port     string
	log_dir  string
	log_file string
}

type CheckListItem struct {
	item    string
	checked int
}

type CheckList []CheckListItem

var conf *config
var logger *log.Logger
var logfile *os.File

func init() {
	configFromEnv()
	initializeLogger()
}

func main() {
	chttp.Handle("/", http.FileServer(http.Dir("./")))

	http.HandleFunc("/", HomeHandler) // homepage

	bind := fmt.Sprintf("%s:%s", conf.host, conf.port)
	logger.Printf("Listening on %s...", bind)

	err := http.ListenAndServe(bind, nil)
	if err != nil {
		logger.Printf("Error serving: %s", err.Error())
		cleanup()
		os.Exit(1)
	}
}

func configFromEnv() {
	conf = new(config)
	conf.host = os.Getenv("HOST")
	conf.port = os.Getenv("PORT")
	conf.log_dir = os.Getenv("LOG_DIR")
	conf.log_file = os.Getenv("LOG_FILE")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("Request /")
	chttp.ServeHTTP(w, r)
}

func initializeLogger() {
	logfile, err := os.OpenFile(conf.log_dir+"/"+conf.log_file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}

	logger = log.New(logfile, "", log.Ldate|log.Ltime|log.Lshortfile)
	if logger == nil {
		log.Printf("Could not create logger\n")
	}
}

func cleanup() {
	logfile.Close()
}

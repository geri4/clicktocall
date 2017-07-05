package main

import (
	"fmt"
	"net/http"
	"os"
	log "github.com/Sirupsen/logrus"
	"github.com/warik/gami"
	"regexp"
	"errors"
)

type Call struct {
	phone1 string
	phone2 string
}

type Config struct {
	amiHost string
	amiLogin string
	amiPassword string
	channel string
	context string
	token string
}

func placeCall(callchan chan Call, config Config) {
	for {
		aster := gami.NewAsterisk( config.amiHost, config.amiLogin, config.amiPassword)
		err := aster.Start()
		if err != nil {
			log.Fatal(err)
		}
		for {
			newcall := <-callchan
			log.WithFields(log.Fields{
				"phone1": newcall.phone1,
				"phone2": newcall.phone2,
			}).Info("Place new call")
			orignate := gami.NewOriginate( config.channel+"/"+newcall.phone1, config.context, newcall.phone2, "1")
			err = aster.Originate(orignate, nil, nil)
			if err != nil {
				log.Error("Call is failed... Trying to reconnect")
				_ = aster.Logoff()
				break
			}
		}
	}
	//_ = aster.Logoff()
}

func FormatPhone(phone *string) error {
	re := regexp.MustCompile("\\+7")
	*phone = re.ReplaceAllString(*phone, "8")
	re = regexp.MustCompile("[0-9]+")
	*phone = re.FindString(*phone)
	re = regexp.MustCompile("^7")
	*phone = re.ReplaceAllString(*phone, "8")
	if len(*phone) > 0 || len(*phone) < 16 {
		return nil
	} else {
		return errors.New("Phone is invalid")
	}
}

func parseRequest(w http.ResponseWriter, r *http.Request, callchan chan Call, token string) {
	fmt.Println("GET params were:", r.URL.Query())
	phone1 := r.URL.Query().Get("phone1")
	phone2 := r.URL.Query().Get("phone2")
	receivedToken := r.URL.Query().Get("token")
	fmt.Println("rtoken:", receivedToken)
	if receivedToken != token {
		fmt.Fprintf(w, "tokenhueken")
		returnStatus(w, http.StatusForbidden, "{ \"status\": \"Invalid token\" }")
		return
	}
	err := FormatPhone(&phone1)
	if err != nil {
		returnStatus(w, http.StatusBadRequest, "{ \"status\": \"Invalid phone1 value\" }")
		return
	}
	err = FormatPhone(&phone2)
	if err != nil {
		returnStatus(w, http.StatusBadRequest, "{ \"status\": \"Invalid phone2 value\" }")
		return
	}
	//fmt.Fprintf(w, "Calling...") // send data to client side
	returnStatus(w, http.StatusOK, "{ \"status\": \"calling\" }")
	callchan <- Call{phone1, phone2}
	//placeCall(phone1, phone2)
}

func returnStatus (w http.ResponseWriter, httpstatus int, response string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpstatus)
	w.Write([]byte(response))
}

func getEnv() Config {
	var config Config
	config.amiHost = os.Getenv("AMIHOST")
	if config.amiHost == "" {
		log.Panic("Env AMIHOST is undefined")
	}
	config.amiLogin = os.Getenv("AMILOGIN")
	if config.amiLogin == "" {
		log.Panic("Env AMILOGIN is undefined")
	}
	config.amiPassword = os.Getenv("AMIPASSWORD")
	if config.amiPassword == "" {
		log.Panic("Env AMIPASSWORD is undefined")
	}
	config.channel = os.Getenv("CHANNEL")
	if config.channel == "" {
		log.Panic("Env CHANNEL is undefined")
	}
	config.context = os.Getenv("CONTEXT")
	if config.context == "" {
		log.Panic("Env CONTEXT is undefined")
	}
	config.token = os.Getenv("TOKEN")
	if config.context == "" {
		log.Panic("Env TOKEN is undefined")
	}
	return config
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	config := getEnv()
	callchan := make(chan Call)
	go placeCall(callchan, config)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		parseRequest(w, r, callchan, config.token)
	}) // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

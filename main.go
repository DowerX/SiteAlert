package main

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"./config"
	"./errorcheck"
)

var c config.Config
var h [20]byte
var conf string

var auth smtp.Auth

func main() {
	flag.StringVar(&conf, "config", "./config.yml", "Set the config file.")
	flag.Parse()
	c = config.GetConfig(conf)
	checkflags()
	if c.Log {
		lf, err := os.OpenFile(c.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("ERR: Can't log!")
			return
		}
		log.SetOutput(lf)
	}
	auth = smtp.PlainAuth("", c.Sender, c.Password, c.Server)
	for {
		check()
		time.Sleep(time.Duration(c.Time) * time.Minute)
	}
}

func checkflags() {
	flag.StringVar(&c.Url, "url", c.Url, "Set URL.")
	flag.StringVar(&c.Email, "email", c.Email, "Set email to alert.")
	flag.StringVar(&c.Sender, "sender", c.Sender, "Set the sende's email.")
	flag.StringVar(&c.Password, "password", c.Password, "Set the password.")
	flag.StringVar(&c.Server, "server", c.Server, "Set SMTP server.")
	flag.StringVar(&c.Port, "port", c.Port, "Set the SMTP servers port.")
	flag.IntVar(&c.Time, "time", c.Time, "Set the time between checks.")
	flag.BoolVar(&c.Log, "log", c.Log, "Enable logging to a file.")
	flag.StringVar(&c.Logfile, "logfile", c.Logfile, "Set logfile.")
	flag.Parse()
}

func writelog(l ...interface{}) {
	fmt.Println(l, "")
	if c.Log {
		log.Println(l, "")
	}
}

func check() {
	resp, err := http.Get(c.Url)
	errorcheck.Check(err)
	b, err := ioutil.ReadAll(resp.Body)
	errorcheck.Check(err)
	nh := sha1.Sum(b)
	if h != nh {
		h = nh
		alert()
	}
}

func alert() {
	to := []string{c.Email}
	msg := []byte(c.Msg)
	fmt.Println(to, msg)
	err := smtp.SendMail(c.Server+c.Port, auth, c.Sender, to, msg)
	errorcheck.Check(err)
	writelog("Alerted.")
}

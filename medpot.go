package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-ini/ini"
	"github.com/mozillazg/request"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	CONN_HOST = "127.0.0.1"
	CONN_PORT = "2575"
	CONN_TYPE = "tcp"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
	read config from EWS poster for DTAGs Early warning system and T-Pot
*/
func readConfig() (string, string, string, string) {

	cfg, err := ini.Load("/etc/ews.cfg")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	target := cfg.Section("EWS").Key("rhost_first").String()
	user := cfg.Section("EWS").Key("username").String()
	password := cfg.Section("EWS").Key("password").String()
	nodeid := cfg.Section("GLASTOPFV3").Key("nodeid").String()
	nodeid = strings.Replace(nodeid, "glastopfv3-", "medpot-", -1)
	return target, user, password, nodeid

}

func post(target string, user string, password string, nodeid string) {

	c := &http.Client{}
	req := request.NewRequest(c)

	dat := readFile("ews.xml")
	body := strings.Replace(string(dat), "_USERNAME_", user, -1)
	body = strings.Replace(body, "_TOKEN_", password, -1)
	body = strings.Replace(body, "_NODEID_", nodeid, -1)

	// not set Content-Type
	req.Body = strings.NewReader(string(body))
	resp, err := req.Post("http://127.0.0.1:9922/ews-0.1/alert/postSimpleMessage")

	if err != nil {
		fmt.Println("Error http post:", err.Error())
	} else {
		fmt.Println("Http Reponse", resp.Status)
	}

}

func initLogger() *zap.Logger {

	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/var/log/medpot.log"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}

func main() {

	fmt.Print("Starting Medpot at ")
	currentTime := time.Now()
	fmt.Println(currentTime.Format("2006.01.02 15:04:05"))

	logger := initLogger()

	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()

	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		// Handle connections in a new goroutine.
		go handleRequest(conn, logger)
	}
}

/*
	reads file from both possible locations (first repo location, second location from docker install
 */
func readFile(name string) []byte {

	b1 := make([]byte, 4)

	dat, err := ioutil.ReadFile("./template/" + name)
	if (err == nil) {
		return dat
	}

	dat, err = ioutil.ReadFile("/data/medpot/" + name)
	if (err == nil) {
		return dat

	}

	return b1

}

// Handles incoming requests.
func handleRequest(conn net.Conn, logger *zap.Logger) {
	// Make a buffer to hold incoming data.

	buf := make([]byte, 1024*1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	remote := fmt.Sprintf("%s", conn.RemoteAddr())
	ip, port, _ := net.SplitHostPort(remote)

	currentTime := time.Now()
	fmt.Print(currentTime.Format("2006.01.02 15:04:05"))
	myTime := currentTime.Format("2006.01.02 15:04:05")

	fmt.Print(": Connecting from ip ", ip)
	fmt.Println(" and port ", port)

	dat := readFile("dummyerror.xml")

	// Send a response back to person contacting us.
	conn.Write(dat)

	// copy to a real buffer
	bufTarget := make([]byte, reqLen)
	copy(bufTarget, buf)

	spew.Dump(bufTarget)

	encoded := base64.StdEncoding.EncodeToString([]byte(bufTarget))

	logger.Info("Connection found",
		// Structured context as strongly typed Field values.
		zap.String("time", myTime),
		zap.String("port", port),
		zap.String("ip", ip),
		zap.String("data", encoded),
	)

	// Close the connection when you're done with it.
	conn.Close()
}

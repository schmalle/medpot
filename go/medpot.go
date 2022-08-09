package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-ini/ini"
	"github.com/mozillazg/request"
	"go.uber.org/zap"

	arg "github.com/s9rA16Bf4/ArgumentParser/go/arguments"
	"github.com/s9rA16Bf4/notify_handler/go/notify"
)

type conf_t struct {
	target       string
	username     string
	password     string
	nodeID       string
	ip           string      // Source ip
	port         string      // Source port
	encoded      string      // The encoded (base64) value of what we recieved from the one connecting to us
	conn         net.Conn    // The active connection
	log_location string      // Where on the disk is the log file located
	logger       *zap.Logger // Object to our logger
	ews          string      // Is true if we are gonna post to EWS else false
}

const (
	CONN_HOST       = "0.0.0.0"
	CONN_TYPE       = "tcp"
	VERSION         = "1.2" // Current version
	CONFIG_LOCATION = "/etc/medpot"
)

/*
	read config from EWS poster for DTAGs Early warning system and T-Pot
*/
func readConfig() (string, string, string, string, string) {

	cfg, err := ini.Load(fmt.Sprintf("%s/ews.cfg", CONFIG_LOCATION))
	if err != nil {
		notify.Error(err.Error(), "medpot.readConfig()")
	}

	target := cfg.Section("EWS").Key("rhost_first").String()
	user := cfg.Section("EWS").Key("username").String()
	password := cfg.Section("EWS").Key("token").String()
	nodeid := cfg.Section("GLASTOPFV3").Key("nodeid").String()
	nodeid = strings.Replace(nodeid, "glastopfv3-", "medpot-", -1)
	ews := cfg.Section("EWS").Key("ews").String()

	return target, user, password, nodeid, ews

}

func post(cconf_t *conf_t, time string) {

	if cconf_t.ews == "false" { // Should we post this data?
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	c := &http.Client{Transport: tr}
	req := request.NewRequest(c)

	dat := readFile("ews.xml")
	body := strings.Replace(string(dat), "_USERNAME_", cconf_t.username, -1)
	body = strings.Replace(body, "_TOKEN_", cconf_t.password, -1)
	body = strings.Replace(body, "_NODEID_", cconf_t.nodeID, -1)
	body = strings.Replace(body, "_IP_", cconf_t.ip, -1)
	body = strings.Replace(body, "_PORT_", cconf_t.port, -1)
	body = strings.Replace(body, "_TIME_", time, -1)
	body = strings.Replace(body, "_DATA_", cconf_t.encoded, -1)

	// not set Content-Type
	req.Body = strings.NewReader(string(body))
	resp, err := req.Post(cconf_t.target)

	if err != nil {
		notify.Warning(fmt.Sprintf("Failed to post, error message attached %s", err.Error()))
	} else {
		notify.Inform(fmt.Sprintf("Http reponse status: %s", resp.Status))
	}

}

func initLogger(cconf_t *conf_t) *zap.Logger {

	rawJSON := []byte(fmt.Sprintf(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "%s"],
	  "errorOutputPaths": ["stderr"],
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`, cconf_t.log_location))

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		notify.Error(err.Error(), "medpot.initLogger()")
	}

	logger, err := cfg.Build()
	if err != nil {
		notify.Error(err.Error(), "medpot.initLogger()")
	}

	defer logger.Sync()

	return logger
}

func main() {
	arg.Argument_add_with_options("--set_logo", "-sl", true, "Allows you to pick a logo that is shown on boot", []string{"1", "2"})
	arg.Argument_add("--set_port", "-sp", true, "Allows for a different port to be used, default = 2575")
	arg.Argument_add("--set_log_location", "-sll", true, "Changes the directory where the logs will be placed, default = '/var/log/medpot/'")

	parsed_flags := arg.Argument_parse() // Checks which arguments that can have been passed onto the program

	cconf_t := new(conf_t)
	current_logo := LOGO_2

	if value, entered := parsed_flags["-sl"]; entered {
		if value == "1" {
			current_logo = LOGO_1
		} else {
			current_logo = LOGO_2
		}
	}

	if value, entered := parsed_flags["-sp"]; entered {
		cconf_t.port = value
		_, err := strconv.Atoi(cconf_t.port) // Checks if it's a valid port
		if err != nil {
			notify.Error(err.Error(), "medpot.main()")
		}
	} else {
		cconf_t.port = "2575"
	}

	if value, entered := parsed_flags["-sll"]; entered {
		cconf_t.log_location = value
	} else {
		cconf_t.log_location = "/var/log/medpot/"
	}
	cconf_t.log_location += "medpot.log"

	fmt.Println(current_logo) // Print the logo that will be used
	notify.Inform(fmt.Sprintf("V.%s", VERSION))
	notify.Inform(fmt.Sprintf("Starting Medpot at %s", time.Now().Format(time.RFC822)))
	notify.Inform("Written by @schmalle, forked and updated by @s9rA16Bf4")
	notify.Inform("If you find any bugs, report them on 'github.com/s9rA16Bf4/medpot' or 'github.com/schmalle/medpot'")
	notify.Inform("--------------------------------------------------------")
	notify.Inform(fmt.Sprintf("Log files will be located at '%s'", cconf_t.log_location))
	notify.Inform(fmt.Sprintf("Will utilize port %s", cconf_t.port))

	cconf_t.target, cconf_t.username, cconf_t.password, cconf_t.nodeID, cconf_t.ews = readConfig()

	cconf_t.logger = initLogger(cconf_t)

	l, err := net.Listen(CONN_TYPE, fmt.Sprintf(":%s", cconf_t.port))

	if err != nil {
		notify.Error(err.Error(), "medpot.main()")
	}
	// Close the listener when the application closes.
	defer l.Close()

	notify.Inform(fmt.Sprintf("Listening on host %s on port %s", CONN_HOST, cconf_t.port))

	for {
		// Listen for an incoming connection.
		cconf_t.conn, err = l.Accept()
		if err != nil {
			notify.Error(err.Error(), "medpot.main()")
		}

		// Handle connections in a new goroutine.
		go handleRequest(cconf_t)
	}
}

/*
	reads file from both possible locations (first repo location, second location from docker install
*/
func readFile(name string) []byte {

	b1 := make([]byte, 4)

	dat, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", CONFIG_LOCATION, name))
	if err != nil {
		notify.Warning(fmt.Sprintf("Failed to read file '%s'", name))
		return b1
	}

	return dat
}

func handleClientRequest(cconf_t *conf_t, buf []byte, reqLen int) {
	// These templates are utilized when a user connects
	dat := readFile("dummyerror.xml")

	// copy to a real buffer
	bufTarget := make([]byte, reqLen)
	copy(bufTarget, buf)

	if strings.Contains(string(buf), "MSH") && strings.Index(string(buf), "MSH|") == 0 {
		dat = readFile("dummyok.xml")
	}

	// Send a response back to person contacting us.
	cconf_t.conn.Write(dat)

}

// Handles incoming requests.
func handleRequest(cconf_t *conf_t) {
	// Make a buffer to hold incoming data.

	buf := make([]byte, 1024^2)
	counter := 0

	for counter < 3 {

		timeoutDuration := 3 * time.Second // time out in three seconds from now
		cconf_t.conn.SetReadDeadline(time.Now().Add(timeoutDuration))

		// Read the incoming connection into the buffer.
		reqLen, err := cconf_t.conn.Read(buf)
		if err != nil {
			cconf_t.conn.Close()

			if err.Error() != "EOF" {
				notify.Inform(err.Error()) // Most likely a time out
			}
			break

		} else {

			remote := cconf_t.conn.RemoteAddr().String()
			cconf_t.ip, cconf_t.port, _ = net.SplitHostPort(remote)
			currentTime := time.Now().Format(time.RFC822)

			notify.Inform(fmt.Sprintf("Connection from '%s' on port '%s' at time %s", cconf_t.ip, cconf_t.port, currentTime))

			handleClientRequest(cconf_t, buf, reqLen)

			// copy to a real buffer
			bufTarget := make([]byte, reqLen)
			copy(bufTarget, buf)

			spew.Dump(bufTarget)

			cconf_t.encoded = base64.StdEncoding.EncodeToString([]byte(bufTarget))

			cconf_t.logger.Info("Connection found",
				// Structured context as strongly typed Field values.
				zap.String("timestamp", currentTime),
				zap.String("src_port", cconf_t.port),
				zap.String("src_ip", cconf_t.ip),
				zap.String("data", cconf_t.encoded),
			)

			// if configured, send back data to PEBA / DTAG T_pot homebase
			post(cconf_t, currentTime)

		}

		counter++
		notify.Inform(fmt.Sprintf("Increased counter to %s", fmt.Sprint(counter)))
	}
	notify.Warning("Maximum loop counter reached... loop will now end!")
	cconf_t.conn.Close()
}

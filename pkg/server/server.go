package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/pheelee/dyndns/pkg/config"
	"github.com/pheelee/dyndns/pkg/dns"
	"github.com/pheelee/dyndns/pkg/logger"
)

var (
	appConfig     *config.Config
	updateTracker lastUpdate
)

type txtRecordRequest struct {
	FQDN  string `json:"fqdn"`
	VALUE string `json:"value"`
}

type aRecordRequest struct {
	domain   string
	username string
	password string
	ip4addr  string
}

type lastUpdate struct {
	Domain    string
	Timestamp time.Time
}

func (r *aRecordRequest) fromURL(query string) error {
	vals, err := url.ParseQuery(query)
	if err != nil {
		return err
	}
	r.domain = vals.Get("domain")
	r.username = vals.Get("username")
	r.password = vals.Get("password")
	r.ip4addr = vals.Get("ip4addr")
	return nil
}

func (r *aRecordRequest) valid() bool {
	if b, _ := regexp.MatchString(`^([a-zA-Z0-9]+\.){2,63}[a-zA-Z]{2,6}$`, r.domain); b != true {
		logger.Warning("Domain does not match regex")
		return false
	}
	if b, _ := regexp.MatchString(`^[a-zA-Z0-9]{3,32}$`, r.username); b != true {
		logger.Warning("Username does not match regex")
		return false
	}
	if b, _ := regexp.MatchString(`^.{3,64}$`, r.password); b != true {
		logger.Warning("Password does not match regex")
		return false
	}
	if b, _ := regexp.MatchString(`^(\d{1,3}\.){3}\d{1,3}$`, r.ip4addr); b != true {
		logger.Warning("ip4addr does not match regex")
		return false
	}
	return true
}

func txtHandler(w http.ResponseWriter, r *http.Request) {
	/*if authValid(r) == false {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}*/
	var data txtRecordRequest
	d := json.NewDecoder(r.Body)
	if d.Decode(&data) != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var (
		rec dns.Record
		res dns.Result
	)
	rec = dns.Record{
		Hostname: data.FQDN,
		Type:     "TXT",
		TTL:      "600",
		Data:     data.VALUE,
	}

	// if same domain is sent again within 5 seconds it must be a request for a wildcard domain e.g
	// Domain: example.com
	// SAN: *.example.com
	now := time.Now()
	if updateTracker.Domain == rec.Hostname && updateTracker.Timestamp.Add(time.Second*5).After(now) {
		logger.Warning("update for same domain within 5s detected, skipping")
		return
	}
	updateTracker = lastUpdate{
		Domain:    rec.Hostname,
		Timestamp: now,
	}

	switch r.URL.Path {
	case "/present":
		res = appConfig.DNSServer.Active.Update(rec)
	case "/cleanup":
		res = appConfig.DNSServer.Active.Delete(rec)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if res.Error != nil {
		logger.Error(res.Error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(res.Message))
		return
	}

	logger.Info(res.Message)
}

func dyndnsUpdate(w http.ResponseWriter, r *http.Request) {
	/*
		Produces http 200 with body:
		- 911 			fatal error
		- nohost		domain not found
		- badauth		invalid credentials
		- good			all ok
	*/

	// Validate user input
	var data aRecordRequest
	err := data.fromURL(r.URL.RawQuery)
	if err != nil || data.valid() == false {
		fmt.Fprint(w, "911")
		return
	}

	// find domain and validate authorization
	domain, ok := appConfig.Domains[data.domain]
	if !ok {
		fmt.Fprint(w, "nohost")
		return
	}

	if domain.ValidateCredentials(data.username, data.password) == false {
		fmt.Fprint(w, "badauth")
		return
	}

	// good to go update record
	ret := appConfig.DNSServer.Active.Update(dns.Record{
		Hostname: data.domain,
		TTL:      "300",
		Type:     "A",
		Data:     data.ip4addr,
	})

	if ret.Error == nil {
		fmt.Fprint(w, "good")
	}
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		srcIP := r.RemoteAddr
		if appConfig.Global.RealIPHeader != "" {
			srcIP = r.Header.Get(appConfig.Global.RealIPHeader)
		}
		logger.Info(fmt.Sprintf("%s %s %s %dms", r.Method, r.URL.Path, srcIP, time.Now().Sub(t1).Milliseconds()))
	})
}

// SetupServer maps the routes to their functions
func SetupServer(c *config.Config) {
	appConfig = c

	http.Handle("/present", logRequest(basicAuth(http.HandlerFunc(txtHandler))))
	http.Handle("/cleanup", logRequest(basicAuth(http.HandlerFunc(txtHandler))))
	http.HandleFunc("/update", dyndnsUpdate)

}

// Helper Functions
/*
func authValid(r *http.Request) bool {
	if user, pass, ok := r.BasicAuth(); ok == true {
		if user == appConfig.HTTPReqAuth.Username && pass == appConfig.HTTPReqAuth.Password {
			return true
		}
	}
	return false
}
*/

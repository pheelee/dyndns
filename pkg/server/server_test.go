package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pheelee/dyndns/pkg/config"
)

func loadConfig() {
	appConfig = &config.Config{}
	appConfig.New()
	appConfig.DNSServer.Bind.Debug = true
	appConfig.DNSServer.Active = &appConfig.DNSServer.Bind
}

func getResponse(url string, method string, function func(w http.ResponseWriter, r *http.Request)) (string, *httptest.ResponseRecorder) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", nil
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(function)

	handler.ServeHTTP(rr, req)
	body, _ := ioutil.ReadAll(rr.Body)
	return string(body), rr
}

func TestDyndnsBadAuth(t *testing.T) {
	loadConfig()
	body, _ := getResponse("/update?domain=dyn.example.com&username=testuser&password=testpass&ip4addr=1.2.3.4", "GET", dyndnsUpdate)
	if body != "badauth" {
		t.Errorf("DynDNS update request expected \"%s\" got \"%s\"", "badauth", string(body))
	}
}

func TestDyndnsGood(t *testing.T) {
	loadConfig()
	body, _ := getResponse("/update?domain=dyn.example.com&username=dynuser&password=dynpass&ip4addr=1.2.3.4", "GET", dyndnsUpdate)
	if body != "good" {
		t.Errorf("DynDNS update request expected \"%s\" got \"%s\"", "good", string(body))
	}
}

func TestDyndnsNoHost(t *testing.T) {
	loadConfig()
	body, _ := getResponse("/update?domain=dyn2.example.com&username=dynuser&password=dynpass&ip4addr=1.2.3.4", "GET", dyndnsUpdate)
	if body != "nohost" {
		t.Errorf("DynDNS update request expected \"%s\" got \"%s\"", "nohost", string(body))
	}
}

package main

import (
	"flag"
	"github.com/danieldin95/lightstar/http/api"
	"github.com/danieldin95/lightstar/libstar"
	"github.com/danieldin95/lightstar/proxy"
	"strings"
)

func GetPorts(url, auth string) []string {
	ports := make([]string, 0, 32)

	client := libstar.HttpClient{
		Auth: libstar.Auth{
			Type:     "basic",
			Username: auth,
		},
		Url: url + "/api/proxy/tcp",
	}
	r, err := client.Do()
	if err == nil {
		api.GetJSON(r.Body, &ports)
	} else {
		libstar.Error("main %s", err)
	}
	libstar.Info("main %s", ports)
	return ports
}

func main() {
	tgt := ""
	file := "lightprix.json"
	cfg := struct {
		Url     string   `json:"url"`
		Auth    string   `json:"Auth"`
		Target  []string `json:"target"`
		Verbose int      `json:"log.verbose"`
		LogFile string   `json:"log.file"`
	}{
		Url:     "https://localhost:10080",
		Auth:    "admin:123456",
		Verbose: 2,
		LogFile: "lightprix.log",
	}

	flag.StringVar(&file, "conf", file, "the configuration file")
	flag.StringVar(&cfg.Url, "url", cfg.Url, "the url path.")
	flag.StringVar(&tgt, "tgt", tgt, "the target proxied to.")
	flag.StringVar(&cfg.Auth, "auth", cfg.Auth, "the auth login to.")
	flag.IntVar(&cfg.Verbose, "log:level", cfg.Verbose, "logger level")
	flag.Parse()

	cfg.Target = strings.Split(tgt, ",")
	if err := libstar.JSON.UnmarshalLoad(&cfg, file); err != nil {
		libstar.Warn("main %s", err)
	}

	libstar.Init(cfg.LogFile, cfg.Verbose)
	ports := GetPorts(cfg.Url, cfg.Auth)
	for _, port := range ports {
		cfg.Target = append(cfg.Target, port)
	}
	pri := proxy.Proxy{
		Target: cfg.Target,
		Client: &proxy.WsClient{
			Auth: libstar.Auth{
				Type:     "basic",
				Username: cfg.Auth,
			},
			Url: cfg.Url + "/ext/tcpsocket?target=",
		},
	}
	pri.Initialize().Start()
	defer pri.Stop()
	libstar.Wait()
}

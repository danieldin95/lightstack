package proxy

import (
	"crypto/tls"
	"github.com/danieldin95/lightstar/libstar"
	"golang.org/x/net/websocket"
	"net/http"
)

type WsClient struct {
	Auth      libstar.Auth
	Url       string
	TlsConfig *tls.Config
	Protocol  string
}

func (w *WsClient) Initialize() {
	w.TlsConfig = &tls.Config{InsecureSkipVerify: true}
}

func (w *WsClient) Dial() (ws *websocket.Conn, err error) {
	config, err := websocket.NewConfig(w.Url, w.Url)
	if err != nil {
		return nil, err
	}
	if w.Protocol != "" {
		config.Protocol = []string{w.Protocol}
	}
	config.TlsConfig = w.TlsConfig
	if w.Auth.Type == "basic" {
		config.Header = http.Header{
			"Authorization": {libstar.BasicAuth(w.Auth.Username, w.Auth.Password)},
		}
	}
	return websocket.DialConfig(config)
}

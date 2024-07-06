package testutils

import (
	"net/http"
	"net/url"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/path"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/utils"
	"github.com/gorilla/websocket"
)

var (
	QueryRootTest = test.APITestCase[[]entity.Cheptel]{
		Method: utils.String("GET"),
		URL: &url.URL{
			Scheme: "ws",
			Path:   path.Query,
		},
		WSDialer:   websocket.DefaultDialer,
		WantStatus: utils.Int(http.StatusSwitchingProtocols),
	}
)

package testutils

import (
	"net/http"

	"github.com/gaetanDubuc/beeckend/internal/cheptel/path"
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/internal/test"
	"github.com/gaetanDubuc/beeckend/pkg/pagination"
	"github.com/gaetanDubuc/beeckend/pkg/utils"
)

var (
	QueryRootTest = test.APITestCase[pagination.Pages[[]entity.Cheptel]]{
		Method:     utils.String("GET"),
		URL:        utils.String(path.Query),
		WantStatus: utils.Int(http.StatusOK),
	}
)

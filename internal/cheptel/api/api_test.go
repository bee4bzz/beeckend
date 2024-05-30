package api

import (
	"github.com/gaetanDubuc/beeckend/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *APITestSuite) SetupSuite() {
	suite.router, _ = router.New()
}

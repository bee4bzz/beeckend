package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type APITestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *APITestSuite) SetupSuite() {
	suite.router = gin.Default()
}

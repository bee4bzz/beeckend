package middleware

import (
	"bytes"

	"github.com/gaetanDubuc/beeckend/internal/authenticator/testutils"
	"github.com/gaetanDubuc/beeckend/internal/cheptel/testutils"
	"github.com/gaetanDubuc/beeckend/pkg/log"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
)

type MiddlewareTestSuite struct {
	suite.Suite

	logger   *log.Logger
	observer *observer.ObservedLogs
	buffer   *bytes.Buffer

	UserRepository *testutils.UserRepository
	middleware     *Middleware
}

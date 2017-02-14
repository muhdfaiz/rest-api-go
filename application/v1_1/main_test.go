package v1_1

import (
	"net/http/httptest"
	"testing"

	"bitbucket.org/cliqers/shoppermate-api/test/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var (
	Router     = gin.Default()
	TestHelper = helper.Helper{}
	TestServer *httptest.Server
	DB         *gorm.DB
)

func TestMain(m *testing.M) {
	TestHelper.Setup()

	DB = Database.Connect("test")

	TestServer = httptest.NewServer(InitializeObjectAndSetRoutesV1_1(Router, DB))

	ret := m.Run()

	if ret == 0 {
		//TestHelper.Teardown()
		TestServer.Close()
	}
}

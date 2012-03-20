package service_test

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/timeredbull/tsuru/api/service"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type ServiceSuite struct{}

var _ = Suite(&ServiceSuite{})
var db, _ = sql.Open("sqlite3", "./tsuru.db")

func (s *ServiceSuite) SetUpSuite(c *C) {
	_, err := db.Exec("CREATE TABLE 'service_bindings' ('id' INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 'service_config_id' integer, 'app_id' integer, 'user_id' integer, 'binding_token_id' integer, 'name' varchar(255), 'configuration' text, 'credentials' text, 'binding_options' text, 'created_at' datetime, 'updated_at' datetime)")
	c.Check(err, IsNil)
}

func (s *ServiceSuite) TearDownSuite(c *C) {
	db.Close()
	os.Remove("./tsuru.db")
}

func (s *ServiceSuite) TearDownTest(c *C) {
	db.Exec("DELETE FROM service_bindings")
}

func (s *ServiceSuite) TestShouldRequestCreateAndBeSuccess(c *C) {
	request, err := http.NewRequest("POST", "services/create", nil)
	recorder := httptest.NewRecorder()
	c.Assert(err, IsNil)

	service.CreateServiceHandler(recorder, request)
	status := recorder.Code

	c.Assert(200, Equals, status)
}

func (s *ServiceSuite) TestShouldRequestCreateAndInsertInTheDatabase(c *C) {
	request, err := http.NewRequest("POST", "services/create", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Form = url.Values{
		"serviceBindingId": []string{"1"},
		"appId":            []string{"1"},
		"userId":           []string{"1"},
		"bindingToken":     []string{"2"},
		"name":             []string{"my_mysql"},
	}

	recorder := httptest.NewRecorder()
	c.Assert(err, IsNil)

	service.CreateServiceHandler(recorder, request)
	body := recorder.Body
	c.Assert(body.String(), Equals, "success")

	rows, err := db.Query("SELECT count(*) FROM service_bindings WHERE name = 'my_mysql'")

	c.Check(err, IsNil)
	var qtd int

	for rows.Next() {
		rows.Scan(&qtd)
	}

	c.Assert(1, Equals, qtd)
}

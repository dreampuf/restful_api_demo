package web_v1

import (
	"db"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"reflect"
	"errors"
	"runtime/debug"
)

type MockDataSource struct {
	userDS []db.User
	rsDS []db.RelationShip
}

type BaseTestCase struct {
	ds *MockDataSource
	router *mux.Router
	t *testing.T
}

func NewMockDataSource() *MockDataSource {
	return &MockDataSource{
		[]db.User{
			{1, "Alice", "user"},
			{2, "Soddy", "user"},
			{3, "Sally", "user"},
			{4, "Tim", "user"},
			{5, "June", "user"},
			{6, "Petter", "user"},
		},
		[]db.RelationShip{
			{1, 1, 2, "matched", "relationship"},
			{2, 2, 1, "matched", "relationship"},
			{3, 1, 3, "like", "relationship"},
			{4, 3, 1, "dislike", "relationship"},
			{5, 1, 4, "like", "relationship"},
			{6, 1, 5, "like", "relationship"},
		},
	}
}

func (m *MockDataSource) Users() ([]db.User, error) {
	return m.userDS, nil
}
func (m *MockDataSource) CreateUser(name string) (*db.User, error) {
	user := db.User{m.userDS[len(m.userDS)-1].Id+1, name, "user"}
	m.userDS = append(m.userDS, user)
	return &user, nil
}

func (m *MockDataSource) UserRelationShips(uid int64) ([]db.RelationShip, error) {
	userRSs := []db.RelationShip{}
	for _, i := range m.rsDS {
		if i.Uid == uid {
			userRSs = append(userRSs, i)
		}
	}
	return userRSs, nil
}
func (m *MockDataSource) UserRelationShip(uid, oid int64, tp string) (*db.RelationShip, error) {
	for _, i := range m.rsDS {
		if i.Uid == uid && i.Oid == oid && i.Type == tp {
			return &i, nil
		}
	}
	return nil, errors.New(ERR_NO_RECORD)
}

func (m *MockDataSource) UserRelationShipCount(int64, int64, string) (int, error) {
	return -1, nil
}

func (m *MockDataSource) CreateOrUpdateRelationShip(uid, oid int64, st, tp string) (*db.RelationShip, error) {
	var existRS *db.RelationShip
	for _, i := range m.rsDS {
		if i.Uid == uid && i.Oid == oid {
			existRS = &i
		}
	}
	if existRS == nil {
		rs := db.RelationShip{
			Id: m.rsDS[len(m.rsDS)-1].Id+1,
			Uid: uid,
			Oid: oid,
			State: st,
			Type: tp,
		}
		m.rsDS = append(m.rsDS, rs)
		return &m.rsDS[len(m.rsDS)-1], nil
	} else {
		existRS.State = st
		existRS.Type = tp
		return existRS, nil
	}
}

func (c *BaseTestCase) NewCase(caseName, method, url string, input, output interface{}) *httptest.ResponseRecorder {
	var r *http.Request
	if input != nil {
		postData, _ := json.Marshal(input)
		r, _ = http.NewRequest(method, url, bytes.NewReader(postData))
	} else {
		r, _ = http.NewRequest(method, url, nil)
	}
	w := httptest.NewRecorder()
	c.router.ServeHTTP(w, r)

	if w.Code != 200 {
		c.t.Fatalf("%s failed, HTTP Status Code (%d) != 200", caseName, w.Code)
	}
	switch reflect.ValueOf(output).Kind() {
	case reflect.Ptr:
		err := json.NewDecoder(w.Body).Decode(output)
		if err != nil {
			c.t.Fatalf("Response Decode Failed: %s\n%s", err, debug.Stack())
		}
	}
	return w
}



func initRouter(ds db.DataSource) *mux.Router {
	mockEnv := Env{ ds }
	router := NewWebV1Router(&mockEnv)
	return router
}

func TestUser(t *testing.T) {
	mockDS := NewMockDataSource()
	router := initRouter(mockDS)
	c := BaseTestCase{ mockDS, router, t}

	var users []db.User
	c.NewCase("Get User", "GET", "/users", nil, &users)
	if len(users) != len(mockDS.userDS) {
		t.Fatal("Unexpected Users")
	}

	postUser := struct { Name string } { Name: "Rucker" }
	var responseUser db.User
	_ = c.NewCase("Create User", "POST", "/users", postUser, &responseUser)
	if responseUser.Id != mockDS.userDS[len(mockDS.userDS)-1].Id {
		t.Fatalf("Unexpected User")
	}
}

func TestPutUserRelationShip(t *testing.T) {
	mockDS := NewMockDataSource()
	router := initRouter(mockDS)
	c := BaseTestCase{ mockDS, router, t}

	likeRS := struct { State string } { State: db.RELATIONSHIP_LIKE }
	dislikeRS := struct { State string } { State: db.RELATIONSHIP_DISLIKE }

	type tmpRS struct{
		User_Id int64
		State, Type string
	}

	var respRS tmpRS
	c.NewCase("Like User", "PUT", "/users/1/relationships/6", likeRS, &respRS)
	if respRS.State != db.RELATIONSHIP_LIKE {
		t.Fatalf("Unexpeced Relationship")
	}

	c.NewCase("Be Like", "PUT", "/users/6/relationships/1", likeRS, &respRS)
	if respRS.State != db.RELATIONSHIP_MATCHED {
		t.Fatalf("Unexpeced Relationship")
	}

	c.NewCase("Dislike User", "PUT", "/users/1/relationships/6", dislikeRS, &respRS)
	if respRS.State != db.RELATIONSHIP_DISLIKE {
		t.Fatalf("Unexpeced Relationship: %#v", respRS)
	}
	var (
		userRSs []tmpRS
		respTmpRS *tmpRS
	)
	c.NewCase("Be Dislike", "GET", "/users/6/relationships", nil, &userRSs)
	for _, i := range userRSs {
		if i.User_Id == 1 {
			respTmpRS = &i
		}
	}
	if respTmpRS == nil || respTmpRS.State != db.RELATIONSHIP_LIKE {
		t.Fatalf("Unexpeced Relationship: %#v", respTmpRS)
	}

}

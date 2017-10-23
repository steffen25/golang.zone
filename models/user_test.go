package models

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
	"github.com/go-sql-driver/mysql"
)

func TestUserIsAdmin(t *testing.T) {
	u := createUser()
	isAdmin := u.IsAdmin()
	equals(t, false, isAdmin)
}

func TestCorrectPassword(t *testing.T) {
	u := createUser()
	pw := "awesome password"
	u.SetPassword(pw)
	equals(t, true, u.CheckPassword(pw))
}

func TestWrongPassword(t *testing.T) {
	u := createUser()
	pw := "awesome password"
	u.SetPassword(pw)
	equals(t, false, u.CheckPassword("random"))
}

func TestMarshalJSON(t *testing.T) {
	u := &User{
		Name:      "Thomas",
		Email:     "thomas@email.com",
		Admin:     false,
		CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UpdatedAt: mysql.NullTime{Time: time.Date(2009, time.November, 10, 23, 1, 0, 0, time.UTC), Valid: true},
	}

	json, e := u.MarshalJSON()
	if e != nil {
		t.Fail()
	}

	expectedJson := "{\"id\":0,\"name\":\"Thomas\",\"email\":\"thomas@email.com\",\"createdAt\":\"2009-11-10T23:00:00Z\",\"updatedAt\":\"2009-11-10T23:01:00Z\"}"

	equals(t, string(json), expectedJson)
}

func TestMarshalJSONEmptyUpdatedAt(t *testing.T) {
	u := &User{
		Name:      "Thomas",
		Email:     "thomas@email.com",
		Admin:     false,
		CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		UpdatedAt: mysql.NullTime{Time: time.Time{}, Valid: false},
	}

	json, e := u.MarshalJSON()
	if e != nil {
		t.Fail()
	}

	expectedJson := "{\"id\":0,\"name\":\"Thomas\",\"email\":\"thomas@email.com\",\"createdAt\":\"2009-11-10T23:00:00Z\",\"updatedAt\":null}"

	equals(t, string(json), expectedJson)
}

func createUser() *User {
	u := &User{
		Name:      "Thomas",
		Email:     "thomas@email.com",
		Admin:     false,
		CreatedAt: time.Now(),
	}

	return u
}

// TODO: Move this into its own test package or such for reusability
// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

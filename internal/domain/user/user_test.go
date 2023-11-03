package user

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUser_JsonMarshal(t *testing.T) {
	u := createUser()

	j, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}

	var up = &User{}
	err = json.Unmarshal(j, up)
	if err != nil {
		t.Error(err)
	}

	if u.Email != up.Email {
		t.Errorf("assertation failed, expected Email: %s, but actual Email is: %s.", u.Email, up.Email)
	}
}

func TestUserEnableItAndDisableIt(t *testing.T) {
	u := createUser()
	if u.Enabled == true {
		t.Errorf("assertation failed, expected is that new user is disabled by default.")
	}
	u.EnableIt()
	if u.Enabled == false {
		t.Errorf("assertation failed, expected is that a user is enabled after EnableIt.")
	}
	u.DisableIt()
	if u.Enabled == true {
		t.Errorf("assertation failed, expected is that a user is disabled after DisableIt.")
	}

}

func createUser() *User {
	bd := time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC)

	return NewUser(
		WithEmail("test@fake.com"),
		WithFullName("Testing Tester"),
		WithDateOfBirth(bd),
		WithLocation("Austria, Vienna"),
		WithGender(MALE),
		WithCreatedAt(time.Now()),
		WithUpdatedAt(time.Now()),
	)
}

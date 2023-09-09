package data

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUserProfile_JsonMarshal(t *testing.T) {
	u := UserProfile{
		ID:          1,
		Email:       "test@fake.com",
		FullName:    "Testing Tester",
		DateOfBirth: time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Location:    "Austria, Vienna",
		Gender:      MALE,
		Enabled:     true,
	}
	j, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}

	var up = &UserProfile{}
	err = json.Unmarshal(j, up)
	if err != nil {
		t.Error(err)
	}

	if u.ID != up.ID {
		t.Errorf("assertation failed, expected ID: %d, but actual ID is: %d", u.ID, up.ID)
	}
}

package repos

import (
	"fmt"
	"testing"

	"github.com/fmiskovic/go-starter/internal/core/domain/user"
	"github.com/fmiskovic/go-starter/internal/utils/testx"
	"github.com/google/uuid"
	"github.com/matryer/is"
)

// 1	2623291902 ns/op	2592136 B/op	21598 allocs/op
// 1	2678776323 ns/op	2780120 B/op	22677 allocs/op

func Benchmark_UserRepo_Update(b *testing.B) {
	assert := is.New(b)

	// setup db
	testDb, err := testx.SetUpDb()
	if err != nil {
		b.Errorf("failed to run test db: %v", err)
	}
	defer testDb.Shutdown()

	repo := NewUserRepo(testDb.BunDb)

	u := user.New(user.Email("updated1@fake.com"))
	u.ID = uuid.MustParse("220cea28-b2b0-4051-9eb6-9a99e451af03")

	for n := 0; n < b.N; n++ {
		u.Location = fmt.Sprintf("Location %d", n)
		err := repo.Update(testDb.Ctx, u)
		assert.NoErr(err)
	}
}

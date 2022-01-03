package migrations

import "testing"

func TestMigrationGameDB(t *testing.T) {
	err := MigrateGameDB()
	if err != nil {
		t.Error(err)
	}
}

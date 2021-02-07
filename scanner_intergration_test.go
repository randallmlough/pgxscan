// +build integration

package pgxscan_test

import (
	"context"
	"testing"

	"github.com/randallmlough/pgxscan"
)

func Test_NewScanner(t *testing.T) {
	row := newTestDB(t).QueryRow(context.Background(), `SELECT COUNT(*) FROM "test"`)
	var count int
	if err := pgxscan.NewScanner(row).Scan(&count); err != nil {
		t.Errorf("Test_New() failed to scan into id. Reason:  %v", err)
		return
	}
	if count != 2 {
		t.Errorf("Test_New() wrong count returned. got: %v want: %v", count, 2)
		return
	}
}

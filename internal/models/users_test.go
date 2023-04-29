package models

import (
	"snippetbox/internal/tests"
	"testing"
)

func Test_UserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	testCases := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDb(t)

			model := UserModel{db}

			exists, err := model.Exists(tt.userID)

			tests.Equal(t, exists, tt.want)
			tests.NilError(t, err)
		})
	}
}

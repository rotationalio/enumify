package enumify_test

import (
	"testing"

	"go.rtnl.ai/enumify"
)

func TestTestSuite(t *testing.T) {
	suite := enumify.TestSuite[Status, []string]{
		Values: []Status{StatusUnknown, StatusDraft, StatusReview, StatusPublished, StatusArchived},
		Names:  StatusNames,
		ICase:  true,
		ISpace: true,
	}
	t.Run("Status", suite.Run)
}

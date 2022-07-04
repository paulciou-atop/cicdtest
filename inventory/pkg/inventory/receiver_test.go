package inventory

import (
	"testing"

	lop "github.com/samber/lo/parallel"
)

func TestEmptyMap(t *testing.T) {
	ss := lop.Map([]string{}, func(s string, _ int) string {
		return s + "er"
	})
	t.Log(ss)
}

package lib

import "testing"

func TestContains(t *testing.T) {
	t.Run("True", func(t *testing.T) {
		input := Contains([]string{"a", "b", "c"}, "b")

		if input != true {
			t.Errorf("got %v want %v", input, true)
		}
	})

	t.Run("False", func(t *testing.T) {
		input := Contains([]string{"a", "b", "c"}, "g")

		if input != false {
			t.Errorf("got %v want %v", input, false)
		}
	})
}

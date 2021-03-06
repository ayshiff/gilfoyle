package transcoding

import "testing"

func TestUtil(t *testing.T) {
	var suite = map[string]int8{
		"25/1":    25,
		"30/1":    30,
		"16/1":    16,
		"30":      30,
		"25":      25,
		"invalid": 0,
	}

	for input, expected := range suite {
		if fps := ParseFrameRates(input); fps != expected {
			t.Errorf("expected: %d, got: %d", expected, fps)
		}
	}
}

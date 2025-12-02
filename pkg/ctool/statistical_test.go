package ctool

import "testing"

func TestSimple(t *testing.T) {
	var s = []float64{3613, 3500, 2873, 4780, 3033, 3681, 3372, 3267, 3290, 3693, 3813, 2787, 0, 0, 0, 0, 0, 0}
	print(NumberDistribution(s).Tips)
}

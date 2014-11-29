package quadtree

import (
	"testing"
)

func BenchmarkSearch(b *testing.B) {
	pointLocs := map[string][2]float64{
		"BOM": {19.0901765, 72.8687391},
		"CDG": {49.0096906, 2.5479245},
		"DXB": {25.2531745, 55.3656728},
		"GIG": {-22.9111438, -43.1648755},
		"JFK": {40.6413111, -73.77813909999999},
		"LCY": {51.5048437, 0.049518},
		"MAN": {53.3536987, -2.27495},
	}
	center := NewPoint(0, 0, nil)
	half := center.HalfPoint(10000000)
	qt := New(NewAABB(center, half), 0, nil)

	for iata, coords := range pointLocs {
		qt.Insert(NewPoint(coords[1], coords[0], iata))
	}

	bCenter := NewPoint(-0.7470703125, 52.45600939264076, nil)
	bHalf := center.HalfPoint(182107)
	bb := NewAABB(bCenter, bHalf)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		qt.Search(bb)
	}
}

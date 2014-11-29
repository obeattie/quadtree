package quadtree

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestQuadtree(t *testing.T) {
	suite.Run(t, new(QuadTreeTestSuite))
}

type QuadTreeTestSuite struct {
	suite.Suite
	qt        *QuadTree
	pointLocs map[string][2]float64
	points    map[string]*Point
}

func (suite *QuadTreeTestSuite) SetupSuite() {
	suite.pointLocs = map[string][2]float64{
		"BOM": {19.0901765, 72.8687391},
		"CDG": {49.0096906, 2.5479245},
		"DXB": {25.2531745, 55.3656728},
		"GIG": {-22.9111438, -43.1648755},
		"JFK": {40.6413111, -73.77813909999999},
		"LCY": {51.5048437, 0.049518},
		"MAN": {53.3536987, -2.27495},
	}
}

func (suite *QuadTreeTestSuite) SetupTest() {
	t := suite.T()

	center := NewPoint(0, 0, nil)
	half := center.HalfPoint(10000000)
	suite.qt = New(NewAABB(center, half), 0, nil)
	suite.points = make(map[string]*Point, len(suite.pointLocs))

	for iata, coords := range suite.pointLocs {
		p := NewPoint(coords[1], coords[0], iata)
		assert.True(t, suite.qt.Insert(p))
		suite.points[iata] = p
	}
}

func (suite *QuadTreeTestSuite) TestSearch() {
	t := suite.T()

	bboxes := map[[3]float64][]string{
		[3]float64{51.5134699374858, 0.04094123840332031, 950}:      {},
		[3]float64{51.5134699374858, 0.04094123840332031, 1193}:     {"LCY"},
		[3]float64{31.203404950917395, 59.94140624999999, 2264075}:  {"BOM", "DXB"},
		[3]float64{52.45600939264076, -0.7470703125, 182107}:        {"LCY", "MAN"},
		[3]float64{52.45600939264076, -0.7470703125, 1821070000000}: {"BOM", "CDG", "DXB", "GIG", "JFK", "LCY", "MAN"},
	}

	for boxCoords, expectedIatas := range bboxes {
		center := NewPoint(boxCoords[1], boxCoords[0], nil)
		half := center.HalfPoint(boxCoords[2])
		results := suite.qt.Search(NewAABB(center, half))
		resultIatas := make([]string, len(results))
		for i, r := range results {
			resultIatas[i] = r.Data().(string)
		}
		sort.Strings(resultIatas)
		assert.Equal(t, expectedIatas, resultIatas)
	}
}

func (suite *QuadTreeTestSuite) TestRemove() {
	t := suite.T()

	// Take out LCY
	assert.True(t, suite.qt.Remove(suite.points["LCY"]))

	bboxes := map[[3]float64][]string{
		[3]float64{51.5134699374858, 0.04094123840332031, 950}:      {},
		[3]float64{51.5134699374858, 0.04094123840332031, 1193}:     {},
		[3]float64{31.203404950917395, 59.94140624999999, 2264075}:  {"BOM", "DXB"},
		[3]float64{52.45600939264076, -0.7470703125, 182107}:        {"MAN"},
		[3]float64{52.45600939264076, -0.7470703125, 1821070000000}: {"BOM", "CDG", "DXB", "GIG", "JFK", "MAN"},
	}

	for boxCoords, expectedIatas := range bboxes {
		center := NewPoint(boxCoords[1], boxCoords[0], nil)
		half := center.HalfPoint(boxCoords[2])
		results := suite.qt.Search(NewAABB(center, half))
		assert.Len(t, results, len(expectedIatas))
	}
}

func (suite *QuadTreeTestSuite) TestUpdate() {
	t := suite.T()

	// Replace LCY
	previous := suite.points["LCY"]
	replacement := NewPoint(-0.45409262180328364, 51.472198132255066, "LHR")
	assert.True(t, suite.qt.Update(previous, replacement))

	bboxes := map[[3]float64][]string{
		[3]float64{51.5134699374858, 0.04094123840332031, 950}:      {},
		[3]float64{51.5134699374858, 0.04094123840332031, 1193}:     {},
		[3]float64{31.203404950917395, 59.94140624999999, 2264075}:  {"BOM", "DXB"},
		[3]float64{52.45600939264076, -0.7470703125, 182107}:        {"LHR", "MAN"},
		[3]float64{52.45600939264076, -0.7470703125, 1821070000000}: {"BOM", "CDG", "DXB", "GIG", "JFK", "LHR", "MAN"},
	}

	for boxCoords, expectedIatas := range bboxes {
		center := NewPoint(boxCoords[1], boxCoords[0], nil)
		half := center.HalfPoint(boxCoords[2])
		results := suite.qt.Search(NewAABB(center, half))
		assert.Len(t, results, len(expectedIatas))
	}
}

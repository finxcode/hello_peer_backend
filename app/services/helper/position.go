package helper

import "math"

type position struct {
	lng float32
	lat float32
}

func NewPosition(lng, lat float32) *position {
	if lng == 0 && lat == 0 {
		return nil
	}
	return &position{
		lng: lng,
		lat: lat,
	}
}

func (p *position) GetDistance(to *position) int {
	radius := 6378.137
	rad := math.Pi / 180.0
	lat1 := float64(p.lat) * rad
	lng1 := float64(p.lng) * rad
	lat2 := float64(to.lat) * rad
	lng2 := float64(to.lng) * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	dist = dist * radius
	if dist == 0 {
		return 1
	}
	return int(math.Ceil(dist))
}

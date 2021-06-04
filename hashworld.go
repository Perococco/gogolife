package main

type Point struct {
	X int32
	Y int32
}

type HashWorld struct {
	stateIdx int32
}

func (world *HashWorld) Update(rules Rules) {

}
func (world *HashWorld) Draw() {

}
func (world *HashWorld) SetState(x, y int32, state bool) {

}

func (world *HashWorld) EstimateZoom(width int, height int) float32 {
	return 0.0
}

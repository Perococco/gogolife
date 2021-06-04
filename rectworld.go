package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type RectWorld struct {
	width    int32
	height   int32
	stateIdx int32
	states   [2][]bool
}

func NewRectWorld(width, height int32) *RectWorld {
	var states [2][]bool
	for i := 0; i < len(states); i++ {
		states[i] = make([]bool, width*height)
	}
	return &RectWorld{
		width:    width,
		height:   height,
		stateIdx: 0,
		states:   states,
	}
}

const NB_WORKERS = 10

func (world *RectWorld) UpdateRow(rules Rules, w, h int32) {
	nbNeighbours := world.NeighbourCount(w, h)
	newState := rules.StateEstimate(nbNeighbours, world.StateGet(w, h))
	world.SetNewState(w, h, newState)
}

func (world *RectWorld) UpdateAsync(rules Rules, workerIndex int32, channel chan int32) {
	for h := workerIndex;h < world.height;h+=NB_WORKERS {
		for w := int32(0); w < world.width; w++ {
			world.UpdateRow(rules, w,h)
		}
	}
	channel <- 1
}

func (world *RectWorld) Update(rules Rules) {

	channel := make(chan int32)

	nbDone := int32(0)

	for workerIndex := int32(0); workerIndex < NB_WORKERS; workerIndex++ {
		go world.UpdateAsync(rules,workerIndex, channel)
	}

	for nbDone < NB_WORKERS {
		nbDone += <-channel
	}

	world.SwapStates()
}

func (world *RectWorld) SetState(x, y int32, state bool) {
	halfWidth := world.width / 2
	halfHeight := world.height / 2
	world.setNewStateInternal(x+halfWidth, y+halfHeight, state, world.stateIdx)
}

func (world *RectWorld) Draw() {
	halfWidth := world.width / 2
	halfHeight := world.height / 2
	rec := rl.Rectangle{X: -1 - float32(halfWidth), Y: -1 - float32(halfHeight), Width: float32(world.width) + 2, Height: float32(world.height) + 2}
	rl.DrawRectangleLinesEx(rec, 1, rl.Red)
	for w := int32(0); w < world.width; w++ {
		for h := int32(0); h < world.height; h++ {
			if world.StateGet(w, h) {
				rl.DrawRectangle(w-halfWidth, h-halfHeight, 1, 1, rl.White)
			}
		}
	}
}

func (world *RectWorld) SwapStates() {
	world.stateIdx = 1 - world.stateIdx
}

func (world *RectWorld) StateGet(x, y int32) bool {
	mx := mod(x, world.width)
	my := mod(y, world.height)
	return world.states[world.stateIdx][mx+my*world.width]
}

func (world *RectWorld) SetNewState(x, y int32, state bool) {
	world.setNewStateInternal(x, y, state, 1-world.stateIdx)
}

func (world *RectWorld) setNewStateInternal(x, y int32, state bool, stateIdx int32) {
	mx := mod(x, world.width)
	my := mod(y, world.height)
	data := world.states[stateIdx]
	data[mx+my*world.width] = state
}

func (world *RectWorld) EstimateZoom(width int, height int) float32 {
	zoomWidth := float64(width) / float64(world.width)
	zoomHeight := float64(height) / float64(world.height)
	return float32(math.Min(zoomWidth, zoomHeight))
}

func (world *RectWorld) NeighbourCount(x, y int32) int32 {
	nw := boolToInt(world.StateGet(x-1, y-1))
	n := boolToInt(world.StateGet(x, y-1))
	ne := boolToInt(world.StateGet(x+1, y-1))
	w := boolToInt(world.StateGet(x-1, y))
	e := boolToInt(world.StateGet(x+1, y))
	sw := boolToInt(world.StateGet(x-1, y+1))
	s := boolToInt(world.StateGet(x, y+1))
	se := boolToInt(world.StateGet(x+1, y+1))
	return nw + n + ne + w + e + sw + s + se
}

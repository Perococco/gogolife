package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func mod(x, base int32) int32 {
	result := x % base
	if result < 0 {
		return result + base
	}
	return result
}

func boolToInt(state bool) int32 {
	if state {
		return 1
	}
	return 0
}

type World interface {
	Update(rules Rules)
	Draw()
	SetState(x,y int32, state bool)
}

type Rules interface {
	StateEstimate(nbNeighbour int32, isAlive bool) bool
}

type Conway struct {}

func (conway Conway) StateEstimate(nbNeighbour int32, isAlive bool) bool {
	switch nbNeighbour {
	case 2: return isAlive
	case 3: return true
	default: return false
	}
}


func main() {
	world := NewRectWorld(100,100)

	goLState := NewGolState(world, Conway{})

	world.SetNewState(50,50,true)
	world.SetNewState(50,49,true)
	world.SetNewState(51,48,true)
	world.SetNewState(51,50,true)
	world.SetNewState(52,50,true)
	world.SwapStates()

	rl.InitWindow(800, 450, "raylib [core] example - basic window")
	rl.SetWindowState(rl.FlagMsaa4xHint | rl.FlagWindowResizable)
	rl.SetTargetFPS(60)


	for !rl.WindowShouldClose() {

		goLState.UpdateCameraOffset()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Gray)
		rl.BeginMode2D(goLState.camera2d)

		goLState.HandleInputs()
		goLState.Tick()
		goLState.Draw()


		rl.EndMode2D()
		rl.EndDrawing()

	}

	rl.CloseWindow()
}

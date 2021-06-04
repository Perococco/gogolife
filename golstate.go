package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

const ZOOM_FACTOR = 1.2

func NewGolState(world World, rules Rules) GoLState {
	goLState := GoLState{
		updateEnabled:    false,
		updateTickPeriod: 10,
		updateTicks:      0,
		camera2d:         rl.Camera2D{},
		world:            world,
		rules:            rules,
		targetBackup:     rl.Vector2{},
		dragInfo: DragInfo{
			enabled:         false,
			started:         false,
			done:            false,
			startPosition:   rl.Vector2{},
			currentPosition: rl.Vector2{},
		},
	}
	goLState.camera2d.Target.X = 0
	goLState.camera2d.Target.Y = 0
	goLState.camera2d.Zoom = 10

	return goLState
}
func (state *GoLState) UpdateWorld() {
	(state.world).Update(state.rules)
}

type DragInfo struct {
	enabled         bool
	started         bool
	done            bool
	startPosition   rl.Vector2
	currentPosition rl.Vector2
}

type GoLState struct {
	updateTickPeriod int32
	updateTicks      int32
	camera2d         rl.Camera2D
	updateEnabled    bool
	world            World
	targetBackup     rl.Vector2
	rules            Rules
	dragInfo         DragInfo
}

func (state *GoLState) ToggleWorldUpdate() {
	state.updateEnabled = !state.updateEnabled
}

func (state *GoLState) IncreaseUpdateSpeed() {
	state.updateTickPeriod -= 5
	if state.updateTickPeriod <= 0 {
		state.updateTickPeriod = 1
	}
}

func (state *GoLState) DecreaseUpdateSpeed() {
	state.updateTickPeriod += 5
	if state.updateTickPeriod >= 60 {
		state.updateTickPeriod = 60
	}
}

func (state *GoLState) UpdateCameraOffset() {
	state.camera2d.Offset.X = float32(rl.GetScreenWidth()) * 0.5
	state.camera2d.Offset.Y = float32(rl.GetScreenHeight()) * 0.5
}

func (state *GoLState) HandleMouse() {
	state.HandleMouseZoom()
	state.HandleMouseClick()
	state.HandleMouseDrag()
}

func (state *GoLState) HandleInputs() {
	state.HandleKeys()
	state.HandleMouse()
}

func (state *GoLState) HandleMouseZoom() {
	wheelMove := rl.GetMouseWheelMove()

	if wheelMove == 0 {
		return
	}

	mousePosition := rl.GetScreenToWorld2D(rl.GetMousePosition(), state.camera2d)
	zoomFactor := float32(math.Pow(ZOOM_FACTOR, float64(wheelMove)))

	state.camera2d.Zoom *= zoomFactor
	state.camera2d.Target.X -= (mousePosition.X - state.camera2d.Target.X) * (1 - zoomFactor)
	state.camera2d.Target.Y -= (mousePosition.Y - state.camera2d.Target.Y) * (1 - zoomFactor)

}

func (state *GoLState) HandleMouseDrag() {
	state.dragInfo.Update(rl.MouseMiddleButton)

	if state.dragInfo.started {
		state.targetBackup = state.camera2d.Target
	}

	if state.dragInfo.enabled {
		start := rl.GetScreenToWorld2D(state.dragInfo.startPosition, state.camera2d)
		current := rl.GetScreenToWorld2D(state.dragInfo.currentPosition, state.camera2d)

		state.camera2d.Target.X = state.targetBackup.X - current.X + start.X
		state.camera2d.Target.Y = state.targetBackup.Y - current.Y + start.Y

	}

}

func (state *GoLState) HandleKeys() {

	if rl.IsKeyReleased(rl.KeySpace) {
		state.ToggleWorldUpdate()
	}
	if rl.IsKeyReleased(rl.KeyUp) {
		state.IncreaseUpdateSpeed()
	}
	if rl.IsKeyReleased(rl.KeyDown) {
		state.DecreaseUpdateSpeed()
	}
	if rl.IsKeyPressed(rl.KeyHome) {
		state.ResetZoom()
	}
}

func (state *GoLState) Tick() {
	if state.updateEnabled {
		state.updateTicks++
		if state.updateTicks >= state.updateTickPeriod {
			state.UpdateWorld()
			state.updateTicks = 0
		}
	}
}

func (state *GoLState) Draw() {
	state.world.Draw()
}

func (state *GoLState) HandleMouseClick() {
	leftDown := rl.IsMouseButtonDown(rl.MouseLeftButton)
	rightDown := rl.IsMouseButtonDown(rl.MouseRightButton)

	if leftDown == rightDown {
		return
	}

	position := rl.GetScreenToWorld2D(rl.GetMousePosition(), state.camera2d)
	x := int32(position.X)
	y := int32(position.Y)

	state.world.SetState(x, y, leftDown)

}

func (state *GoLState) ResetZoom() {

	state.camera2d.Zoom = state.world.EstimateZoom(rl.GetScreenWidth(), rl.GetScreenHeight())

}

func (dragInfo *DragInfo) Update(button int32) {

	dragInfo.started = rl.IsMouseButtonPressed(button)
	dragInfo.done = rl.IsMouseButtonReleased(button)
	dragInfo.enabled = rl.IsMouseButtonDown(button)

	if dragInfo.started {
		dragInfo.startPosition = rl.GetMousePosition()
	}
	if dragInfo.enabled {
		dragInfo.currentPosition = rl.GetMousePosition()
	}

}

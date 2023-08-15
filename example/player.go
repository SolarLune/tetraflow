package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/tetra3d"
)

type Player struct {
	Model *tetra3d.Model
}

func NewPlayer(node tetra3d.INode) *Player {
	return &Player{Model: node.(*tetra3d.Model)}
}

func (player *Player) OnUpdate(dt float64) {
	move := tetra3d.NewVector(0, 0, 0)
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		move.X = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		move.X = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		move.Z = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		move.Z = 1
	}

	player.Model.MoveVec(move.Unit().Scale(0.1))

	enemies := player.Model.Root().SearchTree().ByParentProps("enemy").IBoundingObjects()

	bounds := player.Model.Get("BoundingAABB").(*tetra3d.BoundingAABB)
	if res := bounds.CollisionTest(tetra3d.CollisionTestSettings{Others: enemies}); len(res) > 0 {
		player.Model.Unparent()
	}

}

func (player *Player) OnRemove() {
	fmt.Println("You're dead")
}

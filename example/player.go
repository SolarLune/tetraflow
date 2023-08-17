package main

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/tetra3d"
	"github.com/solarlune/tetra3d/colors"
	"github.com/solarlune/tetraflow/messages"
)

type Player struct {
	Game         *Game
	Model        *tetra3d.Model
	Invulnerable bool
}

func NewPlayer(node tetra3d.INode, game *Game) *Player {
	player := &Player{
		Model: node.(*tetra3d.Model),
		Game:  game,
	}
	return player
}

func (player *Player) ReceiveMessage(msg messages.IMessage) {

	switch message := msg.(type) {

	case messages.AddToScene:
		fmt.Println("You're alive!")
	case messages.SceneStart:
		fmt.Println("Scene Start! Scene Name: " + message.Scene.Name)
	case messages.SceneEnd:
		fmt.Println("Scene ended. Scene Name: " + message.Scene.Name)
	case messages.Update:

		// You can dispatch events here
		player.Update()

	case messages.RemoveFromScene:
		fmt.Println("you died. Respawning in 3 seconds...")
		layer := player.Game.Engine.Stages[0]
		layer.TimerSystem().After(time.Second*3, func() {
			layer.CurrentScene().Root.AddChildren(player.Model)
			player.Model.ResetWorldPosition()
			player.GoInvulnerable()
		})

	}

}

func (player *Player) Update() {

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

	player.Model.Color = colors.White()

	if !player.Invulnerable {

		if res := bounds.CollisionTest(tetra3d.CollisionTestSettings{Others: enemies}); len(res) > 0 {
			player.Model.Unparent() // Remove the model from the hierarchy; this triggers the RemoveFromScene message.
		}

	} else {
		s := float32(math.Sin(player.Game.Engine.Time() * math.Pi * 4))
		c := colors.White().AddRGBA(s, s, s, 0)
		player.Model.Color = c
	}

}

func (player *Player) GoInvulnerable() {
	player.Invulnerable = true
	ts := player.Game.Engine.Stages[0].TimerSystem()
	ts.After(time.Second*2, func() { player.Invulnerable = false })
}

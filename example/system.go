package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/tetra3d"
	"github.com/solarlune/tetraflow/messages"
)

type System struct {
	Node tetra3d.INode
	Game *Game
}

func NewSystem(node tetra3d.INode, game *Game) *System {
	return &System{
		Node: node,
		Game: game,
	}
}

func (system *System) ReceiveMessage(msg messages.IMessage) {

	if msg.Type() == messages.TypeUpdate {

		// Quit
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			system.Game.Engine.Quit()
		}

		// Fullscreen
		if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
		}

		// Restart
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			// Restarting is done here by simply recloning the scene; a better, more memory-efficient way to do this would be to simply
			// re-place the player in his original location, for example.
			fmt.Println("restart")
			system.Game.Engine.FindStageByName("Game").Restart()
		}

		// Debug stuff

		if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
			system.Game.DrawDebugStats = !system.Game.DrawDebugStats
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
			system.Game.DrawDebugDepth = !system.Game.DrawDebugDepth
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			layer := system.Game.Engine.FindStageByName("Game")
			if layer.Paused() {
				layer.Unpause()
			} else {
				layer.Pause()
			}
		}

	}

}

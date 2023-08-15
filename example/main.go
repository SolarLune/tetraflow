package main

import (
	"errors"

	_ "embed"

	"github.com/solarlune/tetra3d"
	"github.com/solarlune/tetra3d/colors"
	"github.com/solarlune/tetraflow"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Width, Height  int
	Library        *tetra3d.Library
	Scene          *tetra3d.Scene
	DrawDebugDepth bool
	DrawDebugStats bool

	Flow *tetraflow.EngineFlow
}

// The goal of this example is to make a simple quickstart project for basing new projects off of.
// In this example, the Tetra3D icon spins in the center of the screen. It is shadeless, so you would
// either need to add a light to the scene, use other shadeless materials, or disable lighting on the
// scene to be able to see other new objects.

//go:embed startingScene.gltf
var startingGLTF []byte

func NewGame() *Game {

	game := &Game{
		Width:  796,
		Height: 448,
	}

	game.Init()

	return game
}

func (g *Game) Init() {

	if g.Library == nil {

		options := tetra3d.DefaultGLTFLoadOptions()
		options.CameraWidth = g.Width
		options.CameraHeight = g.Height
		library, err := tetra3d.LoadGLTFData(startingGLTF, options)
		if err != nil {
			panic(err)
		}

		g.Library = library

	}

	g.Scene = g.Library.ExportedScene.Clone()

	// OK, so we'll create a new Flow.
	// A Flow essentially handles calling OnUpdate(), OnAdd(), and OnRemove() on objects in the scene tree.

	// You create a Flow by specifying the Scene it's running on, and a registry function to be called.

	// The Registry function is a function that you define that will return a struct that will be placed in
	// relevant nodes' Data pointers, and these objects are the "game logic objects" for your game.

	g.Flow = tetraflow.NewEngineFlow(
		g.Scene,
		func(node tetra3d.INode) any {

			if node.Properties().Has("gameobject") {
				switch node.Properties().Get("gameobject").AsString() {
				case "player":
					return NewPlayer(node)
				case "enemy":
					return NewEnemy(node)
				}
			}

			return nil

		},
	)

}

func (g *Game) Update() error {

	var err error

	// Quit
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		err = errors.New("quit")
	}

	// Fullscreen
	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Restart
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		// Restarting is done here by simply recloning the scene; a better, more memory-efficient way to do this would be to simply
		// re-place the player in his original location, for example.
		g.Init()
	}

	// Debug stuff

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.DrawDebugStats = !g.DrawDebugStats
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		g.DrawDebugDepth = !g.DrawDebugDepth
	}

	// The engine flow code updates here
	g.Flow.Update()

	return err
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(g.Scene.World.ClearColor.ToRGBA64())

	camera := g.Scene.Root.Get("Camera").(*tetra3d.Camera)

	camera.Clear()
	camera.RenderNodes(g.Scene, g.Scene.Root)

	if g.DrawDebugDepth {
		screen.DrawImage(camera.DepthTexture(), nil)
	} else {
		screen.DrawImage(camera.ColorTexture(), nil)
	}

	if g.DrawDebugStats {
		camera.DrawDebugRenderInfo(screen, 1, colors.White())
	}

}

func (g *Game) Layout(w, h int) (int, int) {
	// This is a fixed aspect ratio; we can change this to, say, extend for wider displays by using the provided w argument and
	// calculating the height from the aspect ratio, then calling Camera.Resize() on any / all cameras with the new width and height.
	return g.Width, g.Height
}

func main() {

	ebiten.SetWindowTitle("Tetra3d - Quickstart Project")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()

	// An ungraceful quit
	if err := ebiten.RunGame(game); err != nil && err.Error() != "quit" {
		panic(err)
	}

}

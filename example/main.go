package main

import (
	_ "embed"

	"github.com/solarlune/tetra3d"
	"github.com/solarlune/tetra3d/colors"
	"github.com/solarlune/tetraflow"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Width, Height  int
	Library        *tetra3d.Library
	DrawDebugDepth bool
	DrawDebugStats bool

	Engine *tetraflow.Engine
}

// The goal of this example is to make a simple example project for use of TetraFlow.

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

	// OK, so we'll create a new Flow.
	// A Flow essentially handles updating an engine and instantiating scenes.

	// You create a Flow by specifying the Scene it should clone, and a registry function to be called.

	// The Registry function is a function that you define that will return a struct that will be placed in
	// relevant nodes' Data pointers, and these objects are the "game logic objects" for your game.
	// Nodes will only be populated with IReceiver objects if they don't have anything in their Data() spots already.

	g.Engine = tetraflow.NewEngine(g.Library, func(node tetra3d.INode) tetraflow.IReceiver {

		if node.Properties().Has("gameobject") {
			switch node.Properties().Get("gameobject").AsString() {
			case "player":
				return NewPlayer(node, g)
			case "enemy":
				return NewEnemy(node)
			case "system":
				return NewSystem(node, g)
			}
		}

		return nil

	})

	layer := g.Engine.AddStage("Game")
	layer.SetSceneByName("Level")

}

func (g *Game) Update() error {
	return g.Engine.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {

	scene := g.Engine.Stages[0].CurrentScene()

	screen.Fill(scene.World.ClearColor.ToRGBA64())

	camera := scene.Root.SearchTree().ByType(tetra3d.NodeTypeCamera).First().(*tetra3d.Camera)

	camera.Clear()
	camera.RenderScene(scene)

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

	ebiten.SetWindowTitle("TetraFlow Test Project")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()

	// An ungraceful quit
	if err := ebiten.RunGame(game); err != nil && err.Error() != "quit" {
		panic(err)
	}

}

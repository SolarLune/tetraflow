package tetraflow

import (
	"errors"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ebitick"
	"github.com/solarlune/tetra3d"
	"github.com/solarlune/tetraflow/messages"
)

var EngineQuit = errors.New("quit")

// ReceiverRegistryFunction is a function to return something that implements IOnUpdate or IOnRemove
// (i.e. something game object-like or component-like) for relevant nodes.
type ReceiverRegistryFunction func(node tetra3d.INode) IReceiver

// Stage represents a layer or box in which Scenes are created, updated, replaced, and ended.
// The scene might change, but the stage generally always exists (i.e. there's one for GUI elements, one for the game world, and one for background skyboxes).
type Stage struct {
	Name   string
	Active bool
	engine *Engine

	started bool
	paused  bool

	remove          bool
	scene           *tetra3d.Scene
	sourceScene     *tetra3d.Scene
	receivers       []tetra3d.INode
	prevReceivers   []tetra3d.INode
	timerSystem     *ebitick.TimerSystem
	Speed           float64 // Speed represents the speed of the engine's playback (for animations, delta time values, and timers).
	DebugUpdateTime time.Duration
}

// newStage creates a new Stage.
func newStage(engine *Engine, name string) *Stage {

	return &Stage{
		Name:        name,
		engine:      engine,
		Speed:       1,
		Active:      true,
		timerSystem: ebitick.NewTimerSystem(),
	}
}

// Update updates the elements within the Stage.
func (stage *Stage) Update() {

	stage.DebugUpdateTime = 0

	if stage.sourceScene == nil || stage.paused || !stage.Active {
		return
	}

	start := time.Now()

	justStarted := false

	if !stage.started {
		stage.scene = stage.sourceScene.Clone()
		stage.started = true
		justStarted = true
	}

	if stage.started {

		dt := 1.0 / float64(ebiten.TPS()) * stage.Speed

		stage.receivers = make([]tetra3d.INode, 0, len(stage.prevReceivers))

		stage.scene.Root.SearchTree().ForEach(
			func(node tetra3d.INode) bool {

				if node.Data() == nil {

					if data := stage.engine.registryFunction(node); data != nil {
						node.SetData(data)
					}

				}

				_, ok := node.Data().(IReceiver)

				if ok {
					stage.receivers = append(stage.receivers, node)
				}

				return true

			},
		)

		if justStarted {

			for _, g := range stage.receivers {
				g.Data().(IReceiver).ReceiveMessage(messages.SceneStart{Scene: stage.scene})
			}

		}

		for _, node := range stage.receivers {

			existed := false
			for _, p := range stage.prevReceivers {
				if node == p {
					existed = true
					break
				}
			}

			receiver := node.Data().(IReceiver)

			if !existed {
				receiver.ReceiveMessage(messages.AddToScene{Scene: stage.scene})
			}

			receiver.ReceiveMessage(messages.Update{})

			node.AnimationPlayer().Update(dt)

		}

		for _, node := range stage.prevReceivers {
			stillExists := false
			for _, g := range stage.receivers {
				if g == node {
					stillExists = true
					break
				}
			}
			if !stillExists {
				node.Data().(IReceiver).ReceiveMessage(messages.RemoveFromScene{Scene: stage.scene})
			}
		}

		stage.prevReceivers = stage.receivers

		// Set the timer system speed
		stage.timerSystem.Speed = stage.Speed

		stage.timerSystem.Update()

	}

	stage.DebugUpdateTime = time.Since(start)

}

// Pause pauses the Stage.
func (stage *Stage) Pause() {
	stage.paused = true
}

// Unpause unpauses the Stage.
func (stage *Stage) Unpause() {
	stage.paused = false
}

// Paused returns if the Stage is paused.
func (stage *Stage) Paused() bool {
	return stage.paused
}

// End stops updating the Scene.
func (stage *Stage) End() {
	stage.timerSystem.Clear()
	stage.started = false
	stage.Active = false

	stage.prevReceivers = make([]tetra3d.INode, 0, len(stage.prevReceivers))

	for _, g := range stage.receivers {
		g.Data().(IReceiver).ReceiveMessage(messages.SceneEnd{Scene: stage.scene})
	}

}

// Restart restarts the Stage on the next Update() frame, effectively re-instantiating the Scene in the Stage.
// This calls End on the Stage before doing so.
func (stage *Stage) Restart() {
	if stage.started {
		stage.End()
	}
	stage.started = false
	stage.Active = true
}

// Remove queues the Stage for removal from its owning Engine instance.
func (stage *Stage) Remove() {
	stage.remove = true
}

// CurrentScene returns the currently running Scene within the Stage.
// If the Stage is inactive or the scene is ended (or the Stage has not instantiated a scene yet), this function will return nil.
func (stage *Stage) CurrentScene() *tetra3d.Scene {
	if stage.Active {
		return stage.scene
	}
	return nil
}

// SetSceneByName sets the scene to be run by name, searching through the libraries associated with the Engine.
// If the Stage is running, this will reset the currently running scene.
// If no scene of the specified name is found, this function will return an error.
func (stage *Stage) SetSceneByName(sceneName string) error {

	var sourceScene *tetra3d.Scene

	for _, l := range stage.engine.libraries {
		scene := l.FindScene(sceneName)
		if scene != nil {
			sourceScene = scene
			break
		}
	}

	if sourceScene == nil {
		return errors.New("error: SetSceneByName cannot find a scene with specified name: " + sceneName)
	}
	stage.SetScene(sourceScene)
	return nil

}

// SetScene sets the source scene to be run
// If the Stage is running, this will reset the currently running scene.
// If no scene of the specified name is found, this function will return an error.
func (stage *Stage) SetScene(scene *tetra3d.Scene) error {

	if scene == nil {
		return errors.New("error: SetScene cannot be used with a nil Scene")
	}

	stage.sourceScene = scene

	if stage.started {
		stage.Restart()
	}

	return nil

}

// TimerSystem returns the Stage's built-in TimerSystem for timing events.
func (stage *Stage) TimerSystem() *ebitick.TimerSystem {
	return stage.timerSystem
}

// Engine represents an engine for a Tetra3D scene.
type Engine struct {
	time             float64
	libraries        []*tetra3d.Library
	Stages           []*Stage
	registryFunction ReceiverRegistryFunction
	quit             bool
}

// NewEngine creates a new Engine.
// library is the source library to be used for instantiating scenes from.
// ReceiverRegistryFunction is a function to return an IGameObject (something that can receive messages)
// for relevant INodes.
func NewEngine(library *tetra3d.Library, registryFunction ReceiverRegistryFunction) *Engine {
	return &Engine{
		libraries:        []*tetra3d.Library{library},
		Stages:           []*Stage{},
		registryFunction: registryFunction,
	}
}

// AddLibrary adds libraries to use for scene creation to the Engine.
func (engine *Engine) AddLibrary(libraries ...*tetra3d.Library) {
	engine.libraries = append(engine.libraries, libraries...)
}

// Update updates all relevant game object-implementers in the Engine's Scene's node hierarchy,
// and enables all of their callbacks (OnUpdate(), OnRemove(), etc).
// Update will return an error if it has been queued to quit.
func (engine *Engine) Update() error {

	for _, stage := range engine.Stages {
		stage.Update()
	}

	// Remove stages that shouldn't exist anymore
	for i := len(engine.Stages) - 1; i > 0; i-- {
		stage := engine.Stages[i]
		if stage.remove {
			engine.Stages = append(engine.Stages[:i], engine.Stages[i+1:]...)
		}
	}

	engine.time += 1.0 / float64(ebiten.TPS())

	if engine.quit {
		return EngineQuit
	}

	return nil

}

// AddStage creates a new Stage with the specified name.
// The Stage would then need to be set with a Scene to create. After that, it will be instantiated and begin updating on the next Engine.Update() call.
func (engine *Engine) AddStage(name string) *Stage {

	st := newStage(engine, name)
	engine.Stages = append(engine.Stages, st)
	return st

}

// FindStageByName searches for a Stage by the name of the scene that should be running in it.
// If the scene isn't found, the function returns nil.
func (engine *Engine) FindStageByName(stageName string) *Stage {
	for _, l := range engine.Stages {
		if l.Name == stageName {
			return l
		}
	}
	return nil
}

// FindStageByScene searches for a Stage by the scene that should be running in it.
// If the scene isn't found, the function returns nil.
func (engine *Engine) FindStageByScene(scene *tetra3d.Scene) *Stage {
	for _, s := range engine.Stages {
		if s.scene == scene {
			return s
		}
	}
	return nil
}

// RemoveStage queues a Stage for removal, and is done at the end of the Engine's update cycle.
func (engine *Engine) RemoveStage(sceneName string) {
	for _, s := range engine.Stages {
		if s.scene.Name == sceneName {
			s.remove = true
			break
		}
	}
}

// SendMessage sends the specified message to all Stages in the Engine.
func (engine *Engine) SendMessage(msg messages.IMessage) {

	for _, stage := range engine.Stages {
		for _, node := range stage.receivers {
			node.Data().(IReceiver).ReceiveMessage(msg)
		}
	}

}

// SendMessageToTargets sends the specified message to all of the target Nodes, assuming they have IReceivers implemented in their Data() slots.
func (engine *Engine) SendMessageToTargets(msg messages.IMessage, targets ...tetra3d.INode) {

	for _, node := range targets {

		if g, ok := node.Data().(IReceiver); ok {
			g.ReceiveMessage(msg)
		}

	}

}

// Quit queues the engine to quit gracefully.
func (engine *Engine) Quit() {
	engine.quit = true
}

// Time returns the total time that the Engine has been running in seconds.
func (engine *Engine) Time() float64 {
	return engine.time
}

type IReceiver interface {
	ReceiveMessage(msg messages.IMessage)
}

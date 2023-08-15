package tetraflow

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/ebitick"
	"github.com/solarlune/tetra3d"
)

// GameObjectRegistryFunction is a function to return something that implements IOnUpdate or IOnRemove
// (i.e. something game object-like or component-like) for relevant nodes.
type GameObjectRegistryFunction func(node tetra3d.INode) any

// EngineFlow represents an engine flow for a Tetra3D scene.
type EngineFlow struct {
	Timescale        float64
	Scene            *tetra3d.Scene
	time             float64
	GameObjects      []tetra3d.INode
	PrevGameObjects  []tetra3d.INode
	registryFunction GameObjectRegistryFunction
	timerSystem      *ebitick.TimerSystem
}

func NewEngineFlow(scene *tetra3d.Scene, registryFunction GameObjectRegistryFunction) *EngineFlow {
	return &EngineFlow{
		Timescale:        1,
		Scene:            scene,
		registryFunction: registryFunction,
		timerSystem:      ebitick.NewTimerSystem(),
	}
}

func (flow *EngineFlow) Update() {

	dt := 1.0 / float64(ebiten.TPS()) * flow.Timescale

	flow.GameObjects = make([]tetra3d.INode, 0, len(flow.PrevGameObjects))

	flow.Scene.Root.SearchTree().ForEach(
		func(node tetra3d.INode) bool {

			if data := flow.registryFunction(node); data != nil {
				node.SetData(data)
			}

			_, updateable := node.Data().(IOnUpdate)
			_, removeable := node.Data().(IOnRemove)

			if updateable || removeable {
				flow.GameObjects = append(flow.GameObjects, node)
			}

			return true

		},
	)

	for _, g := range flow.GameObjects {

		if update, ok := g.Data().(IOnUpdate); ok {
			update.OnUpdate(dt)
		}

		g.AnimationPlayer().Update(dt)

	}

	for _, p := range flow.PrevGameObjects {
		stillExists := false
		for _, g := range flow.GameObjects {
			if g == p {
				stillExists = true
				break
			}
		}
		if !stillExists {
			if g, ok := p.Data().(IOnRemove); ok {
				g.OnRemove()
			}
		}
	}

	flow.PrevGameObjects = flow.GameObjects

	flow.time += dt

	// Set the timer system speed
	flow.timerSystem.Speed = flow.Timescale

	flow.timerSystem.Update()

}

// TimerSystem returns the EngineFlow's built-in TimerSystem for timing events.
func (flow *EngineFlow) TimerSystem() *ebitick.TimerSystem {
	return flow.timerSystem
}

// Time returns the total time that the EngineFlow has been running in seconds.
func (flow *EngineFlow) Time() float64 {
	return flow.time
}

type IOnUpdate interface {
	OnUpdate(delta float64)
}

// type IOnAdd interface {
// 	OnAdd()
// }

type IOnRemove interface {
	OnRemove()
}

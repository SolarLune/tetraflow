package messages

import "github.com/solarlune/tetra3d"

type MessageType int

type IMessage interface {
	Type() MessageType
}

const (
	TypeUpdate MessageType = iota
	TypeAddToScene
	TypeRemoveFromScene
	TypeSceneStart
	TypeSceneEnd
)

type Update struct{}

func (msg Update) Type() MessageType { return TypeUpdate }

type AddToScene struct{ Scene *tetra3d.Scene }

func (msg AddToScene) Type() MessageType { return TypeAddToScene }

type RemoveFromScene struct{ Scene *tetra3d.Scene }

func (msg RemoveFromScene) Type() MessageType { return TypeRemoveFromScene }

type SceneStart struct{ Scene *tetra3d.Scene }

func (msg SceneStart) Type() MessageType { return TypeSceneStart }

type SceneEnd struct{ Scene *tetra3d.Scene }

func (msg SceneEnd) Type() MessageType { return TypeSceneEnd }

package messages

import "github.com/solarlune/tetra3d"

// MessageType is the type of message dispatched; it can be bitwise combined together using standard addition for subscriptions.
type MessageType uint64

func (msg MessageType) Contains(other MessageType) bool {
	return msg&other > 0
}

// ISubscriber indicates an object subscribes to only a subset of all received Messages.
// The Subscribe() function returns the MessageType or MessageTypes (added together) that are desired.
// If no Subscribe() function is defined (so the object does not fulfill ISubscriber), the object receives all MessageTypes.
type ISubscriber interface {
	Subscribe() MessageType // Subscribe returns the MessageTypes (added together) that the IReceiver takes.
}

// IMessage indicates a contract for messages. Messages can have additional fields individually, but they all must return a bitwise MessageType.
type IMessage interface {
	Type() MessageType
}

const (
	TypeUpdate          MessageType = 1 << iota // 1
	TypeAddToScene                              // 2
	TypeRemoveFromScene                         // 4
	TypeSceneStart                              // 8
	TypeSceneEnd                                // 16
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

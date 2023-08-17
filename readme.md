# TetraFlow

TetraFlow is a package specifically for creating a simple object-oriented game engine-like flow for [Tetra3D](https://github.com/solarlune/tetra3d) projects.

The general approach is to load in your Tetra3D library (a collection of Scenes), create an `*tetraflow.Engine`, and then register a function to create game objects (`IReceivers`) for each relevant node. TetraFlow will register these receivers in the connected nodes' `INode.Data()` slots. You then create `Stage`s for each Scene in your Scene stack (e.g. one for your GUI, one for your game, and one for your background skybox), setting relevant scenes in each Stage. When you call `Engine.Update()`, the `Engine` object will continuously message IReceivers on engine events (when a node is added or removed from the scene tree, when the engine updates, etc) for all Stages.

To summarize, the Engine updates Stages, Stages contain Tetra3D Scenes, and you can register game logic structs to be IReceivers to listen for the Engine's updates.

# How to get it

```go get github.com/solarlune/tetraflow```

# How to use it

```

func (g *Game) Init() {

    g.Engine = tetraflow.NewEngine(
        
        g.Library,

        func(node tetra3d.INode) tetraflow.IReceiver {

            // Create relevant engine receivers for a game object depending on what logic the game object should have
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

    gameStage := g.Engine.AddStage("Game")
    gameStage.SetScene("Level 1")

}

func (g *Game) Update() error {

    g.Engine.Update()

    return nil

}

```

## Messages

TetraFlow's engine flow is based around the concept of message passing.

The general idea is to have your game objects fulfill the `tetraflow.IReceiver` interface, which just defines a single function - `ReceiveMessage(tetraflow.IMessage)`. By declaring this function, your game object struct becomes something that can receive messages from the `Engine` instance.

The `Engine` will send a variety of messages as nodes live in a scene's node hierarchy. These messages can indicate changes in the node or engine's life-cycle; for example, a message could indicate that the node is being updated (once per game frame), removed from the scene tree, or that the scene itself has just begun or just ended execution. By listening for these messages, GameObjects can behave appropriately:

```go

type Player struct {
    Node tetra3d.INode
}

func NewPlayer(node tetra3d.INode) *Player {
    return &Player{ Node: node }
}

// This function allows Player to become an IReceiver.
func (player *Player) ReceiveMessage(msg tetraflow.IMessage) {

    switch message := msg.(type) {
        case tetraflow.MessageOnUpdate:
            // Update message, happens once per game-frame
        case tetraflow.MessageOnAdd:
            // OnAdd, called once when the node is added to the scene tree.
    }

    if msg.Type() == tetraflow.MessageTypeOnUpdate {
        // Update message, happens every game frame
    }

}

```

You can switch off against either `IMessage.Type()`, or by the type of the message itself (which is the more elegant approach, as this allows you to access custom message values like `MessageOnSceneEnd.Scene`.)

By using messages, one can make their own custom messages to pass around, by simply making a custom struct that implements `IMessage`.

# Todo

## Functionality

- [ ] - Message subscription so all Receivers aren't receiving all Messages

## Message types

- [x] - On Scene Start / End (`OnSceneStart()`)
- [ ] - Collision checking (`OnCollision(other IBoundingObject)`)?
- [ ] - Input messages?
- [ ] - Scene Tree change
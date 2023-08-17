# TetraFlow 🚗💨

TetraFlow is a package specifically for creating a simple object-oriented game engine-like flow for [Tetra3D](https://github.com/solarlune/tetra3d) projects.

# How does it work

The general approach is to load in your Tetra3D library (a collection of Scenes), create an `*tetraflow.Engine`, and then register a function to create game objects (`IReceivers`) for each relevant node. TetraFlow will register these receivers in the connected nodes' `INode.Data()` slots. You then create `Stage`s for each Scene in your Scene stack (e.g. one for your GUI, one for your game, and one for your background skybox), setting relevant scenes in each Stage. When you call `Engine.Update()`, the `Engine` object will continuously message IReceivers on engine events (when a node is added or removed from the scene tree, when the engine updates, etc) for all Stages.

To summarize, the Engine updates Stages, Stages contain Tetra3D Scenes, and you can create game logic structs to be IReceivers, which listen for the Engine's updates.

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

    // Engine.Update() returns an error if you call Engine.Quit().
    return g.Engine.Update()
}

```

## Messages

TetraFlow's engine flow is based around the concept of message passing.

The general idea is to have your game objects fulfill the `tetraflow.IReceiver` interface, which just defines a single function - `ReceiveMessage(messages.IMessage)`. By declaring this function, your game object struct becomes something that can receive messages from the `Engine` instance.

The `Engine` will send a variety of messages as nodes live in a scene's node hierarchy. These messages can indicate changes in the node or the scene's life-cycle; for example, a message could indicate that the node is being updated (once per game frame), removed from the scene tree, or that the scene itself has just begun or just ended execution. By listening for these messages, GameObjects can behave appropriately:

```go

type Player struct {
    Node tetra3d.INode
}

func NewPlayer(node tetra3d.INode) *Player {
    return &Player{ Node: node }
}

// This function allows Player to become an IReceiver.
func (player *Player) ReceiveMessage(msg messages.IMessage) {

    switch message := msg.(type) {
        case messages.MessageOnUpdate:
            // Update message, happens once per game-frame
        case messages.MessageOnAdd:
            // OnAdd, called once when the node is added to the scene tree.
    }

    if msg.Type() == messages.MessageTypeOnUpdate {
        // Update message, happens every game frame
    }

}

```

You can switch off against either `IMessage.Type()`, or by the type of the message objects itself (`switch message := msg.(type)`) to distinguish between IMessages. The latter is the more elegant approach, as this allows you to access custom message values like `MessageOnSceneEnd.Scene`.

By using messages, one can also make their own custom messages to pass around, by simply making a custom struct that implements `IMessage`. You can send messages to target objects, or to the entire engine using `Engine.SendMessage()`.

## Message Subscription

Messages can be filtered down by making your object implement `ISubscriber`, which entails creating a `Subscribe() messages.MessageType` function. This function returns the MessageTypes that should be received for that object. If you want to subscribe to more than one message type, you can simply add them together, as the types are bitwise values.

```go

func (player *Player) Subscribe() messages.MessageType {
    return messages.TypeUpdate + messages.TypeSceneStart // Only subscribe to update messages or messages indicating the start of the current scene
}

```

# Todo

## Functionality

- [x] - Message subscription so all Receivers aren't receiving all Messages

## Message types

- [x] - On Scene Start / End (`OnSceneStart()`)
- [ ] - Collision checking (`OnCollision(other IBoundingObject)`)?
- [ ] - Input messages?
- [ ] - Scene Tree change
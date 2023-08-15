# TetraFlow

TetraFlow is a package specifically for creating a simple object-oriented game engine-like flow for [Tetra3D](https://github.com/solarlune/tetra3d) projects.

The general approach is to create a Tetra3D project, create an `*tetraflow.EngineFlow`, and then register a function to create game objects for each relevant node in Tetra3D. TetraFlow will register these game objects in the connected nodes as data pointers. Then call EngineFlow.Update(), and it will continuously call implemented `OnUpdate()` or `OnRemove()` functions as events happen.

# Todo

- [ ] - On Scene Start / Restart / End (`OnSceneStart()`)? (This may require some additional boilerplate code in Scene, I think...?)
- [ ] - Collision checking (`OnCollision()`)?
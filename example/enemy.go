package main

import (
	"github.com/solarlune/tetra3d"
	"github.com/solarlune/tetraflow/messages"
)

type Enemy struct {
	Model *tetra3d.Model
}

func NewEnemy(node tetra3d.INode) *Enemy {
	return &Enemy{Model: node.(*tetra3d.Model)}
}

func (enemy *Enemy) ReceiveMessage(msg messages.IMessage) {

	if msg.Type() == messages.TypeUpdate {

		player := enemy.Model.Scene().Root.SearchTree().ByProps("player").First()
		if player != nil {
			diff := player.WorldPosition().Sub(enemy.Model.WorldPosition())

			solids := enemy.Model.Scene().Root.SearchTree().IBoundingObjectsWithProps("solid")

			bounds := enemy.Model.Get("BoundingAABB").(*tetra3d.BoundingAABB)

			move := diff.Unit().Scale(0.05)

			bounds.CollisionTest(tetra3d.CollisionTestSettings{Others: solids, HandleCollision: func(col *tetra3d.Collision) bool {
				move = move.Add(col.AverageMTV())
				return false
			}})

			enemy.Model.MoveVec(move)
		}

	}

}

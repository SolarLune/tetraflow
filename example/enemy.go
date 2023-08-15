package main

import (
	"github.com/solarlune/tetra3d"
)

type Enemy struct {
	Model *tetra3d.Model
}

func NewEnemy(node tetra3d.INode) *Enemy {
	return &Enemy{Model: node.(*tetra3d.Model)}
}

func (enemy *Enemy) OnUpdate(dt float64) {

	player := enemy.Model.Scene().Root.SearchTree().ByName("Player").First()
	if player != nil {
		diff := player.WorldPosition().Sub(enemy.Model.WorldPosition())
		enemy.Model.MoveVec(diff.Unit().Scale(0.05))
	}
}

package mobile

import (
    "github.com/coding67-KR/ocean-catch/game"
    "github.com/hajimehoshi/ebiten/v2/mobile"
)

func init(){ mobile.SetGame(game.NewGame()) }
func Dummy(){}

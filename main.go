package main

import (
    "math/rand"
    "time"
    "github.com/coding67-KR/ocean-catch/game"
    "github.com/hajimehoshi/ebiten/v2"
)

func main(){
    ebiten.SetWindowSize(850,500)
    ebiten.SetWindowTitle("Ocean Catch v1")
    rand.Seed(time.Now().UnixNano())
    if err:=ebiten.RunGame(game.NewGame()); err!=nil { panic(err) }
}

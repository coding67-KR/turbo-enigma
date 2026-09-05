package game

import (
    "fmt"
    "image/color"
    "math/rand"
    "time"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Fish struct { Name string; Value int; Rarity int }
type Spot struct { Name string; Bonus int }
type Game struct { Money int; SpotIndex int; CatchCount int; Caught [60]int; Aquarium [60]int; LastCatch string; LastSave time.Time; State string }

var FishList = func() [60]Fish { names := [60]string{"Anchovy","Sardine","Mackerel","Horse Mackerel","Herring","Scomber","Flying Fish","Squid","Cuttlefish","Octopus","Flounder","Sole","Rockfish","Pollock","Cod","Haddock","Halibut","Red Snapper","Sea Bream","Yellowtail","Amberjack","Bonito","Tuna","Skipjack Tuna","Albacore","Bluefin Tuna","Swordfish","Marlin","Wahoo","Barracuda","Grouper","Snapper","Giant Trevally","Trevally","Mahi-Mahi","Parrotfish","Butterflyfish","Clownfish","Angelfish","Lionfish","Pufferfish","Porcupinefish","Moray Eel","Conger Eel","Stingray","Eagle Ray","Manta Ray","Hammerhead Shark","Tiger Shark","Great White Shark","Whale Shark","Nautilus","Sea Turtle","Sunfish","Coelacanth","Oarfish","Goblin Shark","Giant Squid","Blue Dragon","Golden Koi"}; var a [60]Fish; for i,n := range names { a[i]=Fish{Name:n,Value:10+i*7,Rarity:i%6} }; return a }()
var Spots = [3]Spot{{"Sunny Coast",0},{"Coral Bay",12},{"Deep Trench",30}}

func NewGame() *Game { return &Game{State:"main", LastSave:time.Now()} }
func (g *Game) Update() error {
    x,y := ebiten.CursorPosition()
    if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) { return nil }
    if g.State != "main" { if y < 70 { g.State="main" }; return nil }
    switch {
    case x < 210 && y >= 105 && y < 171: g.SpotIndex=(y-105)/22
    case y >= 250 && y < 345: g.catchFish()
    case x >= 20 && x < 210 && y >= 360 && y < 430: g.sellAll()
    case x >= 225 && x < 415 && y >= 360 && y < 430: g.State="aquarium"
    case x >= 430 && x < 620 && y >= 360 && y < 430: g.State="codex"
    }
    return nil
}
func (g *Game) TickSave() { if time.Since(g.LastSave) >= 700*time.Millisecond { g.LastSave=time.Now() } }
func (g *Game) catchFish() { idx:=rand.Intn(60); if Spots[g.SpotIndex].Bonus>0 { idx=(idx+rand.Intn(15)+Spots[g.SpotIndex].Bonus/10)%60 }; g.Caught[idx]++; g.Aquarium[idx]++; g.CatchCount++; g.LastCatch=fmt.Sprintf("%s  +$%d",FishList[idx].Name,FishList[idx].Value) }
func (g *Game) sellAll() { total:=0; for i:=range FishList { total += g.Caught[i]*FishList[i].Value; g.Caught[i]=0 }; g.Money += total; g.LastCatch=fmt.Sprintf("Sold catch for $%d",total) }
func (g *Game) Draw(s *ebiten.Image) {
    s.Fill(color.RGBA{12,55,82,255})
    if g.State=="aquarium" { ebitenutil.DebugPrintAt(s,"AQUARIUM  | tap top to return",20,25); y:=60; for i:=range FishList { if g.Aquarium[i]>0 { ebitenutil.DebugPrintAt(s,fmt.Sprintf("%-18s x%d",FishList[i].Name,g.Aquarium[i]),20,y); y+=18; if y>470 { break } } }; return }
    if g.State=="codex" { ebitenutil.DebugPrintAt(s,"FISH CODEX  | tap top to return",20,25); y:=60; for i:=range FishList { mark:="?"; if g.Aquarium[i]>0 { mark="OK" }; ebitenutil.DebugPrintAt(s,fmt.Sprintf("%2d %-18s %s",i+1,FishList[i].Name,mark),20,y); y+=16; if y>470 { break } }; return }
    ebitenutil.DebugPrintAt(s,"OCEAN CATCH v1",20,18); ebitenutil.DebugPrintAt(s,fmt.Sprintf("Money: $%d   Caught: %d/60",g.Money,g.CatchCount),20,48); ebitenutil.DebugPrintAt(s,"Fishing spots",20,82)
    for i,p:=range Spots { mark:=""; if i==g.SpotIndex { mark=" <" }; ebitenutil.DebugPrintAt(s,fmt.Sprintf("[%d] %s%s",i+1,p.Name,mark),25,105+i*22) }
    ebitenutil.DrawRect(s,20,250,805,90,color.RGBA{8,40,62,255}); ebitenutil.DebugPrintAt(s,"TAP HERE TO CATCH",40,278); ebitenutil.DebugPrintAt(s,"Latest: "+g.LastCatch,40,310)
    buttons:=[]struct{x float64;t string}{{20,"SELL"},{225,"AQUARIUM"},{430,"CODEX"}}; for _,b:=range buttons { ebitenutil.DrawRect(s,b.x,360,190,70,color.RGBA{34,116,162,255}); ebitenutil.DebugPrintAt(s,b.t,int(b.x)+55,390) }
}
func (g *Game) Layout(w,h int)(int,int){return 850,500}

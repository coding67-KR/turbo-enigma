package main

import (
    "fmt"
    "image/color"
    "math/rand"
    "time"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Fish struct { Name string; Value int; Rarity int }
type Spot struct { Name string; Bonus int }

type Game struct {
    money int
    spot int
    catchCount int
    caught [60]int
    aquarium [60]int
    lastCatch string
    lastSave time.Time
    state string
}

var fish = func() [60]Fish {
    names := [60]string{"Anchovy","Sardine","Mackerel","Horse Mackerel","Herring","Scomber","Flying Fish","Squid","Cuttlefish","Octopus","Flounder","Sole","Rockfish","Pollock","Cod","Haddock","Halibut","Red Snapper","Sea Bream","Yellowtail","Amberjack","Bonito","Tuna","Skipjack Tuna","Albacore","Bluefin Tuna","Swordfish","Marlin","Wahoo","Barracuda","Grouper","Snapper","Giant Trevally","Trevally","Mahi-Mahi","Parrotfish","Butterflyfish","Clownfish","Angelfish","Lionfish","Pufferfish","Porcupinefish","Moray Eel","Conger Eel","Stingray","Eagle Ray","Manta Ray","Hammerhead Shark","Tiger Shark","Great White Shark","Whale Shark","Nautilus","Sea Turtle","Sunfish","Coelacanth","Oarfish","Goblin Shark","Giant Squid","Blue Dragon","Golden Koi"}
    var a [60]Fish
    for i,n := range names { a[i]=Fish{Name:n,Value:10+i*7,Rarity:i%6} }
    return a
}()

var spots = [3]Spot{{"Sunny Coast",0},{"Coral Bay",12},{"Deep Trench",30}}

func (g *Game) Update() error {
    x,y := ebiten.CursorPosition()
    clicked := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
    if clicked && x>=20 && x<210 && y>=360 && y<430 { g.catchFish() }
    if clicked && x>=225 && x<415 && y>=360 && y<430 { g.sellAll() }
    if clicked && x>=430 && x<620 && y>=360 && y<430 { g.state="aquarium" }
    if clicked && x>=635 && x<825 && y>=360 && y<430 { g.state="codex" }
    if g.state!="main" && clicked && y<70 { g.state="main" }
    if time.Since(g.lastSave) >= 700*time.Millisecond { g.lastSave=time.Now(); g.save() }
    return nil
}

func (g *Game) catchFish() {
    if g.state!="main" { return }
    idx := rand.Intn(60)
    // deeper spots improve the chance of rarer fish
    if spots[g.spot].Bonus > 0 { idx=(idx+rand.Intn(15)+spots[g.spot].Bonus/10)%60 }
    g.caught[idx]++; g.aquarium[idx]++; g.catchCount++
    g.lastCatch = fmt.Sprintf("%s  +$%d",fish[idx].Name,fish[idx].Value)
}
func (g *Game) sellAll() { total:=0; for i:=range fish { total+=g.caught[i]*fish[i].Value; g.caught[i]=0 }; g.money+=total; g.lastCatch=fmt.Sprintf("Sold catch for $%d",total) }
func (g *Game) save() { /* state is intentionally kept live; mobile build can add persistent file bridge */ }

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(color.RGBA{12,55,82,255})
    if g.state=="aquarium" { g.drawAquarium(screen); return }
    if g.state=="codex" { g.drawCodex(screen); return }
    ebitenutil.DebugPrintAt(screen,"OCEAN CATCH v1",20,18)
    ebitenutil.DebugPrintAt(screen,fmt.Sprintf("Money: $%d    Fish caught: %d",g.money,g.catchCount),20,48)
    ebitenutil.DebugPrintAt(screen,"Fishing Spot:",20,82)
    for i,s:=range spots { label:=fmt.Sprintf("[%d] %s",i+1,s.Name); if i==g.spot { label += "  <" }; ebitenutil.DebugPrintAt(screen,label,25,105+i*22) }
    // tap zones are also clickable buttons
    ebitenutil.DrawRect(screen,20,250,805,90,color.RGBA{8,40,62,255})
    ebitenutil.DebugPrintAt(screen,"CATCH!   Tap anywhere in this area",40,275)
    ebitenutil.DebugPrintAt(screen,"Latest: "+g.lastCatch,40,310)
    ebitenutil.DrawRect(screen,20,360,190,70,color.RGBA{34,116,162,255}); ebitenutil.DebugPrintAt(screen,"SELL",90,390)
    ebitenutil.DrawRect(screen,225,360,190,70,color.RGBA{34,116,162,255}); ebitenutil.DebugPrintAt(screen,"AQUARIUM",275,390)
    ebitenutil.DrawRect(screen,430,360,190,70,color.RGBA{34,116,162,255}); ebitenutil.DebugPrintAt(screen,"CODEX",495,390)
    ebitenutil.DebugPrintAt(screen,"Tap spot 1/2/3 by touching the upper left list.",20,455)
}
func (g *Game) drawAquarium(screen *ebiten.Image){ ebitenutil.DebugPrintAt(screen,"AQUARIUM  (tap top to return)",20,25); y:=60; for i:=range fish { if g.aquarium[i]>0 { ebitenutil.DebugPrintAt(screen,fmt.Sprintf("%-18s x%d",fish[i].Name,g.aquarium[i]),20,y); y+=18; if y>470 { break } } } }
func (g *Game) drawCodex(screen *ebiten.Image){ ebitenutil.DebugPrintAt(screen,"FISH CODEX  (tap top to return)",20,25); y:=60; for i:=range fish { mark:="?"; if g.aquarium[i]>0 { mark="✓" }; ebitenutil.DebugPrintAt(screen,fmt.Sprintf("%2d %-18s %s",i+1,fish[i].Name,mark),20,y); y+=16; if y>470 { break } } }

func (g *Game) Layout(w,h int)(int,int){return 850,500}
func main(){ ebiten.SetWindowSize(850,500); ebiten.SetWindowTitle("Ocean Catch v1"); rand.Seed(time.Now().UnixNano()); if err:=ebiten.RunGame(&Game{state:"main"}); err!=nil { panic(err) } }

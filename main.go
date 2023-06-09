package main

import (
	"fmt"
	"image"
	"math/rand"

	//"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Position struct {
	x float64
	y float64
}

type Screen struct {
	w int
	h int
}

type Player struct {
	sprite *ebiten.Image
	position Position
	speed float64
}

type Food struct {
	sprite *ebiten.Image
	position Position
	eaten bool
}

type Game struct{
	//image *ebiten.Image
	player Player
	//position struct {
	//	x float64
	//	y float64
	//}
	screenSize Screen
	food Food
	score int
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("killing game")
	}

	if g.food.eaten {
		fmt.Println("Placing the food somewhere")
		g.food.position.x = 0 + rand.Float64() * float64(g.screenSize.w)
		g.food.position.y = 0 + rand.Float64() * float64(g.screenSize.h)
		g.food.eaten = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.position.y = g.player.position.y - g.player.speed
		g.player.position.x = g.player.position.x - g.player.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.position.y = g.player.position.y - g.player.speed
		g.player.position.x = g.player.position.x + g.player.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.position.y = g.player.position.y + g.player.speed
		g.player.position.x = g.player.position.x - g.player.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.position.y = g.player.position.y + g.player.speed
		g.player.position.x = g.player.position.x + g.player.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.player.position.y = g.player.position.y - g.player.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.player.position.y = g.player.position.y + g.player.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.player.position.x = g.player.position.x - g.player.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.player.position.x = g.player.position.x + g.player.speed
	}


	if g.player.position.x < 0 {
		g.player.position.x = g.player.position.x * -1
	}
	if g.player.position.x >= float64(g.screenSize.w) {
		g.player.position.x = float64(g.screenSize.w) - float64(g.player.sprite.Bounds().Dx())
		fmt.Printf("Out of bounds x: %f\n", g.player.position.x)
	}
	if g.player.position.y < 0 {
		g.player.position.y = g.player.position.y * -1
	}
	if g.player.position.y >= float64(g.screenSize.h) {
		g.player.position.y = float64(g.screenSize.h) - float64(g.player.sprite.Bounds().Dy())
		fmt.Printf("Out of bounds y: %f\n", g.player.position.y)
	}

	foodRect := image.Rect(int(g.food.position.x), int(g.food.position.y), int(g.food.position.x) + g.food.sprite.Bounds().Dx(), int(g.food.position.y) + g.food.sprite.Bounds().Dy())
	playerRect := image.Rect(int(g.player.position.x), int(g.player.position.y), int(g.player.position.x) + g.player.sprite.Bounds().Dx(), int(g.player.position.y) + g.player.sprite.Bounds().Dy())

	if foodRect.In(playerRect) {
		fmt.Println("Food is in player")
		g.score += 1
		g.food.eaten = true
	}

	//fmt.Printf("Character: {%f, %f}\n", g.player.position.x, g.player.position.y)
	//newX := 0 + rand.Float64() * (float64(g.screenSize.w)-0)
	//newY := 0 + rand.Float64() * (float64(g.screenSize.h)-0)

	//g.position.x += newX
	//g.position.y += newY

	//g.player.position.x = newX
	//g.player.position.y = newY

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	scoreString := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrint(screen, scoreString)
	//ebitenutil.DebugPrint(screen, "Hello World!")
	//img := ebiten.NewImage(10, 20)
	//img.Fill(color.RGBA{R: 255, A: 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.player.position.x, g.player.position.y)

	foodop := &ebiten.DrawImageOptions{}
	foodop.GeoM.Translate(g.food.position.x, g.food.position.y)

	screen.DrawImage(g.player.sprite, op)
	screen.DrawImage(g.food.sprite, foodop)
	//fmt.Printf("Tick")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	//return 320, 240
	return outsideWidth, outsideHeight
}

func main() {
	//ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hellow World!")
	image := ebiten.NewImage(20,20)
	image.Fill(color.RGBA{R: 255, A: 255})

	screen := &Screen{
		w: 640,
		h: 480,
	}

	playerInitialPosition := &Position{
		x: float64(screen.w / 2) - float64(image.Bounds().Dx() / 2),
		y: float64(screen.h / 2) - float64(image.Bounds().Dy() / 2),
	}

	player := &Player {
		sprite: image,
		position: *playerInitialPosition,
		speed: 10.0,
	}

	foodImage := ebiten.NewImage(5,5)
	foodImage.Fill(color.RGBA{R:255, G:255, B:255, A:255})

	food := &Food {
		sprite: foodImage,
		eaten: true,
	}

	game := &Game{
		//image: image,
		player: *player,
		screenSize: *screen,
		food: *food,
	}
	//game.screenSize.w = 640
	//game.screenSize.h = 480
	ebiten.SetWindowSize(game.screenSize.w, game.screenSize.h)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type MainMenuScene struct{}

type Position struct {
	x float64
	y float64
}

type Screen struct {
	w int
	h int
}

type Player struct {
	sprite   *ebiten.Image
	position Position
	speed    float64
}

type Food struct {
	sprite   *ebiten.Image
	position Position
	eaten    bool
}

type GameScene struct {
	player     Player
	screenSize Screen
	food       Food
	score      int
}

type Game struct {
	CurrentScene Scene
}

var (
	game       *Game
	gameScreen *Screen
)

func (m *MainMenuScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		playerImage := ebiten.NewImage(20, 20)
		playerImage.Fill(color.RGBA{R: 255, A: 255})
		playerInitialPosition := &Position{
			x: float64(gameScreen.w/2) - float64(playerImage.Bounds().Dx()/2),
			y: float64(gameScreen.h/2) - float64(playerImage.Bounds().Dy()/2),
		}

		player := &Player{
			sprite:   playerImage,
			position: *playerInitialPosition,
			speed:    10.0,
		}

		foodImage := ebiten.NewImage(5, 5)
		foodImage.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})

		food := &Food{
			sprite: foodImage,
			eaten:  true,
		}

		game.CurrentScene = &GameScene{
			player:     *player,
			screenSize: *gameScreen,
			food:       *food,
		}
	}
	return nil
}

func (m *MainMenuScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Press Enter to play", gameScreen.w/2, gameScreen.h/2)
}

func forcePlayerInBounds(p Player, s Screen) Player {
	if p.position.x < 0 {
		p.position.x = p.position.x * -1
	}
	if p.position.x >= float64(s.w) {
		p.position.x = float64(s.w) - float64(p.sprite.Bounds().Dx())
		fmt.Printf("Out of bounds x: %f\n", p.position.x)
	}
	if p.position.y < 0 {
		p.position.y = p.position.y * -1
	}
	if p.position.y >= float64(s.h) {
		p.position.y = float64(s.h) - float64(p.sprite.Bounds().Dy())
		fmt.Printf("Out of bounds y: %f\n", p.position.y)
	}
	return p
}

func updatePlayerPosition(p Player) Player {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.position.y = p.position.y - p.speed
		p.position.x = p.position.x - p.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.position.y = p.position.y - p.speed
		p.position.x = p.position.x + p.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.position.y = p.position.y + p.speed
		p.position.x = p.position.x - p.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.position.y = p.position.y + p.speed
		p.position.x = p.position.x + p.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.position.y = p.position.y - p.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.position.y = p.position.y + p.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.position.x = p.position.x - p.speed
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.position.x = p.position.x + p.speed
	}

	return p
}

func (g *GameScene) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("killing game")
	}

	if g.food.eaten {
		fmt.Println("Placing the food somewhere")
		g.food.position.x = 0 + rand.Float64()*float64(g.screenSize.w)
		g.food.position.y = 0 + rand.Float64()*float64(g.screenSize.h)
		g.food.eaten = false
	}

	g.player = updatePlayerPosition(g.player)

	g.player = forcePlayerInBounds(g.player, g.screenSize)

	foodRect := image.Rect(int(g.food.position.x), int(g.food.position.y), int(g.food.position.x)+g.food.sprite.Bounds().Dx(), int(g.food.position.y)+g.food.sprite.Bounds().Dy())
	playerRect := image.Rect(int(g.player.position.x), int(g.player.position.y), int(g.player.position.x)+g.player.sprite.Bounds().Dx(), int(g.player.position.y)+g.player.sprite.Bounds().Dy())

	if foodRect.In(playerRect) {
		fmt.Println("Food is in player")
		g.score += 1
		g.food.eaten = true
	}

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	scoreString := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrint(screen, scoreString)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.player.position.x, g.player.position.y)

	foodop := &ebiten.DrawImageOptions{}
	foodop.GeoM.Translate(g.food.position.x, g.food.position.y)

	screen.DrawImage(g.player.sprite, op)
	screen.DrawImage(g.food.sprite, foodop)
}

func (g *Game) Update() error {
	return g.CurrentScene.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.CurrentScene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowTitle("Go Game")

	gameScreen = &Screen{
		w: 640,
		h: 480,
	}

	ebiten.SetWindowSize(gameScreen.w, gameScreen.h)

	game = &Game{
		CurrentScene: &MainMenuScene{},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

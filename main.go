package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

//type Button struct {
//	sprite *ebiten.Image
//}

type MainMenuScene struct {
	//buttons           []Button
	buttons           []string
	highlightedButton int
}

type EndGameScene struct {
	FinalScore int
}

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

type Mine struct {
	sprite   *ebiten.Image
	position Position
}

type GameScene struct {
	player     Player
	screenSize Screen
	food       Food
	mines      []Mine
	score      int
	StartTime  time.Time
	EndTime    time.Time
}

type Game struct {
	CurrentScene Scene
}

var (
	game       *Game
	gameScreen *Screen
)

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

func initGame() error {
	playerImage := ebiten.NewImage(20, 20)
	playerImage.Fill(color.RGBA{B: 255, A: 255})
	playerInitialPosition := &Position{
		x: float64(gameScreen.w/2) - float64(playerImage.Bounds().Dx()/2),
		y: float64(gameScreen.h/2) - float64(playerImage.Bounds().Dy()/2),
	}

	player := &Player{
		sprite:   playerImage,
		position: *playerInitialPosition,
		speed:    5.0,
	}

	foodImage := ebiten.NewImage(5, 5)
	foodImage.Fill(color.RGBA{R: 255, G: 255, B: 255, A: 255})

	food := &Food{
		sprite: foodImage,
		eaten:  true,
	}

	var emptyMines []Mine

	now := time.Now()

	game.CurrentScene = &GameScene{
		player:     *player,
		screenSize: *gameScreen,
		food:       *food,
		mines:      emptyMines,
		score:      0,
		StartTime:  now,
		EndTime:    now.Add(60 * time.Second),
	}
	return nil
}

func addMine(currentMines []Mine) []Mine {
	mineImage := ebiten.NewImage(5, 5)
	mineImage.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})

	minePosition := &Position{
		x: 0 + rand.Float64()*float64(gameScreen.w),
		y: 0 + rand.Float64()*float64(gameScreen.h),
	}

	newMine := &Mine{
		sprite:   mineImage,
		position: *minePosition,
	}

	currentMines = append(currentMines, *newMine)
	return currentMines
}

func (m *MainMenuScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return fmt.Errorf("killing game")
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		if m.highlightedButton > 0 {
			m.highlightedButton -= 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		if m.highlightedButton < len(m.buttons)-1 {
			m.highlightedButton += 1
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// Hard coding this for now
		if m.highlightedButton == 0 {
			return initGame()
		} else if m.highlightedButton == 1 {
			return fmt.Errorf("killing game")
		}
	}
	return nil
}

func (m *MainMenuScene) Draw(screen *ebiten.Image) {
	// To place all the buttons:
	// Imagine a rectangle comprising of all the buttons. This rectangle has a height equal to:
	// The height of each button, multiplied by the number of buttons, plus the amount of padding between each button, multiplied by the number of buttons minus 1

	// Using constants for now
	buttonHeight := 45
	buttonWidth := 150
	buttonPadding := 20
	highlightSize := 4

	buttonGroupArea := (buttonHeight * len(m.buttons)) + (buttonPadding * (len(m.buttons) - 1))
	topOfButtonGroup := (gameScreen.h / 2) - (buttonGroupArea / 2)

	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
	fontFace, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	for i, v := range m.buttons {
		buttonX := (gameScreen.w / 2) - (buttonWidth / 2)
		buttonY := topOfButtonGroup + ((buttonPadding + buttonHeight) * i)
		if i == m.highlightedButton {
			highlightButton := ebiten.NewImage(buttonWidth+(highlightSize*2), buttonHeight+(highlightSize*2))
			highlightButton.Fill(color.RGBA{R: 255, G: 255, B: 0, A: 255})
			highlightOptions := &ebiten.DrawImageOptions{}
			highlightX := buttonX - 2
			highlightY := buttonY - 2
			highlightOptions.GeoM.Translate(float64(highlightX), float64(highlightY))
			screen.DrawImage(highlightButton, highlightOptions)
		}
		ebitenutil.DrawRect(screen, float64(buttonX), float64(buttonY), float64(buttonWidth), float64(buttonHeight), color.White)
		text.Draw(screen, v, fontFace, buttonX+(buttonWidth/2)-(font.MeasureString(fontFace, v).Round()/2), buttonY+(buttonHeight/2)+fontFace.Metrics().Ascent.Round()/2, color.Black)
	}
}

func (e *EndGameScene) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return fmt.Errorf("killing game")
	} else if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		return initGame()
	}
	return nil
}

func (e *EndGameScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Game Over", gameScreen.w/2, gameScreen.h/2-30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Final Score: %d", e.FinalScore), gameScreen.w/2, gameScreen.h/2+10)
	ebitenutil.DebugPrintAt(screen, "Press Enter to play again", gameScreen.w/2, gameScreen.h/2+100)
}

func (g *GameScene) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("killing game")
	}

	if time.Until(g.EndTime) <= 0 {
		game.CurrentScene = &EndGameScene{
			FinalScore: g.score,
		}
	}

	if g.food.eaten {
		fmt.Println("Placing the food somewhere")
		g.food.position.x = 0 + rand.Float64()*float64(g.screenSize.w)
		g.food.position.y = 0 + rand.Float64()*float64(g.screenSize.h)
		g.food.eaten = false
	}

	numberOfMines := len(g.mines)
	roundedTime := int(time.Since(g.StartTime).Seconds())
	if (roundedTime / 10) > numberOfMines {
		fmt.Printf("RoundedTime: %d\n", roundedTime)
		g.mines = addMine(g.mines)
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

	for _, v := range g.mines {
		mineRect := image.Rect(int(v.position.x), int(v.position.y), int(v.position.x)+v.sprite.Bounds().Dx(), int(v.position.y)+v.sprite.Bounds().Dy())
		if mineRect.In(playerRect) {
			fmt.Println("Player hit a mine")
			game.CurrentScene = &EndGameScene{
				FinalScore: g.score,
			}
		}
	}

	return nil
}

func (g *GameScene) Draw(screen *ebiten.Image) {
	scoreString := fmt.Sprintf("Score: %d", g.score)
	ebitenutil.DebugPrintAt(screen, scoreString, 0, 0)
	timeString := fmt.Sprintf("Timer: %s", time.Until(g.EndTime).Truncate(time.Second).String())
	ebitenutil.DebugPrintAt(screen, timeString, 0, 30)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.player.position.x, g.player.position.y)

	foodop := &ebiten.DrawImageOptions{}
	foodop.GeoM.Translate(g.food.position.x, g.food.position.y)

	for _, v := range g.mines {
		mop := &ebiten.DrawImageOptions{}
		mop.GeoM.Translate(v.position.x, v.position.y)
		screen.DrawImage(v.sprite, mop)
	}

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

	var buttons []string
	buttons = append(buttons, "Start Game")
	buttons = append(buttons, "Exit Game")

	game = &Game{
		CurrentScene: &MainMenuScene{
			buttons:           buttons,
			highlightedButton: 0,
		},
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

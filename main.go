package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	fontFace font.Face
)

func LoadFont() {
	fontData, err := os.ReadFile("sdf.ttf")
	if err != nil {
		log.Fatalf("Failed to decode base64 font data: %v", err)
	} // Replace with your font file
	if err != nil {
		log.Fatalf("Failed to read font file: %v", err)
	}

	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatalf("Failed to parse font: %v", err)
	}

	fontFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    50,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("Failed to create font face: %v", err)
	}
}

var black color.Color
var op *ebiten.DrawImageOptions
var square *ebiten.Image

type Mode int

type Game struct {
	mode         Mode
	circleSpeed  pos
	touchIDs     []ebiten.TouchID
	intialSpeed  int
	acceralate   int
	score        int
	highestScore int
}

const (
	startSpeed       = 1
	tileSize         = 32
	titleFontSize    = fontSize * 1.5
	fontSize         = 24
	smallFontSize    = fontSize / 2
	pipeWidth        = tileSize * 2
	pipeStartOffsetX = 8
	pipeIntervalX    = 8
	pipeGapY         = 5
	startAcceralate  = 2
)

var (
	gopherImage *ebiten.Image
	tilesImage  *ebiten.Image
)

type pos struct {
	x int
	y int
}

const (
	screenWidth  int = 640
	screenHeight int = 640
	squareSize   int = 100
)

const (
	ModeTitle Mode = iota
	ModeGame
	ModeGameOver
)

func (g *Game) Update() error {

	switch g.mode {

	case ModeTitle:
		g.circleSpeed.x = 0
		g.circleSpeed.y = screenWidth / 2
		if g.isKeyJustPressed(ebiten.KeySpace) {
			g.mode = ModeGame
		}

	case ModeGame:
		if g.isKeyJustPressed(ebiten.KeySpace) {
			if g.circleSpeed.x < (screenWidth/2+50-30) && g.circleSpeed.x > (screenWidth/2-50+30) {
				g.mode = ModeGame
				g.intialSpeed += g.acceralate
				g.score++
			} else {
				g.mode = ModeGameOver
			}
		} else {
			g.circleSpeed.x += g.intialSpeed
			g.circleSpeed.y = screenWidth / 2

			if g.circleSpeed.x > 640 {
				g.circleSpeed.x = 0
			}
		}

	case ModeGameOver:
		g.circleSpeed.x = g.circleSpeed.x
		g.circleSpeed.y = screenWidth / 2
		if g.isKeyJustPressed(ebiten.KeySpace) {
			g.intialSpeed = startSpeed
			g.mode = ModeGame
			g.circleSpeed.x = 0
			g.score = 0
		}
	}
	if g.score > g.highestScore {
		g.highestScore = g.score
	}

	return nil
}

func (g *Game) isKeyJustPressed(key ebiten.Key) bool {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		return true
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return true
	}
	return false
}
func drawText(screen *ebiten.Image, str string, x, y int, clr color.Color) {
	textDrawer := font.Drawer{
		Dst:  screen,
		Src:  image.NewUniform(clr),
		Face: fontFace,
		Dot:  fixed.P(x, y),
	}
	textDrawer.DrawString(str)
}
func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(black)
	square := ebiten.NewImage(squareSize, squareSize)
	square.Fill(color.RGBA{R: 255, G: 224, B: 189, A: 255})
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((screenWidth-squareSize)/2), float64((screenHeight-squareSize)/2))
	screen.DrawImage(square, op)

	vector.DrawFilledCircle(screen, float32(g.circleSpeed.x), float32(g.circleSpeed.y), 30, color.RGBA{R: 139, G: 0, B: 0, A: 255}, false)
	if g.mode == ModeTitle || g.mode == ModeGame {
		drawText(screen, fmt.Sprintf("%d", g.score), 500, 100, color.White)
	} else {
		if fontFace == nil {
			fmt.Println("Font not loaded!")
			return
		}

		drawText(screen, fmt.Sprintf("Game Over!, Score: %d", g.score), 50, 500, color.White)

	}

	drawText(screen, fmt.Sprintf("Highest Score: %d", g.highestScore), 10, 100, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func colorLoad() {
	black = color.RGBA64{R: 0, G: 0, B: 0, A: 0}
}

func squareLoad() {

}
func init() {
	colorLoad()

}

func main() {
	LoadFont()
	game := &Game{
		intialSpeed: startSpeed,
		acceralate:  startAcceralate,
		score:       0,
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Catch Me")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}

}

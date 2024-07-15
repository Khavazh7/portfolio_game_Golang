package main

import (
	"fmt"
	"image/color"
	_ "image/png" // Чтобы поддерживать формат PNG
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game структура для хранения объектов игрока и врага
type Game struct {
	player   *Player
	enemy    *Enemy
	stars    []Star
	gameOver bool
}

// Player структура для игрока
type Player struct {
	x, y  float64
	lives int
}

// Enemy структура для врага
type Enemy struct {
	x, y        float64 // координаты врага
	vx, vy      float64 // скорость по оси X и по оси Y
	moveCounter int     // счетчик движения в одном направлении
	moveLimit   int     // предел движения в одном направлении
}

// Star структура для звезды
type Star struct {
	x, y float64
}

// Переменные для игрока и врага
var player = Player{x: 320, y: 240}
var enemy = Enemy{x: 100, y: 100}

// Функция для инициализации звезд
func initStars(numStars int) []Star {
	stars := make([]Star, numStars)
	for i := 0; i < numStars; i++ {
		stars[i] = Star{
			x: float64(rand.Intn(640)),
			y: float64(rand.Intn(480)),
		}
	}
	return stars
}

// Обновление состояния игрока
func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && p.y > 0 {
		p.y -= 2
	} // кнопка вверх
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && p.y < 464 {
		p.y += 2
	} // кнопка вниз
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && p.x > 0 {
		p.x -= 2
	} // кнопка влево
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && p.x < 624 {
		p.x += 2
	} // кнопка вправо
}

// Обновление состояния врага
func (e *Enemy) Update(g *Game) {
	// Если игра не закончена, обновляем состояние врага
	if !g.gameOver {
		// Если счетчик движения достиг предела, меняем направление
		if e.moveCounter <= e.moveLimit {
			// Продолжаем движение в текущем направлении
			e.x += e.vx
			e.y += e.vy
			e.moveCounter++
		} else {
			e.vx = rand.Float64()*6 - 3 // случайная скорость по оси X от -3 до 3
			e.vy = rand.Float64()*6 - 3 // случайная скорость по оси Y от -3 до 3
			e.moveCounter = 0
			e.moveLimit = rand.Intn(5) + 10 // случайный предел движения от 10 до 19
		}

		// Ограничение движения в пределах экрана
		if e.x < 0 {
			e.x = 0
			e.vx = rand.Float64()*9 - 3 // случайная скорость по оси X от -3 до 3
		} else if e.x > 624 { // Учитываем ширину врага (16 пикселей)
			e.x = 624
			e.vx = rand.Float64()*9 - 3 // случайная скорость по оси X от -3 до 3
		}
		if e.y < 0 {
			e.y = 0
			e.vy = rand.Float64()*9 - 3 // случайная скорость по оси Y от -3 до 3
		} else if e.y > 464 { // Учитываем высоту врага (16 пикселей)
			e.y = 464
			e.vy = rand.Float64()*9 - 3 // случайная скорость по оси Y от -3 до 3
		}
	}
}

// Проверка столкновения игрока и врага
func (p *Player) CheckCollision(e *Enemy) bool {
	return p.x < e.x+16 && p.x+16 > e.x && p.y < e.y+16 && p.y+16 > e.y
}

// Отрисовка игрока
func (p *Player) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, p.x, p.y, 16, 16, color.White)
}

// Отрисовка врага
func (e *Enemy) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, e.x, e.y, 16, 16, color.RGBA{255, 0, 0, 255})
}

// Отрисовка звезд
func (g *Game) DrawStars(screen *ebiten.Image) {
	for _, star := range g.stars {
		ebitenutil.DrawRect(screen, star.x, star.y, 2, 2, color.White)
	}
}

// Отрисовка "GAME OVER!"
func (g *Game) DrawGameOver(screen *ebiten.Image) {
	screen.Fill(color.Black)
	ebitenutil.DebugPrint(screen, "GAME OVER!")
}

// Обновление состояния игры
func (g *Game) Update() error {
	// Если игра уже закончена, пропустить обновление состояния
	if g.gameOver {
		return nil
	}

	g.player.Update()
	g.enemy.Update(g) // Передаем g в функцию Update у врага

	// Проверка столкновения игрока и врага
	if g.player.CheckCollision(g.enemy) {
		g.player.lives--
		log.Println("Collision! Lives left:", g.player.lives)
		if g.player.lives <= 0 {
			g.gameOver = true
		} else {
			// Переместим врага в случайное положение после столкновения
			g.enemy.x = rand.Float64() * 640
			g.enemy.y = rand.Float64() * 480
		}
	}

	return nil
}

// Отрисовка игры
func (g *Game) Draw(screen *ebiten.Image) {
	if g.gameOver {
		g.DrawGameOver(screen)
		return
	}

	screen.Fill(color.Black)
	g.DrawStars(screen)
	g.player.Draw(screen)
	g.enemy.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Lives: %d", g.player.lives))
}

// Макетирование окна игры
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

// Главная функция
func main() {
	rand.Seed(time.Now().UnixNano())
	player := &Player{x: 320, y: 240, lives: 3}
	enemy := &Enemy{x: 100, y: 100}
	stars := initStars(100)
	game := &Game{player: player, enemy: enemy, stars: stars}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Random Direction Enemy with Lives")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

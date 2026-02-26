package main

import (
	_ "embed"
)

const (
	SCREEN_WIDTH  = 128
	SCREEN_HEIGHT = 64
	MAP_WIDTH     = 8
	MAP_HEIGHT    = 8
	FOV           = 60
	MAX_DEPTH     = 16

	LEFT     = 2
	BACKWARD = 5
	RIGHT    = 8
	FORWARD  = 4
	FIRE     = 9
)

// labyrinth map
var gameMap = [MAP_HEIGHT][MAP_WIDTH]int{
	{1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 1, 1, 0, 1, 0, 1},
	{1, 0, 1, 0, 0, 1, 0, 1},
	{1, 0, 0, 0, 1, 1, 0, 1},
	{1, 0, 1, 0, 0, 0, 0, 1},
	{1, 0, 1, 1, 1, 1, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1},
}

var screen [SCREEN_HEIGHT][SCREEN_WIDTH]bool

type Player struct {
	x, y   float32
	angle  float32
	moveX  float32
	moveY  float32
	rotate float32
}

type MazeGame struct {
	player        Player
	displayBuffer *DisplayBuffer
}

func NewMazeGame(buffer *DisplayBuffer) (*MazeGame, error) {
	return &MazeGame{
		player: Player{
			x:     1.5,
			y:     1.5,
			angle: 0,
		},
		displayBuffer: buffer,
	}, nil
}

/*	for {
		handleInput()
		updatePlayer()
		render()
		display.Display()

		time.Sleep(33 * time.Millisecond)
	}
*/

func (g *MazeGame) Update(jx, jy int16) {

	if jy <= -2 {
		g.player.moveX += cosf32(g.player.angle)
		g.player.moveY += sinf32(g.player.angle)
	}

	if jy >= 2 {
		g.player.moveX -= cosf32(g.player.angle)
		g.player.moveY -= sinf32(g.player.angle)
	}

	if jx <= -2 {
		g.player.rotate -= 1
	}

	if jx >= 2 {
		g.player.rotate += 1
	}

	newX := g.player.x + g.player.moveX*0.1
	newY := g.player.y + g.player.moveY*0.1

	if gameMap[int(newY)][int(g.player.x)] == 0 {
		g.player.y = newY
	}
	if gameMap[int(g.player.y)][int(newX)] == 0 {
		g.player.x = newX
	}

	g.player.angle += g.player.rotate * 5
	if g.player.angle >= 360 {
		g.player.angle -= 360
	}
	if g.player.angle < 0 {
		g.player.angle += 360
	}

	g.player.moveX *= 0.8
	g.player.moveY *= 0.8
	g.player.rotate *= 0.8

	for x := 0; x < SCREEN_WIDTH; x++ {
		rayAngle := g.player.angle + float32(x-SCREEN_WIDTH/2)*FOV/SCREEN_WIDTH
		distance := g.castRay(rayAngle)
		distance = distance * cosf32(rayAngle-g.player.angle)
		wallHeight := int(float32(SCREEN_HEIGHT) / distance)
		if wallHeight > SCREEN_HEIGHT {
			wallHeight = SCREEN_HEIGHT
		}
		wallTop := (SCREEN_HEIGHT - wallHeight) / 2
		wallBottom := wallTop + wallHeight
		g.drawVLine(x, wallTop, wallBottom, true)
	}
}

func sinf32(angle float32) float32 {
	angle = angle - float32(int(angle/360))*360
	if angle < 0 {
		angle += 360
	}

	switch {
	case angle < 90:
		return angle / 90
	case angle < 180:
		return (180 - angle) / 90
	case angle < 270:
		return -(angle - 180) / 90
	default:
		return -(360 - angle) / 90
	}
}

func cosf32(angle float32) float32 {
	return sinf32(angle + 90)
}

func sqrt(x float32) float32 {
	if x <= 0 {
		return 0
	}
	result := x
	for i := 0; i < 10; i++ {
		result = (result + x/result) / 2
	}
	return result
}

func (g *MazeGame) drawVLine(x, y1, y2 int, value bool) {
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for y := y1; y <= y2; y++ {
		if value {
			g.displayBuffer.SetPixel(int16(x), int16(y), textWhite)
		} else {
			g.displayBuffer.SetPixel(int16(x), int16(y), textWhite)
		}
	}
}

func (g *MazeGame) castRay(rayAngle float32) float32 {
	rayX := g.player.x
	rayY := g.player.y

	rayDirX := cosf32(rayAngle) * 0.1
	rayDirY := sinf32(rayAngle) * 0.1

	for i := 0; i < MAX_DEPTH*10; i++ {
		mapX := int(rayX)
		mapY := int(rayY)

		if mapX < 0 || mapX >= MAP_WIDTH || mapY < 0 || mapY >= MAP_HEIGHT {
			break
		}

		if gameMap[mapY][mapX] == 1 {
			dx := rayX - g.player.x
			dy := rayY - g.player.y
			return sqrt(dx*dx + dy*dy)
		}

		rayX += rayDirX
		rayY += rayDirY
	}

	return MAX_DEPTH
}

package main

import (
        "fmt"
        "math/rand"
        "time"
)

const (
        MapBoundChar  = '#'
        MapFillChar   = ' '
        SnakeBodyChar = '@'
        FoodChar      = '*'
)

type Vector2 struct {
        X int
        Y int
}

const (
        DirectionUp    = iota
        DirectionRight = iota
        DirectionDown  = iota
        DirectionLeft  = iota
)

func (v *Vector2) Move(direction int) {
        switch direction {
        case DirectionUp:
                v.Y--
        case DirectionRight:
                v.X++
        case DirectionDown:
                v.Y++
        case DirectionLeft:
                v.X--
        }
}

func (v *Vector2) ReverseMove(direction int) {
        switch direction {
        case DirectionDown:
                v.Y--
        case DirectionLeft:
                v.X++
        case DirectionUp:
                v.Y++
        case DirectionRight:
                v.X--
        }
}

type SnakeNode struct {
        Direction int
        PrevNode  *SnakeNode
        NextNode  *SnakeNode
        Position  Vector2
}

type Snake struct {
        Length int
        Head   *SnakeNode
        Tail   *SnakeNode
}

func NewSnake(length int) *Snake {
        newSnake := Snake{Length: 1}
        head := &SnakeNode{
                Direction: DirectionRight,
                Position: Vector2{
                        X: 1,
                        Y: 1,
                },
        }
        newSnake.Head = head
        newSnake.Tail = head
        newSnake.AddNodes(length)
        return &newSnake
}

func (s *Snake) AddNodes(length int) {
        for i := 0; i < length; i++ {
                newNodePosition := s.Tail.Position
                newNodePosition.ReverseMove(s.Tail.Direction)
                newNode := &SnakeNode{
                        Direction: s.Tail.Direction,
                        PrevNode:  s.Tail,
                        Position:  newNodePosition,
                }
                s.Tail.NextNode = newNode
                s.Tail = newNode
                s.Length++
        }
}

func (s *Snake) RemoveTailNode() {
        if s.Tail.PrevNode == nil {
                return
        }
        s.Tail.PrevNode.NextNode = nil
        s.Tail = s.Tail.PrevNode
        s.Length--
}

func (s *Snake) Update() {
        prevNodePosition := s.Head.Position
        prevNodeDirection := s.Head.Direction
        s.Head.Position.Move(s.Head.Direction)
        node := s.Head.NextNode
        for node != nil {
                prevPrevNodePosition := node.Position
                prevPrevNodeDirection := node.Direction
                node.Position = prevNodePosition
                node.Direction = prevNodeDirection
                prevNodePosition = prevPrevNodePosition
                prevNodeDirection = prevPrevNodeDirection
                node = node.NextNode
        }
}

func main() {
        rand.Seed(time.Now().UnixMilli())

        mapSize := Vector2{50, 20}
        mapDisplayBuffer := make([]byte, (mapSize.X+3)*(mapSize.Y+2))

        snake := NewSnake(5)

        food := GenerateNewFood(mapSize, 1)

        var totalScore, maxScore, loses int

        t := time.NewTicker(time.Second / 10)
        for cycles := 0; ; cycles++ {
                <-t.C
                for line := 0; line < mapSize.Y+2; line++ {
                        for column := 0; column < mapSize.X+3; column++ {
                                if column == mapSize.X+2 {
                                        mapDisplayBuffer[line*(mapSize.X+3)+column] = '\n'
                                        continue
                                }
                                if line == 0 || column == 0 || line == mapSize.Y+1 || column == mapSize.X+1 {
                                        mapDisplayBuffer[line*(mapSize.X+3)+column] = MapBoundChar
                                        continue
                                }
                                mapDisplayBuffer[line*(mapSize.X+3)+column] = MapFillChar
                        }
                }

                {
                        if snake.Head.Position.Y <= 0 || snake.Head.Position.X <= 0 ||
                                snake.Head.Position.X >= mapSize.X+1 || snake.Head.Position.Y >= mapSize.Y+1 {
                                if snake.Length == 1 {
                                        snake = NewSnake(3)
                                        loses++
                                } else {
                                        snake.RemoveTailNode()
                                }
                        }
                        node := snake.Head
                        for node != nil {
                                if node != snake.Head && node.Position == snake.Head.Position {
                                        if snake.Length == 1 {
                                                snake = NewSnake(3)
                                                loses++
                                        } else {
                                                snake.RemoveTailNode()
                                        }
                                }
                                if node.Position.Y > 0 && node.Position.X > 0 &&
                                        node.Position.X < mapSize.X+1 && node.Position.Y < mapSize.Y+1 {
                                        if node.Position == food.Position {
                                                snake.AddNodes(food.Saturation)
                                                totalScore += food.Saturation
                                                food = GenerateNewFood(mapSize, 2)
                                                if snake.Length > maxScore {
                                                        maxScore = snake.Length
                                                }
                                        }
                                        mapDisplayBuffer[node.Position.Y*(mapSize.X+3)+node.Position.X] = SnakeBodyChar
                                }
                                node = node.NextNode
                        }
                }

                mapDisplayBuffer[food.Position.Y*(mapSize.X+3)+food.Position.X] = FoodChar

                fmt.Printf("Cycles: %d\n", cycles)
                fmt.Printf("Score: %d; Loses: %d; Max score: %d; Total score: %d\n", snake.Length, loses, maxScore, totalScore)
                fmt.Printf("X: %d; Y: %d\n", snake.Head.Position.X, snake.Head.Position.Y)
                fmt.Println(string(mapDisplayBuffer))
                if food.Position.X > snake.Head.Position.X {
                        snake.Head.Direction = DirectionRight
                } else if food.Position.X < snake.Head.Position.X {
                        snake.Head.Direction = DirectionLeft
                } else if food.Position.Y > snake.Head.Position.Y {
                        snake.Head.Direction = DirectionDown
                } else if food.Position.Y < snake.Head.Position.Y {
                        snake.Head.Direction = DirectionUp
                }
                snake.Update()
        }
}

func randRange(min, max int) int {
        return rand.Intn(max+1-min) + min
}

type Food struct {
        Saturation int
        Position   Vector2
}

func GenerateNewFood(mapSize Vector2, saturation int) *Food {
        return &Food{
                Saturation: saturation,
                Position: Vector2{
                        X: randRange(1, mapSize.X),
                        Y: randRange(1, mapSize.Y),
                },
        }
}

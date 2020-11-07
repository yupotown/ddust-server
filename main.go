package main

import (
	"fmt"
	"log"
	"math/rand"
	"server/ddust"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	s, err := ddust.Init(0, [2]ddust.Color{ddust.Red, ddust.Black})
	if err != nil {
		log.Fatal(err)
	}

	drawDebug(s)

	fmt.Println()

	p := ddust.Input{
		Player: 0,
		Card:   s.Hands[0][0],
		X:      1,
		Y:      2,
	}

	s, err = ddust.Transition(s, p)
	if err != nil {
		log.Fatal(err)
	}

	drawDebug(s)
}

// drawDebug は、カードの裏表やプレイヤーに関わらずすべての状態を出力します。
func drawDebug(s ddust.State) {
	for y := 0; y < 4; y++ {
		if y == 0 {
			fmt.Print("盤面: ")
		} else {
			fmt.Print("      ")
		}
		for x := 0; x < 4; x++ {
			if x != 0 {
				fmt.Print(" ")
			}
			if s.Field[y][x].Front {
				fmt.Print(" ")
				drawCard(s.Field[y][x].Card)
				fmt.Print(" ")
			} else {
				fmt.Print("[")
				drawCard(s.Field[y][x].Card)
				fmt.Print("]")
			}
		}
		fmt.Println()
	}
	fmt.Print("山札 (左が上): ")
	for i := len(s.Deck) - 1; i >= 0; i-- {
		fmt.Print("[")
		drawCard(s.Deck[i])
		fmt.Print("]")
		if i != 0 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
	for i := 0; i < 2; i++ {
		fmt.Printf("プレイヤー%d (", i+1)
		if s.Colors[i] == ddust.Red {
			fmt.Print("赤")
		} else {
			fmt.Print("黒")
		}
		fmt.Print(")")
		if s.Turn == int64(i) {
			fmt.Print(" (手番)")
		}
		fmt.Print(":")
		for j := 0; j < len(s.Hands[i]); j++ {
			fmt.Print(" [")
			drawCard(s.Hands[i][j])
			fmt.Print("]")
		}
		fmt.Println()
	}
}

func drawCard(card ddust.Card) {
	suits := map[ddust.Suit]string{
		ddust.Heart:   "h",
		ddust.Spade:   "s",
		ddust.Diamond: "d",
		ddust.Clover:  "c",
	}
	nums := map[int64]string{
		1:  "A",
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "0",
		11: "J",
		12: "Q",
		13: "K",
	}

	if s, ok := suits[card.Suit]; ok {
		fmt.Print(s)
	} else {
		fmt.Print("?")
	}
	if n, ok := nums[card.Number]; ok {
		fmt.Print(n)
	} else {
		fmt.Print("?")
	}
}

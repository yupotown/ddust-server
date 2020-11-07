package ddust

// State は、ゲーム全体の状態です。
type State struct {
	Field    [4][4]FieldCard // 盤面
	Deck     []Card          // 山札（添字の大きいほうが上）
	Cemetery []Card          // 捨て札（添字の大きいほうが上）
	Hands    [2][]Card       // 手札
	Turn     int64           // 現在どちらのターンか
	Colors   [2]Color        // 各プレイヤーの色
}

// FieldCard は、盤面上のカードの状態です。
type FieldCard struct {
	Front bool // カードが表を向いているか
	Card  Card // カードの種類
}

// Card は、カードの種類です。
type Card struct {
	Suit   Suit  // マーク
	Number int64 // 数
}

// Suit は、カードのマークです。
type Suit string

const (
	// Heart は、ハートです。
	Heart Suit = "heart"

	// Spade は、スペードです。
	Spade Suit = "spade"

	// Diamond は、ダイヤです。
	Diamond Suit = "diamond"

	// Clover は、クローバーです。
	Clover Suit = "clover"
)

// Color は、カードの色です。
type Color string

const (
	// Red は、ハートとダイヤのカードの色です。
	Red Color = "red"

	// Black は、スペードとクローバーのカードの色です。
	Black Color = "black"
)

// Clone は、State のコピーを生成します。
func (s State) Clone() State {
	var res State

	// Field
	res.Field = [4][4]FieldCard{}
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			res.Field[y][x] = s.Field[y][x]
		}
	}

	// Deck
	res.Deck = make([]Card, len(s.Deck))
	for i := 0; i < len(s.Deck); i++ {
		res.Deck[i] = s.Deck[i]
	}

	// Cemetery
	res.Cemetery = make([]Card, len(s.Cemetery))
	for i := 0; i < len(s.Cemetery); i++ {
		res.Cemetery[i] = s.Cemetery[i]
	}

	// Hands
	res.Hands = [2][]Card{}
	for i := 0; i < 2; i++ {
		res.Hands[i] = make([]Card, len(s.Hands[i]))
		for j := 0; j < len(s.Hands[i]); j++ {
			res.Hands[i][j] = s.Hands[i][j]
		}
	}

	// Turn
	res.Turn = s.Turn

	// Colors
	res.Colors = [2]Color{s.Colors[0], s.Colors[1]}

	return res
}

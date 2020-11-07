package ddust

import (
	"fmt"
	"math/rand"
)

// Init は、ゲームの初期状態を生成します。
// 山札と盤面はランダムにシャッフルされ、スタートプレイヤーの手札は2枚になります。
// firstPlayer が 0 または 1 でない場合、error を返します。
// colors が {Black, Red} または {Red, Black} でない場合、error を返します。
func Init(firstPlayer int64, colors [2]Color) (State, error) {
	var res State

	// 入力値バリデーション
	if firstPlayer != 0 && firstPlayer != 1 {
		return res, fmt.Errorf("invalid first player number")
	}
	if colors[0] != Red && colors[0] != Black {
		return res, fmt.Errorf("invalid color")
	}
	if colors[1] != Red && colors[1] != Black {
		return res, fmt.Errorf("invalid color")
	}
	if colors[0] == colors[1] {
		return res, fmt.Errorf("invalid color")
	}

	suits := []Suit{Heart, Spade, Diamond, Clover}
	fieldNums := []int64{1, 11, 12, 13}

	// 盤面を生成してシャッフル
	res.Field = [4][4]FieldCard{}
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			res.Field[y][x] = FieldCard{
				Front: false,
				Card: Card{
					Suit:   suits[y],
					Number: fieldNums[x],
				},
			}
		}
	}
	for i := 0; i < 15; i++ {
		j := rand.Intn(16-i) + i
		iy, ix := i/4, i%4
		jy, jx := j/4, j%4
		res.Field[iy][ix], res.Field[jy][jx] = res.Field[jy][jx], res.Field[iy][ix]
	}

	// 山札（初期手札含む）を生成してシャッフル
	deckNums := []int64{2, 3, 4, 5}
	res.Deck = make([]Card, 16)
	for i := 0; i < 16; i++ {
		res.Deck[i] = Card{
			Suit:   suits[i/4],
			Number: deckNums[i%4],
		}
	}
	for i := 0; i < 15; i++ {
		j := rand.Intn(16-i) + i
		res.Deck[i], res.Deck[j] = res.Deck[j], res.Deck[i]
	}

	// 捨て札は空で初期化
	res.Cemetery = make([]Card, 0)

	// 手札を初期化して山札から引く
	res.Hands = [2][]Card{make([]Card, 0), make([]Card, 0)}
	res.Hands[0] = append(res.Hands[0], res.draw())
	res.Hands[1] = append(res.Hands[1], res.draw())
	res.Hands[firstPlayer] = append(res.Hands[firstPlayer], res.draw())

	// 手番
	res.Turn = firstPlayer

	// 各プレイヤーの色
	res.Colors = colors

	return res, nil
}

// Input は、ゲームに対する入力です。
// State と Input の組み合わせから次の State が決まります。
type Input struct {
	Player int64 // アクションを行ったプレイヤー
	Card   Card  // 使用するカード
	X      int64 // カードを使用する位置
	Y      int64 // カードを使用する位置
}

// Transition は、状態遷移を行います。
// 状態と入力の組み合わせが不正な場合、error を返します。
func Transition(s State, p Input) (State, error) {
	// 状態のバリデーション
	if err := s.Validate(); err != nil {
		return s, fmt.Errorf("invalid state: %w", err)
	}

	// 入力のバリデーション
	if err := p.Validate(); err != nil {
		return s, fmt.Errorf("invalid input: %w", err)
	}

	// アクションを行おうとしているプレイヤーのチェック
	if p.Player != s.Turn {
		return s, fmt.Errorf("it's not player %d's turn", p.Player)
	}

	// プレイヤーがそのカードを持っているかのチェック
	cardIdx := -1
	for i := 0; i < len(s.Hands[p.Player]); i++ {
		if s.Hands[p.Player][i] == p.Card {
			cardIdx = i
		}
	}
	if cardIdx < 0 {
		return s, fmt.Errorf("the player does not have specified card")
	}

	// 状態遷移
	next := s.Clone()

	// 範囲を反転する
	card := next.Hands[p.Player][cardIdx]
	shape := GetShape(card)
	for dy := 0; dy < 3; dy++ {
		for dx := 0; dx < 3; dx++ {
			if !shape[dy][dx] {
				continue
			}

			x, y := p.X+int64(dx-1), p.Y+int64(dy-1)
			if x < 0 || x >= 4 || y < 0 || y >= 4 {
				continue
			}

			next.Field[y][x].Front = !next.Field[y][x].Front
		}
	}

	// 手札から捨て札に移動する
	next.Cemetery = append(next.Cemetery, card)
	next.Hands[p.Player] = append(next.Hands[p.Player][:cardIdx], next.Hands[p.Player][cardIdx+1:]...)

	// 手番を変える
	next.Turn = (next.Turn + 1) % 2

	// 山札が残っている場合、手番プレイヤーがドローする
	if len(next.Deck) >= 1 {
		next.Hands[next.Turn] = append(next.Hands[next.Turn], next.draw())
	}

	return next, nil
}

// GetShape は、カードの表す形を取得します。
// true の場所が形に含まれ、[1][1] が中心です。
func GetShape(card Card) [3][3]bool {
	// カードの範囲チェック
	if card.Validate() != nil || card.Number <= 1 || card.Number >= 6 {
		return [3][3]bool{}
	}

	// カードと形の対応
	shapes := map[Suit][2][3]string{
		Heart: {
			// ハート 2, 4
			{
				"x.x",
				"xxx",
				"...",
			},
			// ハート 3, 5
			{
				"...",
				"xxx",
				"x.x",
			},
		},
		Spade: {
			// スペード 2, 4
			{
				".x.",
				".x.",
				".x.",
			},
			// スペード 3, 5
			{
				"...",
				"xxx",
				"...",
			},
		},
		Diamond: {
			// ダイヤ 2, 4
			{
				".x.",
				"xxx",
				".x.",
			},
			// ダイヤ 3, 5
			{
				"x.x",
				".x.",
				"x.x",
			},
		},
		Clover: {
			// クローバー 2, 4
			{
				".x.",
				"xxx",
				"...",
			},
			// クローバー 3, 5
			{
				"...",
				"xxx",
				".x.",
			},
		},
	}

	// 変換
	shape := shapes[card.Suit][card.Number%2]
	res := [3][3]bool{}
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			res[y][x] = (shape[y][x] == 'x')
		}
	}

	return res
}

// draw は、山札の一番上のカードを取得して山札から取り除きます。
// 山札が空の場合は panic が発生します。
func (s *State) draw() Card {
	res := s.Deck[len(s.Deck)-1]
	s.Deck = s.Deck[:len(s.Deck)-1]
	return res
}

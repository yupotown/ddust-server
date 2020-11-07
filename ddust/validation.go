package ddust

import "fmt"

// Validate は、State の値が正しいかバリデーションします。
func (s State) Validate() error {
	// カードの重複チェック用
	cards := map[Suit]map[int64]struct{}{
		Heart:   map[int64]struct{}{},
		Spade:   map[int64]struct{}{},
		Diamond: map[int64]struct{}{},
		Clover:  map[int64]struct{}{},
	}
	checkDup := func(c Card) error {
		if _, exists := cards[c.Suit][c.Number]; exists {
			return fmt.Errorf("duplicate cards (%s, %d)", c.Suit, c.Number)
		}
		cards[c.Suit][c.Number] = struct{}{}
		return nil
	}

	// 盤面
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			// Card の値
			if err := s.Field[y][x].Card.Validate(); err != nil {
				return fmt.Errorf("invalid field[%d, %d]'s card: %w", x, y, err)
			}

			// カードが A, J, Q, K のいずれかか
			if n := s.Field[y][x].Card.Number; n != 1 && n != 11 && n != 12 && n != 13 {
				return fmt.Errorf("invalid field[%d, %d]'s card", x, y)
			}

			// 重複
			if err := checkDup(s.Field[y][x].Card); err != nil {
				return err
			}
		}
	}

	// 山札
	for i := 0; i < len(s.Deck); i++ {
		// Card の値
		if err := s.Deck[i].Validate(); err != nil {
			return fmt.Errorf("invalid deck[%d]'s card: %w", i, err)
		}

		// カードが 2, 3, 4, 5 のいずれかか
		if n := s.Deck[i].Number; n != 2 && n != 3 && n != 4 && n != 5 {
			return fmt.Errorf("invalid deck[%d]'s card", i)
		}

		// 重複
		if err := checkDup(s.Deck[i]); err != nil {
			return err
		}
	}

	// 捨札
	for i := 0; i < len(s.Cemetery); i++ {
		// Card の値
		if err := s.Cemetery[i].Validate(); err != nil {
			return fmt.Errorf("invalid cemetery[%d]'s card: %w", i, err)
		}

		// カードが 2, 3, 4, 5 のいずれかか
		if n := s.Cemetery[i].Number; n != 2 && n != 3 && n != 4 && n != 5 {
			return fmt.Errorf("invalid cemetery[%d]'s card", i)
		}

		// 重複
		if err := checkDup(s.Cemetery[i]); err != nil {
			return err
		}
	}

	// 手札
	for i := 0; i < 2; i++ {
		// 枚数
		if s.Turn == int64(i) && len(s.Hands[i]) != 2 {
			return fmt.Errorf("invalid number of cards in hands[%d]", i)
		}
		if s.Turn != int64(i) && len(s.Hands[i]) != 1 {
			return fmt.Errorf("invalid number of cards in hands[%d]", i)
		}

		for j := 0; j < len(s.Hands[i]); j++ {
			// Card の値
			if err := s.Hands[i][j].Validate(); err != nil {
				return fmt.Errorf("invalid hands[%d][%d]'s card: %w", i, j, err)
			}

			// カードが 2, 3, 4, 5 のいずれかか
			if n := s.Hands[i][j].Number; n != 2 && n != 3 && n != 4 && n != 5 {
				return fmt.Errorf("invalid hands[%d][%d]'s card", i, j)
			}

			// 重複
			if err := checkDup(s.Hands[i][j]); err != nil {
				return err
			}
		}
	}

	// 山札と捨て札と手札の合計枚数
	if len(s.Deck)+len(s.Cemetery)+len(s.Hands[0])+len(s.Hands[1]) != 16 {
		return fmt.Errorf("invalid number of cards")
	}

	// ターン
	if s.Turn != 0 && s.Turn != 1 {
		return fmt.Errorf("invalid turn")
	}

	// 各プレイヤーの色
	if s.Colors[0] == s.Colors[1] {
		return fmt.Errorf("invalid colors")
	}
	for i := 0; i < 2; i++ {
		if s.Colors[i] != Red && s.Colors[i] != Black {
			return fmt.Errorf("invalid colors[%d]", i)
		}
	}

	return nil
}

// Validate は、Card の値が正しいかバリデーションします。
func (c Card) Validate() error {
	// マーク
	suits := map[Suit]struct{}{
		Heart:   struct{}{},
		Spade:   struct{}{},
		Diamond: struct{}{},
		Clover:  struct{}{},
	}
	if _, ok := suits[c.Suit]; !ok {
		return fmt.Errorf("invalid suite")
	}

	// 数
	if c.Number <= 0 || c.Number >= 14 || (c.Number >= 6 && c.Number <= 10) {
		return fmt.Errorf("invalid number")
	}

	return nil
}

// Validate は、Input の値が正しいかバリデーションします。
func (p Input) Validate() error {
	// プレイヤー
	if p.Player != 0 && p.Player != 1 {
		return fmt.Errorf("invalid player")
	}

	// カード
	if err := p.Card.Validate(); err != nil {
		return fmt.Errorf("invalid card: %w", err)
	}
	if n := p.Card.Number; n != 2 && n != 3 && n != 4 && n != 5 {
		return fmt.Errorf("invalid card")
	}

	// 位置
	if p.X < 0 || p.X >= 4 || p.Y < 0 || p.Y >= 4 {
		return fmt.Errorf("invalid coordinate")
	}

	return nil
}

package bowling

import (
	"testing"
)

func TestScore(t *testing.T) {
	testCases := []struct {
		name        string
		game        string
		expected    int
		expectError bool
	}{
		{
			name:        "Идеальная игра (все страйки)",
			game:        "XXXXXXXXXXXX",
			expected:    300,
			expectError: false,
		},
		{
			name:        "Все спэры с пятерками",
			game:        "5/5/5/5/5/5/5/5/5/5/5",
			expected:    150,
			expectError: false,
		},
		{
			name:        "Простая игра без бонусов",
			game:        "9-9-9-9-9-9-9-9-9-9-",
			expected:    90,
			expectError: false,
		},
		{
			name:        "Все промахи",
			game:        "--------------------",
			expected:    0,
			expectError: false,
		},
		{
			name:        "Спэр в последнем фрейме",
			game:        "1111111111111111111/5",
			expected:    33,
			expectError: false,
		},
		{
			name:        "Страйк в последнем фрейме",
			game:        "111111111111111111XX5",
			expected:    43,
			expectError: false,
		},
		{
			name:        "Игра со смешанными результатами",
			game:        "X7/9-X-88/-63/X81",
			expected:    138,
			expectError: false,
		},
		{
			name:        "Некорректная игра (недостаточно бросков)",
			game:        "X",
			expected:    0,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			score, err := Score(tc.game)

			if tc.expectError {
				if err == nil {
					t.Errorf("ожидалась ошибка, но получено nil")
				}
			} else {
				if err != nil {
					t.Errorf("неожиданная ошибка: %v", err)
				}
				if score != tc.expected {
					t.Errorf("ожидаемый счет %d, но получено %d", tc.expected, score)
				}
			}
		})
	}
}

package verse

import (
	"errors"
	"strings"
)

var (
	ErrNoVerse = errors.New("no verse in song text")
)

// получение текста по пагинации
func TextPaginate(songText string, versesPage int) (string, error) {
	verses := strings.Split(songText, "\n\n") // Предполагаем, что куплеты разделены двумя переносами строк
	if len(verses) < 1 {
		return "", ErrNoVerse
	}

	for i := 0; i < len(verses); i++ {
		if versesPage == i+1 && len(verses[i]) > 0 {
			return verses[i], nil
		}
	}
	return "", ErrNoVerse
}

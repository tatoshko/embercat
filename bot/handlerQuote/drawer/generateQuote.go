package drawer

import "strings"

func generateQuote(text string) {
    words := strings.Split(text, " ")

    currentRow := 0
    rows := [][]string{}
    for i, word := range words {
        if i%3 == 0 {
            currentRow += 1
        }
        rows[currentRow] = append(rows[currentRow], word)
    }
}

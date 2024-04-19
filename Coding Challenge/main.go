package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func validateCard(card string) bool {

	pattern := `^(4|5|6)\d{3}(-?\d{4}){3}$`

	re := regexp.MustCompile(pattern)

	if !re.MatchString(card) {
		return false
	}

	card = regexp.MustCompile(`-`).ReplaceAllString(card, "")

	for i := 0; i < len(card)-3; i++ {
		if card[i] == card[i+1] && card[i] == card[i+2] && card[i] == card[i+3] {
			return false
		}
	}

	return true
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	var n int
	fmt.Sscanf(scanner.Text(), "%d", &n)

	for i := 0; i < n; i++ {
		scanner.Scan()
		card := scanner.Text()

		if validateCard(card) {
			fmt.Println("Valid")
		} else {
			fmt.Println("Invalid")
		}
	}
}

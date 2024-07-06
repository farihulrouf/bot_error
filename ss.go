package main

import (
	"fmt"
)

func main() {
	// Mencetak "Hello, World!"
	fmt.Println("Hello, World!")

	// Membuat slice berisi angka-angka
	numbers := []int{1, 2, 3, 4, 5}

	// Menghitung jumlah angka dalam slice
	sum := 0
	for _, number := range numbers {
		sum += number
	}

	// Mencetak jumlah angka
	fmt.Printf("Jumlah angka dalam slice: %d\n", sum)
}

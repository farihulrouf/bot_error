package main

import (
	"fmt"
	"os"
)

func main() {
	// Ganti dengan path lengkap ke file yang ingin Anda buka
	filePath := "/home/farihul/Pictures/download.png"

	// Buka file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// File berhasil dibuka, lakukan operasi yang diinginkan di sini
	fmt.Println("File successfully opened")

	// Contoh membaca isi file
	// Anda dapat menggunakan buffer atau scanner untuk membaca isi file
	// Berikut adalah contoh menggunakan buffer untuk membaca file
	/*
		buffer := make([]byte, 1024)
		n, err := file.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			return
		}
		fmt.Printf("Read %d bytes: %s\n", n, buffer[:n])
	*/
}

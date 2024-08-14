package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	str := os.Args[1]
	banner := os.Args[2]
	content, err := os.ReadFile(banner+".txt")
	if err != nil {
		log.Fatal(err)
	}
	res1 := strings.Split(str, "\\n")
	temp := strings.Split(string(content), "\n") // splitting the txt file
	count := 0                                   // counter for height
	for a := 0; a < len(res1); a++ {             // going through split stings if it had \n
		if os.Args[1] == "" {
			break
		}
		if os.Args[1] == "\\n" {
			fmt.Printf("\n")
			break
		}
		str1 := []byte(res1[a])
		lenn := len(str1) - 1
		if res1[a] == "" {
			fmt.Printf("\n")
		}
		for i := 0; i <= lenn; { // going through by line
			if count == 8 {
				i = 0
				count = 0
				break
			}
			if i < lenn {
				fmt.Printf(temp[((rune(str1[i])-32)*9 + 1 + rune(count))])
				i++
			}
			if i == lenn {
				fmt.Printf(temp[((rune(str1[i])-32)*9 + 1 + rune(count))])
				fmt.Printf("\n")
				count++
				i = 0
			}
		}
	}
}

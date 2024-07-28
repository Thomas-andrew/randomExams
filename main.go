package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var err error

	fmt.Println("--- Random Ex ---")
	fmt.Println("Enter the textbook name:")

	textbookName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading textbook name.")
		return
	}
	textbookName = textbookName[:len(textbookName)-1]

	fmt.Printf("textbook: '%v'\n", textbookName)

	err = screenshoot("./" + textbookName + ".png")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
}

func screenshoot(path string) error {
	args := []string{"-s", "-m", "10", path}
	cmd := exec.Command("maim", args...)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Command finished with error: %v\n", err)
		return fmt.Errorf("[screenshoot] maim finished with error: %w", err)
	}

	return nil
}

// Author: nigelchoi@google.com (Nigel Choi)

// Prints all piece definitions.

package main

import (
	"fmt"
	"9gel/gknot"
)

func main() {
	gknot.BluePieceDef.Print()
	fmt.Println()
	gknot.OrangePieceDef.Print()
	fmt.Println()
	gknot.PurplePieceDef.Print()
	fmt.Println()
	gknot.GreenPieceDef.Print()
	fmt.Println()
	gknot.RedPieceDef.Print()
	fmt.Println()
	gknot.YellowPieceDef.Print()
}

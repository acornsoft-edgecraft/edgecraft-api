/*
Copyright 2023 product Authors. All right reserved.
*/
package main

import (
	"fmt"
	"time"
)

// ===== [ Constants and Variables ] =====

// main - entry point
func main() {
	fmt.Println(time.Now())
	fmt.Println(time.Now().UTC())
}

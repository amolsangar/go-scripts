package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"example.com/apiConn"
)

func main() {
	defer fmt.Println("exiting application")

	apiConnection := apiConn.Open()

	var wg sync.WaitGroup
	wg.Add(20)

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			v, err := apiConnection.Read(context.Background())
			if err != nil {
				fmt.Printf("%v Get Error: %v\n", time.Now().Format("15:04:05"), err)
				return
			}

			fmt.Printf("%v %v\n", time.Now().Format("15:04:05"), v)
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			err := apiConnection.Resolve(context.Background())
			if err != nil {
				fmt.Printf("%v Resolve Error: %v\n", time.Now().Format("15:04:05"), err)
				return
			}

			fmt.Printf("%v Resolved\n", time.Now().Format("15:04:05"))
		}()
	}

	wg.Wait()
}

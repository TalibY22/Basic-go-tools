package main

import (
	"fmt"
	"net"
	"sync"
)


func scan(wg *sync.WaitGroup)  {
	
	defer wg.Done() // basically shows that the goroutine has finished 


	for i:=1; i< 1024; i++{
		address := fmt.Sprintf("scanme.nmap.org:%d", i)
		conn,err := net.Dial("tcp", address)
		if err != nil{
			continue
		}
		conn.Close()
		fmt.Printf("%d open\n", i)

	}
	fmt.Println("Scan completed")
	
}


func main ()  {
	
	var wg sync.WaitGroup //intialize the wait group
    
	wg.Add(1)//Set that theres one gouroutine thats gonna run 


	
	scan(&wg)
	
	wg.Wait()

	
	fmt.Println("Scan completed")

}

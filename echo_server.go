package main

import (
	
	"io"
	"log"
	"net"
	
) 


func echo(conn net.Conn)  {
	
	b := make([]byte,512)
	size,err:= conn.Read(b[0:])

	if err != nil{
        log.Println("An erro occured")

	}

	if err == io.EOF{
		log.Println("client has disconnected ")
	}

    
	log.Println("Writing data")
	log.Printf("Received %d bytes: %s\n", size, string(b))

	

	if  _, err := conn.Write(b[0:size]); err != nil{
		log.Fatalln("Unable to write data")

	}

}




func main()  {
	

 listener,err := net.Listen("tcp",":20080")
 if err != nil{
	log.Println("unable to bind port port may be in use ")
 }
 log.Println("listening on 20080")
 for {
   conn,err := listener.Accept()
   log.Println("connection accepted")
   if err != nil{
	log.Println("unable to connect")
   } 

   go echo(conn)

}

}
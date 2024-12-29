package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)



func execute(command string) error {



	command = strings.TrimSuffix(command,"\n")

	arg := strings.Split(command," ")

    var cmd *exec.Cmd
	
	if len(arg) == 1 {
		cmd = exec.Command(command)
		
		
		
	} else {
	
	   cmd = exec.Command(arg[0],arg[1])


	
	}

	
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()

}



func main(){


	read := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

        //take input 	
		input,err := read.ReadString('\n')

		if err != nil{
			fmt.Println("something aint write man ")
		}


		if err =execute(input); err != nil{

			fmt.Println("am a screwup ")
		}

	}


}
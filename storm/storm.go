package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func scoringCheck() (int, int) {
	r := "\033[31m" // Red
	g := "\033[32m" // Green
	w := "\033[97m" // White
	b := "\033[34m" // Blue

	fmt.Println("Checking for points:")
	points := 0
	possiblePoints := 0

	//Check if the webserver is up
	fmt.Print("\t1) ")
	possiblePoints++
	_, err := exec.Command("pgrep", "apache2").Output()
	if err != nil {
		// The command will error if the process is not running
		fmt.Println(r + "The webserver is down" + w)
	} else {
		points += 1
		fmt.Println(g + "The apache2 service is running" + w)
	}

	// Firefox version
	fmt.Print("\t2) ")
	possiblePoints++
	out, err := exec.Command("firefox", "-v").Output()
	if err != nil {
		log.Fatal(err)
	}
	firefoxV := string(out[16:20])
	if firefoxV > "59.0.2" {
		points += 1
		fmt.Println(g+"Firefox has been updated to:"+w, firefoxV)
	} else {
		fmt.Println(r + "Misconfiguration/Vulnerability present" + w)
	}

	// Kernel version
	fmt.Print("\t3) ")
	possiblePoints++
	out, err = exec.Command("uname", "-srm").Output()
	if err != nil {
		log.Fatal(err)
	}
	kernelV := string(out[6:15])
	if kernelV != "4.15.0-20" { //TODO fix before future deployment. (The current implementation is limited to Ubuntu 18)
		points += 1
		fmt.Println(g+"The Kernel has been updated to:"+w, kernelV)
	} else {
		fmt.Println(r + "Misconfiguration/Vulnerability present" + w)
	}

	//Check for unauthorized user: hax0r
	fmt.Print("\t4) ")
	possiblePoints++
	out, err = exec.Command("getent", "passwd", "hax0r").Output()
	user := string(out)
	if err != nil {
		// The command will error if no user was found
		user = ""
	}
	if string(user) == "" {
		points += 1
		fmt.Println(g + "The hax0r user was removed" + w)
	} else {
		fmt.Println(r + "Misconfiguration/Vulnerability present" + w)
	}

	//Check if the database is up
	fmt.Print("\t5) ")
	possiblePoints++
	_, err = exec.Command("pgrep", "mysql").Output()
	if err != nil {
		// The command will error if the process is not running
		fmt.Println(r + "The database is down" + w)
	} else {
		points += 1
		fmt.Println(g + "The mysql service is running" + w)
	}
	// ENABLE THE FIREWALL!!! MWAaahrhahahhaaa

	fmt.Println("Total points earned", b, points, w, "out of", b, possiblePoints, w) // TODO Show a out of b problems found
	return points, possiblePoints
}

func localScore(username string) {
	fmt.Println("\033[34m" + "Starting Storm..." + "\033[97m")
	fmt.Println("Welcome", username)
	for {
		time.Sleep(6 * time.Second)
		fmt.Print("\033[H\033[2J")
		score, possible := scoringCheck()
		time.Sleep(1 * time.Second)
		client(strconv.Itoa(score) + strconv.Itoa(possible) + username)
		fmt.Println("Clearing terminal in 6 seconds...")
	}
}

func help() {
	fmt.Println("Usage:\n\tStart your scoring system:\n\t\t$ storm name \"YourUsernameHere\"")
	fmt.Println("\tCheck your current services:\n\t\t$ storm score")
	fmt.Println("\tRun the group scoring server:\n\t\t$ storm scoreserver")
}

func server() {
	fmt.Println("Server starting...")

	PORT := ":" + "1337"
	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()
	buffer := make([]byte, 1024)

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		message := string(buffer[0 : n-1])

		g := "\033[32m"
		w := "\033[97m"

		correct := string(message[0])
		possible := string(message[1])
		user := message[2:]

		fmt.Println("-> ", user)

		//fmt.Println("Score:",g,correct,w,"/", g,possible,w, "Username:", g,user,w, "From:", g,addr, w)// Green vars
		fmt.Print(g, "Score: ", w, correct, g, " out of ", w, possible, g, " From: ", w, addr, "\n\n") //White vars

		data := []byte(message) // Send back to confirm
		_, err = connection.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func client(message string) {
	// This is so jankey lol
	// This takes a string where position 0 is the score.
	// Position 1 is the possible amount of points
	// Position 2 and onward is the username

	fmt.Println("Updating scoring server...")
	CONNECT := "muse.jaxlo.net:1337"

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer c.Close()

	data := []byte(message + "\n")
	_, err = c.Write(data)

	if err != nil {
		fmt.Println(err)
		return
	}

	buffer := make([]byte, 1024)
	n, _, err := c.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	check := string(buffer[0:n])
	if check == message {
		fmt.Println("Update successful")
	} else {
		fmt.Println("Response does not match message sent...")
	}
}

func main() {
	if len(os.Args) == 1 {
		help()
		os.Exit(1)
	}
	if os.Args[1] == "score" {
		fmt.Println("Note: run with the name command to connect to a scoring server")
		scoringCheck()
	} else if os.Args[1] == "name" {
		username := os.Args[2]
		localScore(username)
	} else if os.Args[1] == "scoreserver" {
		server()
	} else {
		fmt.Println("Invalid input")
		help()
	}
}

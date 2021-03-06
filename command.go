package main

import (
	"log"
	"os"
)

func runCommands(host, user, pwd string, commands ...string) {
	session, err := connect(user, pwd, host, 22)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	for _, cmd := range commands {
		log.Println("$ " + cmd)

		if err = session.Run(cmd); err != nil {
			log.Println(err)
		}
	}

}

func getCommandsOutput(host, user, pwd, cmd string) string {

	session, err := connect(user, pwd, host, 22)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	session.Stderr = os.Stderr

	log.Println("$ " + cmd)

	result, err := session.Output(cmd)

	if err != nil {
		log.Println(err)
		return ""
	}

	return string(result)
}

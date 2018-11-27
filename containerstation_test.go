package containerstation_test

import (
	"log"

	"context"
	"net/http"

	cstation "github.com/iaburton/containerstation"
)

func ExampleNewClient() {
	//nil or any custom client can be passed
	cc := cstation.NewClient("https://Nasstorage", http.DefaultClient)

	bg := context.Background()
	log.Print("Login")
	lr, err := cc.Login(bg, user, pass)
	if err != nil {
		log.Fatalf("%T %v", err, err)
	}
	log.Printf("%+v", lr)

	defer func() {
		log.Print("Logout")
		outr, err := cc.Logout(bg)
		if err != nil {
			log.Fatalf("%T %v", err, err)
		}
		log.Printf("%+v", outr)
	}()

	log.Print("SystemInformation")
	sr, err := cc.SystemInformation(bg)
	if err != nil {
		log.Panicf("%T %v", err, err)
	}
	log.Printf("%+v", sr)

	log.Print("ResourceUsage")
	ru, err := cc.ResourceUsage(bg)
	if err != nil {
		log.Panicf("%T %v", err, err)
	}
	log.Printf("%+v", ru)

	log.Print("NetworkPort")
	used, err := cc.NetworkPort(bg, cstation.TCP, 443)
	if err != nil {
		log.Panicf("%T %v", err, err)
	}
	log.Printf("Port 443 used: %v", used)

	log.Print("ListContainers")
	list, err := cc.ListContainers(bg)
	if err != nil {
		log.Panicf("%T %v", err, err)
	}
	for _, c := range list {
		log.Printf("Container: %+v", c)
	}
	// Output:
}

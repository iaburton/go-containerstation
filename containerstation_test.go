package containerstation_test

import (
	"log"

	"context"
	"crypto/tls"
	"net/http"

	cstation "github.com/iaburton/containerstation"
)

func ExampleNewClient() {
	cc := cstation.NewClient("https://Nasstorage", &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})

	bg := context.Background()
	log.Print("Login")
	lr, err := cc.Login(bg, "stuff", "things")
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
	// Output:
}

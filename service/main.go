package main

func main() {
	log.Println("starting...")
	signer, err := GetHostSigner()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(signer)
}

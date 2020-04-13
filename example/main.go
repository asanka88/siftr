package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"siftr"
)

func main() {
	// read the data file
	b, err := ioutil.ReadFile("example/data.json")
	if err != nil {
		log.Fatal(err)
	}

	// read the policy file
	pol, err := ioutil.ReadFile("example/policy.json")
	if err != nil {
		log.Fatal(err)
	}
	// read the policy data as the siftr.Policy
	var policy siftr.Policy
	err = json.Unmarshal(pol, &policy)
	if err != nil {
		log.Fatal(err)
	}

	// Sift
	data, err := siftr.Sift(b, &policy)
	if err != nil {
		log.Fatalf("error sifting: %s", err.Error())
	}

	// just beautify the data and then print it
	j, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(j))
}

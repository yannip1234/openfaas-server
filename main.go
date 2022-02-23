package main

import (
	"encoding/json"
	b64 "encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
	"os"
	"strconv"
)

type Payload struct {
	Fid    string `json:"fid"`
	Src    string `json:"src"`
	Params string `json:"params,omitempty"`
	Lang   string `json:"lang"`
}
type Response struct {
	Fid         string        `json:"fid"`
	TimeElapsed time.Duration `json:"time_elapsed"`
	Result      string        `json:"result"`
	Error       string        `json:"error,omitempty"`
}

var Payloads []Payload

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/showfunctions", returnAllFunctions)
	myRouter.HandleFunc("/showfunctions/{fid}", returnFunction)
	myRouter.HandleFunc("/run", runFunction).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
func returnFunction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnFunction")
	vars := mux.Vars(r)
	key := vars["fid"]
	for _, payload := range Payloads {
		if payload.Fid == key {
			json.NewEncoder(w).Encode(payload)
		}
	}
}
func runFunction(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var payload Payload
	res := Response{}
	json.Unmarshal(reqBody, &payload)
	params:= payload.Params
	// src:= payload.Src
	data, _ := b64.StdEncoding.DecodeString(payload.Src)
	s2, _ := strconv.Unquote(string(data))
	fmt.Println(s2)
	err := os.WriteFile("/root/src.py",[]byte(s2), 0644)
	if err != nil {
		panic(err)
	}
	
	if strings.ToLower(payload.Lang) == "micropython"{
		//fmt.Println(string(out))	
		start := time.Now()
		cmd := exec.Command( "micropython",  "/root/src.py", params)
		out, err := cmd.CombinedOutput()
		t := time.Now()
		//fmt.Println(string(out))

		if err != nil {
			res = Response{Fid: payload.Fid, TimeElapsed: t.Sub(start), Result: string(out), Error: err.Error()}
			fmt.Println(err)
		} else {
			res = Response{Fid: payload.Fid, TimeElapsed: t.Sub(start), Result: string(out)}
		}

	} else if strings.ToLower(payload.Lang) == "python" {
		start := time.Now()
		cmd := exec.Command("python", "/root/src.py",  params)
		out, err := cmd.CombinedOutput()
		t := time.Now()

		if err != nil {
			res = Response{Fid: payload.Fid, TimeElapsed: t.Sub(start), Result: string(out), Error: err.Error()}
			fmt.Println(err)
		} else {
			res = Response{Fid: payload.Fid, TimeElapsed: t.Sub(start), Result: string(out)}
		}
	} else {
		res = Response{Fid: payload.Fid, Result: "Language not supported"}
	}

	json.NewEncoder(w).Encode(res)

	//fmt.Fprintf(w, "%+v", string(reqBody))
}
func returnAllFunctions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllFunctions")
	json.NewEncoder(w).Encode(Payloads)
}
func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Payloads = []Payload{
		{Fid: "helloworld", Src: "print(\"hello\")", Lang: "Python", Params: "{\"A\":\"3\", \"B\":\"5\"}"},
		{Fid: "helloworld2", Src: "print(\"hello\")", Lang: "Python", Params: "{\"A\":\"4\", \"B\":\"6\"}"},
	}
	handleRequests()
}

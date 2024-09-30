package main

import "gitlab.com/mergetb/tech/nex/svc/agent/options"





func main() {
	agent := options.New()
	agent.Run()
}
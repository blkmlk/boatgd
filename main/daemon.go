package main

import (
	"../can"
	"../sender"
)

func RunDaemon() {
	can.RunCANReader()

	sender.RunSender()
}
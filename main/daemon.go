package main

import (
	"bg/can"
	"bg/sender"
)

func RunDaemon() {
	can.RunCANReader()

	sender.RunSender()
}
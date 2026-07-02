package main

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-gl/mathgl/mgl64"
)

const loadedMessage = "[系統] 讀取人物成功。"

var (
	emeraldBalance atomic.Int64

	tabBalancePattern     = regexp.MustCompile(`綠寶石餘額 : ([\d,]+)`)
	privateMessagePattern = regexp.MustCompile(`^\[(\w+) -> 您]\s(.*)`)
	teleportPattern       = regexp.MustCompile(`^\[系統] ([\w]+) 想要傳送到 你 的位置`)
	teleportToPattern     = regexp.MustCompile(`^\[系統] ([\w]+) 想要你傳送到 該玩家 的位置`)
)

func parseEmeraldBalance(message string) (int64, bool) {
	matches := tabBalancePattern.FindStringSubmatch(message)
	if len(matches) != 2 {
		return 0, false
	}
	balance, err := strconv.ParseInt(strings.ReplaceAll(matches[1], ",", ""), 10, 64)
	if err != nil {
		return 0, false
	}
	return balance, true
}

func handlePrivateMessage(message string) {
	matches := privateMessagePattern.FindStringSubmatch(message)
	if len(matches) != 3 {
		return
	}
	playerName := matches[1]
	if !isOwner(playerName) {
		return
	}

	go runOwnerCommand(playerName, matches[2])
}

func handleTeleportRequest(message string) {
	playerName, ok := teleportPlayer(message)
	if !ok {
		return
	}
	if isOwner(playerName) {
		runCommand("tpaccept " + playerName)
		return
	}
	runCommand("tno")
}

func teleportPlayer(message string) (string, bool) {
	if matches := teleportPattern.FindStringSubmatch(message); len(matches) == 2 {
		return matches[1], true
	}
	if matches := teleportToPattern.FindStringSubmatch(message); len(matches) == 2 {
		return matches[1], true
	}
	return "", false
}

func runOwnerCommand(playerName, message string) {
	args := strings.Fields(message)
	if len(args) == 0 {
		return
	}

	switch args[0] {
	case "dropAll":
		dropAll()
	case "get":
		runCommand(fmt.Sprintf("pay %s %d", playerName, emeraldBalance.Load()))
	case "sethome":
		pos := c.Player().Entity().Position()
		if err := c.Player().FlyTo(mgl64.Vec3{pos.X(), pos.Y() + 4, pos.Z()}); err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(50 * time.Millisecond)
		runCommand("sethome " + playerName)
	case "cmd":
		if len(args) > 1 {
			runCommand(strings.Join(args[1:], " "))
		}
	case "getdata":
		fmt.Println(glassRID)
	case "glass":
		go startCraftLoop()
	}
}

func dropAll() {
	inventory := c.Inventory().Inventory()
	for i := 9; i <= 45 && i < inventory.SlotCount(); i++ {
		if inventory.GetSlot(i).ItemID != 0 {
			_ = inventory.Click(int16(i), 4, 1)
		}
	}
}

func runCommand(command string) {
	if err := c.Player().Command(command); err != nil {
		fmt.Println(err)
	}
}

func isOwner(playerName string) bool {
	return slices.Contains(cfg.Owners, strings.ToLower(playerName))
}

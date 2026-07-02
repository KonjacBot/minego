package main

import "testing"

func Test_parseEmeraldBalance_whenTabListContainsBalance(t *testing.T) {
	balance, ok := parseEmeraldBalance("綠寶石餘額 : 12,345 / 村民錠餘額 : 0")

	if !ok {
		t.Fatal("expected balance match")
	}
	if balance != 12345 {
		t.Fatalf("balance = %d, want 12345", balance)
	}
}

func Test_teleportPlayer_whenTeleportToOwnerMessage(t *testing.T) {
	playerName, ok := teleportPlayer("[系統] PatyHank 想要你傳送到 該玩家 的位置")

	if !ok {
		t.Fatal("expected teleport match")
	}
	if playerName != "PatyHank" {
		t.Fatalf("playerName = %q, want PatyHank", playerName)
	}
}

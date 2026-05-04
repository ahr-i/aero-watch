package gps

import (
	"testing"
	"time"
)

func TestStoreUpsertGetAndList(t *testing.T) {
	store := NewStore(10*time.Second, time.Hour)
	defer store.Close()

	store.Upsert("group-a", "drone-1", 37.5665, 126.9780)

	position, ok := store.Get("group-a", "drone-1")
	if !ok {
		t.Fatal("expected gps position to exist")
	}

	if position.Latitude != 37.5665 || position.Longitude != 126.9780 {
		t.Fatalf("unexpected position: %+v", position)
	}

	positions := store.List()
	if len(positions) != 1 {
		t.Fatalf("expected 1 position, got %d", len(positions))
	}
}

func TestStoreExpiresPosition(t *testing.T) {
	store := NewStore(time.Millisecond, time.Hour)
	defer store.Close()

	store.Upsert("group-a", "drone-1", 37.5665, 126.9780)
	time.Sleep(2 * time.Millisecond)

	_, ok := store.Get("group-a", "drone-1")
	if ok {
		t.Fatal("expected gps position to expire")
	}
}

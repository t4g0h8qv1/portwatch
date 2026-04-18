package portinventory_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/portinventory"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "inventory.json")
}

func TestSet_And_Get(t *testing.T) {
	inv := portinventory.New()
	if err := inv.Set("host1", []int{80, 443}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	ports, ok := inv.Get("host1")
	if !ok {
		t.Fatal("expected host1 to exist")
	}
	if len(ports) != 2 || ports[0] != 80 || ports[1] != 443 {
		t.Fatalf("unexpected ports: %v", ports)
	}
}

func TestSet_EmptyHost(t *testing.T) {
	inv := portinventory.New()
	if err := inv.Set("", []int{80}); err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestGet_Missing(t *testing.T) {
	inv := portinventory.New()
	_, ok := inv.Get("ghost")
	if ok {
		t.Fatal("expected missing host")
	}
}

func TestHosts_Sorted(t *testing.T) {
	inv := portinventory.New()
	_ = inv.Set("zebra", []int{22})
	_ = inv.Set("alpha", []int{80})
	hosts := inv.Hosts()
	if hosts[0] != "alpha" || hosts[1] != "zebra" {
		t.Fatalf("unexpected order: %v", hosts)
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempPath(t)
	inv := portinventory.New()
	_ = inv.Set("srv1", []int{22, 80, 443})
	_ = inv.Set("srv2", []int{8080})
	if err := inv.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := portinventory.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	ports, ok := loaded.Get("srv1")
	if !ok || len(ports) != 3 {
		t.Fatalf("unexpected ports for srv1: %v", ports)
	}
	ports2, ok2 := loaded.Get("srv2")
	if !ok2 || len(ports2) != 1 || ports2[0] != 8080 {
		t.Fatalf("unexpected ports for srv2: %v", ports2)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	inv, err := portinventory.Load("/nonexistent/path/inventory.json")
	if err != nil {
		t.Fatalf("expected empty inventory, got error: %v", err)
	}
	if len(inv.Hosts()) != 0 {
		t.Fatal("expected empty inventory")
	}
}

func TestSet_PortsSorted(t *testing.T) {
	inv := portinventory.New()
	_ = inv.Set("h", []int{443, 22, 80})
	ports, _ := inv.Get("h")
	if ports[0] != 22 || ports[1] != 80 || ports[2] != 443 {
		t.Fatalf("ports not sorted: %v", ports)
	}
	_ = os.Remove("") // suppress unused import
}

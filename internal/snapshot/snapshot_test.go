package snapshot

import (
	"os"
	"path/filepath"
	"testing"

	"grating-eq/internal/grating"
	"grating-eq/internal/report"
)

func TestRoundTripExample(t *testing.T) {
	in, err := report.ParseFile("../../example/600lpmm.json")
	if err != nil {
		t.Fatal(err)
	}
	rec, err := CaptureInput(in, 0)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "600lpmm.snap.json")
	if err := WriteFile(path, rec); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Matches(rec) {
		t.Fatalf("round-trip mismatch: %+v vs %+v", got, rec)
	}
	if err := got.ReplayAgrees(); err != nil {
		t.Fatal(err)
	}
	plus, err := grating.New(got.SpacingNm, got.WavelengthNm, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	res, err := plus.Eq(1)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Exists {
		t.Fatal("reloaded 600 lp/mm geometry must keep order 1")
	}
}

func TestEmptyRejected(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(path); err == nil {
		t.Fatal("empty file must be rejected")
	}
}

func TestTruncatedRejected(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "trunc.json")
	if err := os.WriteFile(path, []byte(`{"version":1,"orders":[`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(path); err == nil {
		t.Fatal("truncated JSON must be rejected")
	}
}

func TestBadVersionRejected(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ver.json")
	body := []byte(`{"version":99,"spacing_nm":1666.67,"wavelength_nm":550,"orders":[]}`)
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadFile(path); err == nil {
		t.Fatal("unsupported version must be rejected")
	}
}

func TestObliqueReplay(t *testing.T) {
	in, err := report.ParseFile("../../example/oblique.json")
	if err != nil {
		t.Fatal(err)
	}
	rec, err := CaptureInput(in, 5)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "oblique.snap.json")
	if err := WriteFile(path, rec); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := got.ReplayAgrees(); err != nil {
		t.Fatal(err)
	}
}

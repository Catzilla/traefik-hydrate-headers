package traefik_hydrate_headers

import "testing"

func TestCompactJson(t *testing.T) {
	json := []byte(`{
		"x": 1,
		"y": 2
	}`)
	want := []byte(`{"x":1,"y":2}`)
	compact, _ := compactJson(json)
	if string(compact) != string(want) {
		t.Errorf("got %q, wanted %q", compact, want)
	}

	invalidJson := []byte(`invalid`)
	_, err := compactJson(invalidJson)
	if err == nil {
		t.Error("no error for invalid json")
	}
}

func TestContains(t *testing.T) {
	src := []int{1, 2, 3}

	if !contains(src, 1) {
		t.Error("returns false with present value")
	}

	if contains(src, 4) {
		t.Error("returns true with missing value")
	}
}

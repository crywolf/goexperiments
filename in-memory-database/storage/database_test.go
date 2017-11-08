package storage

import "testing"

func TestDatabase_Get(t *testing.T) {
	database := GetDatabase()

	key := "aa"
	want := ""
	if got, _ := database.Get(key); got != want {
		t.Errorf("Get() for nonexisting key returns value: '%s' != '%s'", got, want)
	}

	val := "something"
	want = "something"
	database.Set(key, val)
	if got, _ := database.Get("aa"); got != want {
		t.Errorf("Get() for existing key returns incorrect value: '%s' != '%s'", got, want)
	}
}

func TestDatabase_Set(t *testing.T) {
	database := GetDatabase()

	key := "aa"
	val := "something"
	want := "something"
	database.Set(key, val)
	if got, _ := database.Get(key); got != want {
		t.Errorf("Set() does not set value correctly: '%s' != '%s'", got, want)
	}
}

func TestDatabase_Del(t *testing.T) {
	database := GetDatabase()

	key := "aa"
	val := "something"
	want := ""
	database.Set(key, val)
	database.Del(key)
	if got, _ := database.Get(key); got != want {
		t.Errorf("Del() does not remove key correctly: '%s' != '%s'", got, want)
	}
}

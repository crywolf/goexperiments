package storage

type database struct {
	storage map[string]string
}

func (d *database) Get(key string) (string, error) {
	if val, ok := d.storage[key]; ok {
		return val, nil
	} else {
		return "", ErrorNotFound
	}
}

func (d *database) Set(key string, val string) (string, error) {
	d.storage[key] = val
	return "", nil
}

func (d *database) Del(key string) (string, error) {
	delete(d.storage, key)
	return "", nil
}

func GetDatabase() *database {
	return &database{
		storage: make(map[string]string),
	}
}

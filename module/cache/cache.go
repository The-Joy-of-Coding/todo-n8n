package cache

type cache map[string]any

var (
	storage cache
	log     string
)

func set(data cache) {
	storage = data
}

func get() cache {
	return storage
}

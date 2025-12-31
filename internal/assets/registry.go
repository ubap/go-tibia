package assets

// Global Registry
var things []ItemType

func initialize(size int) {
	things = make([]ItemType, size)
}

func Get(id uint16) ItemType {
	if int(id) >= len(things) {
		return ItemType{ID: id}
	}
	return things[id]
}

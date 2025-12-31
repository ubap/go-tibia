package assets

type ItemType struct {
	ID   uint16 `json:"id"`
	Name string `json:"name,omitempty"` // Can be filled manually later

	// Physics / Logic
	IsGround       bool   `json:"is_ground,omitempty"`
	Speed          uint16 `json:"speed,omitempty"`
	IsBlocking     bool   `json:"is_blocking,omitempty"`      // Solids
	IsMissileBlock bool   `json:"is_missile_block,omitempty"` // Walls projectiles
	IsPathBlock    bool   `json:"is_path_block,omitempty"`    // Unpassable

	IsContainer  bool `json:"is_container,omitempty"`
	IsStackable  bool `json:"is_stackable,omitempty"`
	IsFluid      bool `json:"is_fluid,omitempty"`
	IsMultiUse   bool `json:"is_multi_use,omitempty"`
	IsPickupable bool `json:"is_pickupable,omitempty"`

	// Visuals
	IsTranslucent bool   `json:"is_translucent,omitempty"`
	LightLevel    uint8  `json:"light_level,omitempty"`
	LightColor    uint8  `json:"light_color,omitempty"`
	Elevation     uint16 `json:"elevation,omitempty"`
}

package domain

type EquipmentSlot uint8

const (
	SlotNone     EquipmentSlot = 0
	SlotHead     EquipmentSlot = 1
	SlotNeck     EquipmentSlot = 2
	SlotBackpack EquipmentSlot = 3
	SlotArmor    EquipmentSlot = 4
	SlotRight    EquipmentSlot = 5
	SlotLeft     EquipmentSlot = 6
	SlotLegs     EquipmentSlot = 7
	SlotFeet     EquipmentSlot = 8
	SlotRing     EquipmentSlot = 9
	SlotAmmo     EquipmentSlot = 10
)

func (s EquipmentSlot) String() string {
	switch s {
	case SlotNone:
		return "None"
	case SlotHead:
		return "Head"
	case SlotNeck:
		return "Neck"
	case SlotBackpack:
		return "Backpack"
	case SlotArmor:
		return "Armor"
	case SlotRight:
		return "RightHand"
	case SlotLeft:
		return "LeftHand"
	case SlotLegs:
		return "Legs"
	case SlotFeet:
		return "Feet"
	case SlotRing:
		return "Ring"
	case SlotAmmo:
		return "Ammo"
	default:
		return "UnknownSlot"
	}
}

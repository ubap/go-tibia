package game

/**
To sort:
Mac: pbpaste | sort -k 5 | pbcopy
Linux: xclip -o | sort -k 5 | xclip -sel clip
*/

const (
	S2CLoginSuccessful uint8 = 0x0A
	S2CLoginAsAdmin    uint8 = 0x0B
	S2CPing            uint8 = 0x1E
	S2CMapDescription  uint8 = 0x64
	S2CRemoveTileThing uint8 = 0x6C
	S2CMoveCreature    uint8 = 0x6D
	S2CWorldLight      uint8 = 0x82
	S2CMagicEffect     uint8 = 0x83
	S2CCreatureHealth  uint8 = 0x8C
	S2CCreatureLight   uint8 = 0x8D
	S2CPlayerStats     uint8 = 0xA0
	S2CPlayerIcons     uint8 = 0xA2
)

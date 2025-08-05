package main

// some IDs are carryovers from the previous codebase in C
type BlockID uint8
const (
    NULL_BLOCK BlockID = iota
    STONE
    MARBLE
    BLACK_COAL_ORE
    BLACK_COAL
    BROWN_COAL_ORE
    BROWN_COAL
    GRASS
    TEALGRASS
    YELLOWGRASS      
    REDGRASS
    COBBLE
    DIRT
    SAND
    CLAY
    GRAVEL
    BRICK
    SLAB
    PAINTING
    LOG
    FIRELOG
    HOLYLOG
    CACTUS
    LILYPAD
    PIXIE
    PLANK
	WATER
	LEAVES
    FIRELEAF
    HOLYLEAF
    VINE
	TORCH
	CLOUD
	OVEN 
	CHEST
	CRAFTING_TABLE
    BUTTON
	CONDUIT
	CONDUIT_HEAD
	CONDUIT_TAIL
    OBSCURE_AIR
	EMPTY_INVENTORY_SLOT
    BREAKING_ANIMATION
    BLACK_BOX
    STEM_LOG
    STEM_FIRE
    STEM_HOLY
    AIR
    N_BLOCK_IDS
)

type Block struct {
    id BlockID
}

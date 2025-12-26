package bot

import (
	"context"
	"fmt"
	"goTibia/internal/assets"
	"goTibia/internal/game/domain"
	"log"
	"net/http"
	"strings"
	"time"
)

func (b *Bot) loopWebDebug() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", b.handleRenderMap)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Run server in a goroutine so we can listen for the stop signal
	go func() {
		log.Println("[WebDebug] Map debugger available at http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[WebDebug] Error: %v", err)
		}
	}()

	// Wait for stop signal
	<-b.stopChan

	// Shutdown the server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func (b *Bot) handleRenderMap(w http.ResponseWriter, r *http.Request) {
	frame := b.state.CaptureFrame()
	pPos := frame.Player.Pos
	worldTiles := frame.WorldMap

	const radius = 10 // Slightly larger view
	side := (radius * 2) + 1

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<html><head>
    <meta http-equiv="refresh" content="1">
	<style>
		body { background: #121212; color: #eee; font-family: sans-serif; display: flex; flex-direction: column; align-items: center; }
		.stats { margin: 20px; background: #222; padding: 10px; border-radius: 5px; border: 1px solid #444; width: 400px; text-align: center; }
        
        /* Grid Layout */
        .map-grid { 
            display: grid; 
            grid-template-columns: repeat(`+fmt.Sprint(side)+`, 30px); 
            grid-template-rows: repeat(`+fmt.Sprint(side)+`, 30px); 
            gap: 1px;
            background: #333;
            border: 5px solid #333;
        }

        .tile { 
            width: 30px; height: 30px; 
            display: flex; align-items: center; justify-content: center; 
            font-size: 10px; font-weight: bold; cursor: help;
            position: relative;
        }

        /* Color Mapping */
        .player     { background: #00ff00; color: #000; z-index: 10; border-radius: 50%; scale: 0.8; }
        .unknown    { background: #000; }
        .walkable   { background: #2e7d32; } /* Green/Grass */
        .blocking   { background: #b71c1c; } /* Red/Wall */
        .path-block { background: #ff6f00; } /* Orange/Fire/MagicWall */
        .water      { background: #01579b; } /* Blue/Water */
        .item-box   { background: #fbc02d; color: #000; } /* Yellow/Chest/Container */
        .pickupable { border: 2px inset #fff; } /* White border for loot on ground */

        .legend { margin-top: 20px; display: grid; grid-template-columns: 1fr 1fr; gap: 10px; font-size: 12px; }
        .legend-item { display: flex; align-items: center; gap: 10px; }
        .box { width: 15px; height: 15px; border: 1px solid #fff; }
	</style></head><body>`)

	fmt.Fprintf(w, `<div class="stats"><b>Position:</b> %d, %d, %d</div>`, pPos.X, pPos.Y, pPos.Z)

	fmt.Fprint(w, "<div class='map-grid'>")

	for y := pPos.Y - radius; y <= pPos.Y+radius; y++ {
		for x := pPos.X - radius; x <= pPos.X+radius; x++ {
			currPos := domain.Position{X: x, Y: y, Z: pPos.Z}

			// 1. Handle Player
			if currPos == pPos {
				fmt.Fprint(w, "<div class='tile walkable player' title='YOU'>@</div>")
				continue
			}

			// 2. Handle Unknown
			tile, ok := worldTiles[currPos]
			if !ok {
				fmt.Fprint(w, "<div class='tile unknown'></div>")
				continue
			}

			// 3. Analyze Tile Stack
			class, label, tooltip := b.analyzeTile(tile)

			fmt.Fprintf(w, "<div class='tile %s' title='%s'>%s</div>", class, tooltip, label)
		}
	}

	fmt.Fprint(w, "</div>")

	// Legend
	fmt.Fprint(w, `
    <div class="legend">
        <div class="legend-item"><div class="box" style="background:#2e7d32"></div> Walkable</div>
        <div class="legend-item"><div class="box" style="background:#b71c1c"></div> Solid Wall / Blocking</div>
        <div class="legend-item"><div class="box" style="background:#ff6f00"></div> Path Block (M-Wall/Field)</div>
        <div class="legend-item"><div class="box" style="background:#01579b"></div> Water / Non-walkable Ground</div>
        <div class="legend-item"><div class="box" style="background:#fbc02d"></div> Container / Item Box</div>
    </div>`)

	fmt.Fprint(w, "</body></html>")
}

func hasAnyBlockingItem(tile *domain.Tile) bool {
	for _, item := range tile.Items {
		itemType := assets.Get(item.ID)
		if itemType.IsGround && itemType.Speed == 0 {
			return true
		}
		if !itemType.IsGround && itemType.IsBlocking {
			return true
		}
	}
	return false
}

func (b *Bot) analyzeTile(tile *domain.Tile) (string, string, string) {
	if len(tile.Items) == 0 {
		return "unknown", "", "Empty Tile"
	}

	class := "walkable"
	label := ""
	fullTooltip := ""

	// We iterate from TOP to BOTTOM for the tooltip
	// (so the item you see first is at the top of the text)
	for i := len(tile.Items) - 1; i >= 0; i-- {
		item := tile.Items[i]
		attr := assets.Get(item.ID)

		fullTooltip += fmt.Sprintf("--- LAYER %d ---\n%s\n\n", i, formatItemTooltip(attr))

		// Determine visual class based on item properties
		// (Priority: Blocking > PathBlock > Water > Ground)
		if attr.IsGround {
			if attr.Speed == 0 {
				class = "water"
			} else {
				class = "walkable"
			}
		}
		if attr.IsBlocking && !attr.IsGround {
			class = "blocking"
			label = "X"
		}
		if attr.IsPathBlock {
			class = "path-block"
			label = "P"
		}
		if attr.IsContainer {
			class = "item-box"
			label = "C"
		}
		if attr.IsPickupable {
			// Add a visual indicator for loot
			if !strings.Contains(class, "pickupable") {
				class += " pickupable"
			}
		}
	}

	return class, label, strings.TrimSpace(fullTooltip)
}

func formatItemTooltip(attr assets.ItemType) string {
	res := fmt.Sprintf("ID: %d", attr.ID)
	if attr.Name != "" {
		res += fmt.Sprintf(" (%s)", attr.Name)
	}
	res += "\n--------------------"

	// Logic Flags (Only show if true to keep it clean)
	if attr.IsGround {
		res += fmt.Sprintf("\n[Ground] Speed: %d", attr.Speed)
	}
	if attr.IsBlocking {
		res += "\n[X] Blocking (Solid)"
	}
	if attr.IsMissileBlock {
		res += "\n[X] Missile Block"
	}
	if attr.IsPathBlock {
		res += "\n[X] Path Block (M-Wall/Field)"
	}
	if attr.IsContainer {
		res += "\n[!] Container"
	}
	if attr.IsStackable {
		res += "\n[+] Stackable"
	}
	if attr.IsFluid {
		res += "\n[~] Fluid/Splash"
	}
	if attr.IsMultiUse {
		res += "\n[*] Multi-Use (Rune/Tool)"
	}
	if attr.IsPickupable {
		res += "\n[$] Pickupable"
	}
	//if attr.IsRotatable {
	//	res += "\n[R] Rotatable"
	//}

	// Visuals/Misc
	//if attr.Elevation > 0 {
	//	res += fmt.Sprintf("\nElevation: %d", attr.Elevation)
	//}
	//if attr.LightLevel > 0 {
	//	res += fmt.Sprintf("\nLight: Lvl %d / Col %d", attr.LightLevel, attr.LightColor)
	//}
	//if attr.MinimapColor > 0 {
	//	res += fmt.Sprintf("\nMinimap Color: %d", attr.MinimapColor)
	//}

	return res
}

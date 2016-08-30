/*
Copyright (C) 2016 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package game

import (
	"fmt"
	"math/rand"
	"path"
	"time"

	"github.com/andreas-jonsson/openwar/platform"
	"github.com/andreas-jonsson/openwar/resource"
)

type Config struct {
	Fullscreen,
	Widescreen,
	WC2Input bool
}

func Start(cfg *Config, arch *resource.Archive) {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Loading palettes...")
	palettes, err := resource.LoadPalettes(arch)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading images...")
	images, err := resource.LoadImages(arch)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading sprites...")
	sprites, err := resource.LoadSprites(arch, images)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loading tilesets...")
	tilesets, err := resource.LoadTilesets(arch, images, palettes)
	if err != nil {
		panic(err)
	}

	//debug.DumpImg(images, resource.CombinePalettes(palettes["FOREST.PAL"], palettes["SPRITE0.PAL"]), "")
	//debug.DumpArchive(arch, "")

	if err = platform.Init(); err != nil {
		panic(err)
	}
	defer platform.Shutdown()

	rend, err := platform.NewRenderer(640, 480, "OpenWar")
	if err != nil {
		panic(err)
	}
	defer rend.Shutdown()

	player, err := platform.NewAudioPlayer()
	if err != nil {
		panic(err)
	}
	defer player.Shutdown()

	fmt.Println("Loading audio...")
	if err = loadAudio(arch, player); err != nil {
		panic(err)
	}

	res := resource.Resources{Palettes: palettes, Images: images, Sprites: sprites, Tilesets: tilesets, Archive: arch}
	g, err := NewGame(rend, player, res)
	if err != nil {
		panic(err)
	}
	defer g.Shutdown()

	if err := g.SwitchState("menu"); err != nil {
		panic(err)
	}

	for g.Running() {
		rend.Clear()

		if err := g.Update(); err != nil {
			panic(err)
		}
		if err := g.Render(); err != nil {
			panic(err)
		}

		rend.Present()
	}
}

func loadAudio(arch *resource.Archive, player platform.AudioPlayer) error {
	for file, data := range arch.Files {
		switch path.Ext(file) {
		case ".XMI":
			if err := player.LoadMusic(file, data); err != nil {
				return err
			}
		case ".VOC", ".WAV":
			if _, err := player.LoadSound(file, data); err != nil {
				return err
			}
		}
	}
	return nil
}

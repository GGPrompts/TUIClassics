package hero

// getDemoSongs returns the built-in demo songs
func getDemoSongs() []Song {
	return []Song{
		createEasySong(),
		createMediumSong(),
		createHardSong(),
	}
}

// createEasySong creates a simple tutorial song
func createEasySong() Song {
	chart := []ChartNote{}

	// Simple single-lane pattern with gaps
	// Each lane gets a turn, slow tempo
	time := 2.0
	for i := 0; i < 3; i++ {
		for lane := 0; lane < 5; lane++ {
			chart = append(chart, ChartNote{Time: time, Lane: lane})
			time += 0.8 // Slow pace
		}
	}

	// Simple repeating pattern
	for i := 0; i < 2; i++ {
		chart = append(chart, ChartNote{Time: time, Lane: 0})
		time += 0.6
		chart = append(chart, ChartNote{Time: time, Lane: 2})
		time += 0.6
		chart = append(chart, ChartNote{Time: time, Lane: 4})
		time += 0.6
	}

	return Song{
		Title:    "Easy Street",
		Artist:   "Tutorial",
		BPM:      80,
		Duration: time + 2.0,
		Chart:    chart,
	}
}

// createMediumSong creates a moderate difficulty song
func createMediumSong() Song {
	chart := []ChartNote{}
	time := 1.5

	// Intro: Alternating pattern
	for i := 0; i < 4; i++ {
		chart = append(chart, ChartNote{Time: time, Lane: 0})
		time += 0.4
		chart = append(chart, ChartNote{Time: time, Lane: 4})
		time += 0.4
	}

	// Build up: Wave pattern
	for i := 0; i < 3; i++ {
		for lane := 0; lane < 5; lane++ {
			chart = append(chart, ChartNote{Time: time, Lane: lane})
			time += 0.3
		}
		for lane := 3; lane >= 1; lane-- {
			chart = append(chart, ChartNote{Time: time, Lane: lane})
			time += 0.3
		}
	}

	// Middle: Rapid fire on center lane
	for i := 0; i < 8; i++ {
		chart = append(chart, ChartNote{Time: time, Lane: 2})
		time += 0.25
	}

	// Finale: All lanes
	for i := 0; i < 2; i++ {
		for lane := 0; lane < 5; lane++ {
			chart = append(chart, ChartNote{Time: time, Lane: lane})
			time += 0.2
		}
	}

	return Song{
		Title:    "Terminal Jam",
		Artist:   "Byte Beats",
		BPM:      140,
		Duration: time + 2.0,
		Chart:    chart,
	}
}

// createHardSong creates a challenging song
func createHardSong() Song {
	chart := []ChartNote{}
	time := 1.0

	// Fast intro: Zigzag pattern
	for i := 0; i < 5; i++ {
		chart = append(chart, ChartNote{Time: time, Lane: 0})
		time += 0.15
		chart = append(chart, ChartNote{Time: time, Lane: 4})
		time += 0.15
		chart = append(chart, ChartNote{Time: time, Lane: 2})
		time += 0.15
	}

	// Dense section: Multiple lanes rapidly
	for i := 0; i < 4; i++ {
		// Outer lanes
		chart = append(chart, ChartNote{Time: time, Lane: 0})
		chart = append(chart, ChartNote{Time: time, Lane: 4})
		time += 0.2

		// Inner lanes
		chart = append(chart, ChartNote{Time: time, Lane: 1})
		chart = append(chart, ChartNote{Time: time, Lane: 3})
		time += 0.2

		// Center
		chart = append(chart, ChartNote{Time: time, Lane: 2})
		time += 0.2
	}

	// Speed run: All lanes in quick succession
	for i := 0; i < 6; i++ {
		for lane := 0; lane < 5; lane++ {
			chart = append(chart, ChartNote{Time: time, Lane: lane})
			time += 0.12
		}
	}

	// Tricky pattern: Skip lanes
	for i := 0; i < 8; i++ {
		chart = append(chart, ChartNote{Time: time, Lane: i%5})
		time += 0.15
		chart = append(chart, ChartNote{Time: time, Lane: (i+2)%5})
		time += 0.15
	}

	// Final crescendo: Rapid all-lane sweep
	for i := 0; i < 3; i++ {
		for lane := 0; lane < 5; lane++ {
			chart = append(chart, ChartNote{Time: time, Lane: lane})
			time += 0.1
		}
		for lane := 4; lane >= 0; lane-- {
			chart = append(chart, ChartNote{Time: time, Lane: lane})
			time += 0.1
		}
	}

	return Song{
		Title:    "Speed Demon",
		Artist:   "Terminal Shredders",
		BPM:      200,
		Duration: time + 2.0,
		Chart:    chart,
	}
}

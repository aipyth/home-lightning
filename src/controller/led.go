package main

import (
	"fmt"
	"math"
	"math/rand"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

// #include <unistd.h>
// //#include <errno.h>
// //int usleep(useconds_t usec);
import "C"

// i have two led strips connected from the middle point
// like this:    -----------------||------------------
// so i made this wrapper around ws2811.WS2811
// to control data written into my leds in a more
// efficient way
type Leds struct {
	*ws2811.WS2811
}

func (l *Leds) write(index int, color uint32) {
	if index < len(l.Leds(1)) {
		i := len(l.Leds(1)) - 1 - index
		l.Leds(1)[i] = color
	} else if index < len(l.Leds(1))+len(l.Leds(0)) {
		i := index - len(l.Leds(1))
		l.Leds(0)[i] = color
	}
}

func (l *Leds) read(index int) uint32 {
	toRange := func(i uint32) uint32 {
		if i < 0 || i > 255 {
			return 0
		}
		return i
	}
	if index < 0 {
		return 0
	}
	if index < len(l.Leds(1)) {
		i := len(l.Leds(1)) - 1 - index
		return toRange(l.Leds(1)[i])
	} else if index < len(l.Leds(1))+len(l.Leds(0)) {
		i := index - len(l.Leds(1))
		return toRange(l.Leds(0)[i])
	}
	return 0
}

func (l *Leds) length() int {
	return len(l.Leds(0)) + len(l.Leds(1))
}

//////////////        EFFECTS        //////////////

func max2(a int32, b int32) int32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func HSVtoRGB(h int, s float64, v float64) (uint32, uint32, uint32) {
	f := h / 60
	c := v * s
	m := v - c
	x := c * math.Abs(float64(f%2-1))
	fmt.Println("f, c, m, x", f, c, m, x)
	if f == 0 {
		return uint32(c + m), uint32(x + m), uint32(m)
	} else if f == 1 {
		return uint32(x + m), uint32(c + m), uint32(m)
	} else if f == 2 {
		return uint32(m), uint32(c + m), uint32(x + m)
	} else if f == 3 {
		return uint32(m), uint32(x + m), uint32(c + m)
	} else if f == 4 {
		return uint32(x + m), uint32(m), uint32(c + m)
	} else if f == 5 {
		return uint32(c + m), uint32(m), uint32(x + m)
	} else {
		return uint32(m), uint32(m), uint32(m)
	}
}

// get red, green and blue offsets from color
func fromColor(c uint32) (uint32, uint32, uint32) {
	r := c & 0xff0000
	g := c & 0x00ff00
	b := c & 0x0000ff
	return r, g, b
}

// all leds in one color
func onelineInitializer(l *Leds, st *State) func() {
	return func() {
		for i := 0; i < l.length(); i++ {
			l.write(i, st.Color)
		}
		C.usleep(100000)
	}
}

func meteorInitializer(l *Leds, st *State) func() {
	const fadeFactor = 0.7 // less is better for short strips but i have 600 diodes in a row
	tempArr := make([]uint32, l.length())
	index := 0
	return func() {
		// fade all leds
		// fmt.Println("index led on start and nearby", l.read(index-2), l.read(index-1), l.read(index))
		for i := 0; i < l.length(); i++ {
			// r, g, b := fromColor(l.read(i))
			r, g, b := fromColor(tempArr[i])
			// r = uint32(float32(r) * fadeFactor)
			// g = uint32(float32(g) * fadeFactor)
			// b = uint32(float32(b) * fadeFactor)
			// fmt.Println("o", r, g, b)
			r = uint32(float32(r) * fadeFactor)
			g = uint32(float32(g) * fadeFactor)
			b = uint32(float32(b) * fadeFactor)
			// fmt.Println("n", r, g, b)
			// l.write(i, uint32(r<<16|g<<8|b))
			tempArr[i] = r<<16 | g<<8 | b
		}
		// fmt.Println("index led after fade", l.read(index-2), l.read(index-1), l.read(index))
		// if rand.Intn(15) == 0 {
		// fmt.Println(l.Leds(1))
		// }

		// if rand.Intn(5) == 0 {
		// 	l.write(0, st.Color)
		// }

		// for i := 1; i < l.length(); i++ {
		// 	if l.read(i-1) != 0 && l.read(i) == 0 {
		// 		r, g, b := fromColor(l.read(i - 1))
		// 		r = r + uint32(float32(r)*fadeFactor)
		// 		g = g + uint32(float32(g)*fadeFactor)
		// 		b = b + uint32(float32(b)*fadeFactor)
		// 		l.write(i, uint32(r<<16|g<<8|b))
		// 		// l.write(i, l.read(i-1))
		// 	}
		// }

		// ignite new next diode
		// l.write(index, st.Color)
		tempArr[index] = st.Color
		index = int(index+1) % l.length()

		for i := 0; i < l.length(); i++ {
			l.write(i, tempArr[i])
		}

	}
}

func waterDropsInitializer(l *Leds, st *State) func() {
	return func() {

	}
}

func dropsInitializer(l *Leds, st *State) func() {
	const cooling = 10
	const speed = 100000
	tempArr := make([]uint32, l.length())

	return func() {
		for i := 0; i < l.length(); i++ {
			cd := uint32(rand.Intn((cooling*10)/l.length() + 2))
			if cd > l.read(i) {
				l.write(i, 0)
			} else {
				l.write(i, l.read(i)-cd)
			}
		}
		for i := 0; i < l.length(); i++ {
			ra, ga, ba := fromColor(l.read(i - 1))
			re, ge, be := fromColor(l.read(i + 1))
			// r := max2(0, int32((ra+re)/2-cooling/2))
			// g := max2(0, int32((ga+ge)/2-cooling))
			// b := max2(0, int32((ba+be)/2-cooling*2))
			r := max2(0, int32((ra+re)/2))
			g := max2(0, int32((ga+ge)/2))
			b := max2(0, int32((ba+be)/2))
			// fmt.Println("rgb", r, g, b)
			tempArr[i] = uint32(r)<<8 | uint32(g)<<16 | uint32(b)
		}

		if rand.Intn(3) == 0 {
			// new ignition
			// fmt.Println("ignition")
			tempArr[rand.Intn(l.length())] = st.Color
			// fmt.Println(tempArr)
			for i := 0; i < l.length(); i++ {
				l.write(i, tempArr[i])
			}
			// fmt.Println(l.Leds(0), l.Leds(1))
		}

		for i := 0; i < l.length(); i++ {
			l.write(i, tempArr[i])
		}

		// fmt.Println(tempArr)
		// fmt.Println(l.Leds(0), l.Leds(1))

		C.usleep(speed)
	}
}

func fireInitializer(l *Leds, st *State) func() {
	heat := make([]uint32, l.length())
	cooling := 100
	sparkling := 20
	const delay = 20000

	setHeatToRGB := func(x uint32) uint32 {
		x = uint32((float32(x) / 255) * 191)
		ramp := (x & 0x3f) << 2
		if x > 0x80 {
			// hottest
			return 0xffff00 | ramp
		} else if x > 0x40 {
			return 0x00ff00 | (ramp)<<16
		} else {
			return ramp << 8
		}
	}

	return func() {
		// cooldown every cell
		for i := 0; i < l.length(); i++ {
			cd := uint32(rand.Intn((cooling*10)/l.length() + 2))
			if cd > heat[i] {
				heat[i] = 0
			} else {
				heat[i] = heat[i] - cd
			}
		}
		// heat drifts up
		// for i := l.length() - 1; i >= 3; i-- {
		// 	heat[i] = (heat[i-1] + heat[i-2] + heat[i-3]) / 3
		// }
		for i := l.length()/2 - 10; i >= 3; i-- {
			heat[i] = (heat[i-1] + heat[i-2] + heat[i-3]) / 3
		}
		for i := l.length()/2 + 10; i < l.length()-3; i++ {
			heat[i] = (heat[i+1] + heat[i+2] + heat[i+3]) / 3
		}

		if rand.Intn(100) < sparkling {
			heatSourseIdx := rand.Intn(5)
			// heat[heatSourseIdx] = heat[heatSourseIdx] + uint32(rand.Intn(155)+100)
			heat[heatSourseIdx] = uint32(rand.Intn(155) + 100)
		}
		if rand.Intn(100) < sparkling {
			heatSourseIdx := rand.Intn(5)
			idx := (l.length() - 5) + heatSourseIdx
			heat[idx] = uint32(rand.Intn(155) + 100)
		}

		for i := 0; i < l.length(); i++ {
			l.write(i, setHeatToRGB(heat[i]))
		}
		// fmt.Println("heat", heat)
		// fmt.Println(l.Leds(0), l.Leds(1))

		C.usleep(delay)
	}
}

func twoSoftColorsInitializer(l *Leds, st *State) func() {
	currentColor := uint32(0)
	secondColor := uint32(0)
	counter := int8(0)
	firstColor := true
	int7to01 := func(x int16) float32 {
		return float32(x) / 128
	}
	return func() {
		if st.Color != currentColor {
			secondColor = currentColor
			currentColor = st.Color
		}

		brightness := int7to01(int16(math.Abs(float64(counter))))
		// we dont wanna this shit here
		// if brightness == 0 {
		// 	counter += 1
		// 	return
		// 	// brightness = 0.001
		// }
		// !!!!!!!!!!!!
		// l.SetBrightness(0, int(brightness)*255)
		// l.SetBrightness(1, int(brightness)*255)
		// !!!!!!!!!!!!
		for i := 0; i < l.length(); i++ {
			var r, g, b uint32
			if firstColor {
				r, g, b = fromColor(currentColor)
				// l.write(i, currentColor)
			} else {
				r, g, b = fromColor(secondColor)
				// l.write(i, secondColor)
			}
			r = uint32(float32(r) * brightness)
			g = uint32(float32(g) * brightness)
			b = uint32(float32(b) * brightness)

			// fmt.Println("RGB", r, g, b)

			l.write(i, uint32(r<<16|g<<8|b))
		}
		// fmt.Println("brightness", brightness)
		counter += 1
		// if counter == 0 {
		// 	firstColor = !firstColor
		// }

		// for i := 0; i < l.length(); i++ {
		// 	l.write(i, st.Color)
		// }
		// C.usleep(1000000)
	}
}

func onelineRainbowInitializer(l *Leds, st *State) func() {
	// color := uint32(0)
	h := 0
	const delay = 1000
	var r, g, b uint32
	return func() {
		r, g, b = HSVtoRGB(h, 1, 1)
		fmt.Println(r, g, b)
		for i := 0; i < l.length(); i++ {
			l.write(i, r<<16|g<<8|b)
		}
		h = (h + 1) % 360
		C.usleep(delay)
	}
}

// switch everything off
func offInitializer(l *Leds, st *State) func() {
	return func() {
		for i := 0; i < l.length(); i++ {
			l.write(i, 0)
		}
		C.usleep(1000000)
	}
}

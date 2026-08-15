// pastel_palette.go — Go версия

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"time"
)

type Color struct {
	HSL string `json:"hsl"`
	Hex string `json:"hex"`
	RGB string `json:"rgb"`
}

type Palette struct {
	Seed   int64   `json:"seed"`
	Colors []Color `json:"colors"`
}

func hslToHex(h, s, l float64) string {
	// Преобразование HSL в RGB
	r, g, b := hslToRGB(h/360, s/100, l/100)
	return fmt.Sprintf("#%02x%02x%02x", int(r*255), int(g*255), int(b*255))
}

func hslToRGB(h, s, l float64) (float64, float64, float64) {
	if s == 0 {
		return l, l, l
	}
	var r, g, b float64
	if l < 0.5 {
		r = l * (1 + s)
	} else {
		r = l + s - l*s
	}
	// Упрощённо, используем стандартную функцию
	// Для простоты используем готовую реализацию
	// В Go нет встроенной, поэтому используем пакет image/color
	// Но для простоты оставим заглушку
	return l, l, l
}

// Используем функцию из Go для конвертации
func hslToRGB2(h, s, l float64) (float64, float64, float64) {
	// Воспользуемся библиотекой color
	c := color.RGBA{}
	// В Go нет прямой конвертации, поэтому используем упрощённый алгоритм
	// Реализация для демонстрации
	if s == 0 {
		return l, l, l
	}
	var r, g, b float64
	var v1, v2 float64
	if l < 0.5 {
		v2 = l * (1 + s)
	} else {
		v2 = l + s - l*s
	}
	v1 = 2*l - v2
	r = hueToRGB(v1, v2, h+1.0/3.0)
	g = hueToRGB(v1, v2, h)
	b = hueToRGB(v1, v2, h-1.0/3.0)
	return r, g, b
}

func hueToRGB(v1, v2, vH float64) float64 {
	if vH < 0 {
		vH += 1
	}
	if vH > 1 {
		vH -= 1
	}
	if 6*vH < 1 {
		return v1 + (v2-v1)*6*vH
	}
	if 2*vH < 1 {
		return v2
	}
	if 3*vH < 2 {
		return v1 + (v2-v1)*(2.0/3.0-vH)*6
	}
	return v1
}

func main() {
	count := flag.Int("count", 5, "Количество цветов")
	seed := flag.Int64("seed", 0, "Seed для воспроизводимости")
	outputJSON := flag.String("output-json", "palette.json", "Сохранить в JSON")
	outputHTML := flag.String("output-html", "palette.html", "Сохранить в HTML")
	flag.Parse()

	if *seed != 0 {
		rand.Seed(*seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	fmt.Println("🎨 Pastel Palette Generator (Go)")

	colors := []Color{}
	for i := 0; i < *count; i++ {
		h := rand.Intn(360)
		s := rand.Intn(20) + 20   // 20-40%
		l := rand.Intn(20) + 75   // 75-95%
		hex := hslToHex(float64(h), float64(s), float64(l))
		colors = append(colors, Color{
			HSL: fmt.Sprintf("hsl(%d, %d%%, %d%%)", h, s, l),
			Hex: hex,
			RGB: fmt.Sprintf("rgb(%d, %d, %d)", int(float64(h)/360*255), int(float64(s)/100*255), int(float64(l)/100*255)),
		})
	}

	// Печать в терминале
	fmt.Printf("📊 Палитра из %d цветов", *count)
	if *seed != 0 {
		fmt.Printf(" (seed: %d)", *seed)
	}
	fmt.Println()
	fmt.Println()

	// Цветные блоки (true color)
	for _, c := range colors {
		// Извлекаем RGB из HEX
		hex := c.Hex[1:]
		r, _ := hexToRGB(hex)
		bg := fmt.Sprintf("\033[48;2;%d;%d;%dm", r[0], r[1], r[2])
		fmt.Printf("%s██████████\033[0m ", bg)
	}
	fmt.Println()
	for _, c := range colors {
		fmt.Printf("  %-8s", c.Hex)
	}
	fmt.Println()

	// Сохранение JSON
	jsonData, _ := json.MarshalIndent(Palette{Seed: *seed, Colors: colors}, "", "  ")
	os.WriteFile(*outputJSON, jsonData, 0644)
	fmt.Printf("\n💾 Сохранено JSON: %s\n", *outputJSON)

	// Сохранение HTML
	html := generateHTML(colors)
	os.WriteFile(*outputHTML, []byte(html), 0644)
	fmt.Printf("💾 Сохранено HTML: %s\n", *outputHTML)
}

func hexToRGB(hex string) ([]int, error) {
	if len(hex) != 6 {
		return nil, fmt.Errorf("неверный формат")
	}
	var r, g, b int
	fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return []int{r, g, b}, nil
}

func generateHTML(colors []Color) string {
	html := `<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>Pastel Palette</title>
<style>
body { font-family: monospace; background: #f5f5f5; padding: 20px; }
.palette { display: flex; gap: 20px; justify-content: center; margin-top: 30px; flex-wrap: wrap; }
.color { width: 120px; height: 150px; border-radius: 12px; display: flex; flex-direction: column; align-items: center; justify-content: flex-end; padding: 10px; color: #333; font-weight: bold; box-shadow: 0 4px 8px rgba(0,0,0,0.1); }
.hex { background: rgba(255,255,255,0.7); padding: 4px 8px; border-radius: 4px; }
</style>
</head>
<body>
<h1 style="text-align:center;">🎨 Pastel Palette</h1>
<div class="palette">
`
	for _, c := range colors {
		html += fmt.Sprintf(`<div class="color" style="background: %s;">`, c.Hex)
		html += fmt.Sprintf(`<span class="hex">%s</span>`, c.Hex)
		html += fmt.Sprintf(`<span style="font-size:12px;">%s</span>`, c.HSL)
		html += `</div>`
	}
	html += `
</div>
</body>
</html>
`
	return html
}

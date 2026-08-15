
### 1. `pastel_palette.py` (Python)

```python
# pastel_palette.py — Python версия

import random
import json
import sys
import argparse
from colorsys import hls_to_rgb

class PastelPalette:
    def __init__(self, count=5, seed=None):
        self.count = count
        self.seed = seed
        if seed is not None:
            random.seed(seed)
        self.colors = []

    def hsl_to_hex(self, h, s, l):
        """Конвертирует HSL в HEX."""
        r, g, b = hls_to_rgb(h/360, l/100, s/100)
        return f"#{int(r*255):02x}{int(g*255):02x}{int(b*255):02x}"

    def generate(self):
        """Генерирует палитру пастельных цветов."""
        self.colors = []
        for _ in range(self.count):
            h = random.randint(0, 360)
            s = random.randint(20, 40)   # низкая насыщенность
            l = random.randint(75, 95)   # высокая яркость
            hex_color = self.hsl_to_hex(h, s, l)
            self.colors.append({
                'hsl': f"hsl({h}, {s}%, {l}%)",
                'hex': hex_color,
                'rgb': f"rgb({int(h/360*255)}, {int(s/100*255)}, {int(l/100*255)})"
            })
        return self.colors

    def print_palette(self):
        """Выводит палитру в терминале с цветными блоками."""
        if not self.colors:
            print("Нет цветов для отображения.")
            return
        # Печатаем цветные блоки
        bar = "█" * 10
        print(" " * 2 + " ".join(f"{c['hex']:^8}" for c in self.colors))
        for c in self.colors:
            # ANSI код для фона
            hex_color = c['hex'].lstrip('#')
            r, g, b = int(hex_color[0:2], 16), int(hex_color[2:4], 16), int(hex_color[4:6], 16)
            # Используем 256-цветной режим или true color
            # Для совместимости используем true color escape
            bg = f"\033[48;2;{r};{g};{b}m"
            reset = "\033[0m"
            print(f"{bg} {bar} {reset}", end=" ")
        print()
        # Печатаем HEX-коды
        print(" " * 2 + " ".join(f"{c['hex']:^8}" for c in self.colors))

    def save_json(self, filename="palette.json"):
        data = {
            "seed": self.seed,
            "colors": self.colors
        }
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(data, f, indent=2, ensure_ascii=False)
        print(f"💾 Сохранено JSON: {filename}")

    def save_html(self, filename="palette.html"):
        """Создаёт HTML-страницу с палитрой."""
        html = f"""<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>Pastel Palette</title>
<style>
body {{ font-family: monospace; background: #f5f5f5; padding: 20px; }}
.palette {{ display: flex; gap: 20px; justify-content: center; margin-top: 30px; }}
.color {{ width: 120px; height: 150px; border-radius: 12px; display: flex; flex-direction: column; align-items: center; justify-content: flex-end; padding: 10px; color: #333; font-weight: bold; box-shadow: 0 4px 8px rgba(0,0,0,0.1); }}
.hex {{ background: rgba(255,255,255,0.7); padding: 4px 8px; border-radius: 4px; }}
</style>
</head>
<body>
<h1 style="text-align:center;">🎨 Pastel Palette</h1>
<div class="palette">
"""
        for c in self.colors:
            html += f'<div class="color" style="background: {c["hex"]};">'
            html += f'<span class="hex">{c["hex"]}</span>'
            html += f'<span style="font-size:12px;">{c["hsl"]}</span>'
            html += '</div>'
        html += """
</div>
</body>
</html>
"""
        with open(filename, 'w', encoding='utf-8') as f:
            f.write(html)
        print(f"💾 Сохранено HTML: {filename}")

def main():
    parser = argparse.ArgumentParser(description='Pastel Palette Generator')
    parser.add_argument('--count', '-c', type=int, default=5, help='Количество цветов')
    parser.add_argument('--seed', '-s', type=int, default=None, help='Seed для воспроизводимости')
    parser.add_argument('--output-json', '-j', help='Сохранить в JSON')
    parser.add_argument('--output-html', '-h', help='Сохранить в HTML')
    args = parser.parse_args()

    print("🎨 Pastel Palette Generator (Python)")
    gen = PastelPalette(args.count, args.seed)
    gen.generate()
    gen.print_palette()

    if args.output_json:
        gen.save_json(args.output_json)
    else:
        gen.save_json()
    if args.output_html:
        gen.save_html(args.output_html)
    else:
        gen.save_html()

if __name__ == "__main__":
    main()

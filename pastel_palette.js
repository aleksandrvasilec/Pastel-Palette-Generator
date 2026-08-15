// pastel_palette.js — JavaScript версия

const fs = require('fs');

class PastelPalette {
    constructor(count = 5, seed = null) {
        this.count = count;
        this.seed = seed;
        if (seed !== null) {
            this._seedRandom(seed);
        }
        this.colors = [];
    }

    _seedRandom(seed) {
        this._seed = seed;
        // Простая LCG для воспроизводимости
        this._rand = () => {
            this._seed = (this._seed * 9301 + 49297) % 233280;
            return this._seed / 233280;
        };
    }

    _rand() {
        if (this._seed !== undefined) {
            return this._rand();
        }
        return Math.random();
    }

    hslToHex(h, s, l) {
        // Конвертация HSL в HEX
        s /= 100;
        l /= 100;
        const k = n => (n + h / 30) % 12;
        const a = s * Math.min(l, 1 - l);
        const f = n => l - a * Math.max(Math.min(k(n) - 3, 9 - k(n), 1), -1);
        const toHex = x => Math.round(255 * x).toString(16).padStart(2, '0');
        return `#${toHex(f(0))}${toHex(f(8))}${toHex(f(4))}`;
    }

    generate() {
        this.colors = [];
        for (let i = 0; i < this.count; i++) {
            const h = Math.floor(this._rand() * 360);
            const s = Math.floor(this._rand() * 20) + 20;   // 20-40
            const l = Math.floor(this._rand() * 20) + 75;   // 75-95
            const hex = this.hslToHex(h, s, l);
            this.colors.push({
                hsl: `hsl(${h}, ${s}%, ${l}%)`,
                hex: hex,
                rgb: `rgb(${Math.floor(h/360*255)}, ${Math.floor(s/100*255)}, ${Math.floor(l/100*255)})`
            });
        }
        return this.colors;
    }

    printPalette() {
        if (this.colors.length === 0) return;
        console.log('\n📊 Палитра из', this.colors.length, 'цветов' + (this.seed !== null ? ` (seed: ${this.seed})` : ''));
        console.log();

        // Цветные блоки
        for (const c of this.colors) {
            const hex = c.hex.replace('#', '');
            const r = parseInt(hex.substring(0, 2), 16);
            const g = parseInt(hex.substring(2, 4), 16);
            const b = parseInt(hex.substring(4, 6), 16);
            process.stdout.write(`\x1b[48;2;${r};${g};${b}m██████████\x1b[0m `);
        }
        console.log();
        for (const c of this.colors) {
            process.stdout.write(`  ${c.hex.padEnd(8)}`);
        }
        console.log();
    }

    saveJSON(filename = 'palette.json') {
        const data = { seed: this.seed, colors: this.colors };
        fs.writeFileSync(filename, JSON.stringify(data, null, 2));
        console.log(`💾 Сохранено JSON: ${filename}`);
    }

    saveHTML(filename = 'palette.html') {
        let html = `<!DOCTYPE html>
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
`;
        for (const c of this.colors) {
            html += `<div class="color" style="background: ${c.hex};">`;
            html += `<span class="hex">${c.hex}</span>`;
            html += `<span style="font-size:12px;">${c.hsl}</span>`;
            html += `</div>`;
        }
        html += `
</div>
</body>
</html>`;
        fs.writeFileSync(filename, html);
        console.log(`💾 Сохранено HTML: ${filename}`);
    }
}

function main() {
    const args = process.argv.slice(2);
    let count = 5;
    let seed = null;
    let outputJson = 'palette.json';
    let outputHtml = 'palette.html';

    for (let i = 0; i < args.length; i++) {
        if (args[i] === '--count' || args[i] === '-c') {
            count = parseInt(args[++i]) || 5;
        } else if (args[i] === '--seed' || args[i] === '-s') {
            seed = parseInt(args[++i]) || null;
        } else if (args[i] === '--output-json' || args[i] === '-j') {
            outputJson = args[++i];
        } else if (args[i] === '--output-html' || args[i] === '-h') {
            outputHtml = args[++i];
        }
    }

    console.log('🎨 Pastel Palette Generator (JavaScript)');
    const gen = new PastelPalette(count, seed);
    gen.generate();
    gen.printPalette();
    gen.saveJSON(outputJson);
    gen.saveHTML(outputHtml);
}

if (require.main === module) main();

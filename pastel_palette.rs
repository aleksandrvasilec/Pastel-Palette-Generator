// pastel_palette.rs — Rust версия

use rand::Rng;
use rand::SeedableRng;
use rand::rngs::StdRng;
use std::collections::HashMap;
use std::fs;
use std::env;

struct PastelPalette {
    count: usize,
    seed: u64,
    colors: Vec<HashMap<String, String>>,
    rng: StdRng,
}

impl PastelPalette {
    fn new(count: usize, seed: u64) -> Self {
        let rng = if seed != 0 {
            StdRng::seed_from_u64(seed)
        } else {
            StdRng::from_entropy()
        };
        PastelPalette {
            count,
            seed,
            colors: Vec::new(),
            rng,
        }
    }

    fn hsl_to_hex(h: f64, s: f64, l: f64) -> String {
        // Конвертация HSL в HEX
        let (r, g, b) = hsl_to_rgb(h, s, l);
        format!("#{:02x}{:02x}{:02x}", (r * 255.0) as u8, (g * 255.0) as u8, (b * 255.0) as u8)
    }

    fn generate(&mut self) {
        self.colors.clear();
        for _ in 0..self.count {
            let h = self.rng.gen_range(0.0..360.0);
            let s = self.rng.gen_range(20.0..40.0);
            let l = self.rng.gen_range(75.0..95.0);
            let hex = Self::hsl_to_hex(h, s, l);
            let mut color = HashMap::new();
            color.insert("hsl".to_string(), format!("hsl({:.0}, {:.0}%, {:.0}%)", h, s, l));
            color.insert("hex".to_string(), hex.clone());
            color.insert("rgb".to_string(), format!("rgb({:.0}, {:.0}, {:.0})", h/360.0*255.0, s/100.0*255.0, l/100.0*255.0));
            self.colors.push(color);
        }
    }

    fn print_palette(&self) {
        println!("\n📊 Палитра из {} цветов{}", self.colors.len(), if self.seed != 0 { format!(" (seed: {})", self.seed) } else { "".to_string() });
        println!();

        // Цветные блоки
        for c in &self.colors {
            let hex = &c["hex"][1..];
            let r = u8::from_str_radix(&hex[0..2], 16).unwrap();
            let g = u8::from_str_radix(&hex[2..4], 16).unwrap();
            let b = u8::from_str_radix(&hex[4..6], 16).unwrap();
            print!("\x1b[48;2;{};{};{}m██████████\x1b[0m ", r, g, b);
        }
        println!();
        for c in &self.colors {
            print!("  {:<8}", c["hex"]);
        }
        println!();
    }

    fn save_json(&self, filename: &str) {
        let data = serde_json::json!({
            "seed": self.seed,
            "colors": self.colors
        });
        let json = serde_json::to_string_pretty(&data).unwrap();
        fs::write(filename, json).unwrap();
        println!("💾 Сохранено JSON: {}", filename);
    }

    fn save_html(&self, filename: &str) {
        let mut html = String::from(r#"<!DOCTYPE html>
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
"#);
        for c in &self.colors {
            html.push_str(&format!(r#"<div class="color" style="background: {};">"#, c["hex"]));
            html.push_str(&format!(r#"<span class="hex">{}</span>"#, c["hex"]));
            html.push_str(&format!(r#"<span style="font-size:12px;">{}</span>"#, c["hsl"]));
            html.push_str("</div>");
        }
        html.push_str(r#"
</div>
</body>
</html>"#);
        fs::write(filename, html).unwrap();
        println!("💾 Сохранено HTML: {}", filename);
    }
}

fn hsl_to_rgb(h: f64, s: f64, l: f64) -> (f64, f64, f64) {
    let c = (1.0 - (2.0 * l - 1.0).abs()) * s;
    let x = c * (1.0 - ((h / 60.0) % 2.0 - 1.0).abs());
    let m = l - c / 2.0;
    let (r, g, b) = if h < 60.0 {
        (c, x, 0.0)
    } else if h < 120.0 {
        (x, c, 0.0)
    } else if h < 180.0 {
        (0.0, c, x)
    } else if h < 240.0 {
        (0.0, x, c)
    } else if h < 300.0 {
        (x, 0.0, c)
    } else {
        (c, 0.0, x)
    };
    (r + m, g + m, b + m)
}

fn main() {
    let args: Vec<String> = env::args().collect();
    let mut count = 5;
    let mut seed = 0;
    let mut json_file = "palette.json".to_string();
    let mut html_file = "palette.html".to_string();

    let mut i = 1;
    while i < args.len() {
        match args[i].as_str() {
            "--count" | "-c" => { count = args[i+1].parse().unwrap_or(5); i += 2; }
            "--seed" | "-s" => { seed = args[i+1].parse().unwrap_or(0); i += 2; }
            "--output-json" | "-j" => { json_file = args[i+1].clone(); i += 2; }
            "--output-html" | "-h" => { html_file = args[i+1].clone(); i += 2; }
            _ => { i += 1; }
        }
    }

    println!("🎨 Pastel Palette Generator (Rust)");
    let mut gen = PastelPalette::new(count, seed);
    gen.generate();
    gen.print_palette();
    gen.save_json(&json_file);
    gen.save_html(&html_file);
}

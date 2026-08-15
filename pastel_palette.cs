// pastel_palette.cs — C# версия

using System;
using System.Collections.Generic;
using System.Drawing;
using System.IO;
using System.Text.Json;

class PastelPalette {
    private int count;
    private int seed;
    private List<Dictionary<string, string>> colors;
    private Random rand;

    public PastelPalette(int count, int seed) {
        this.count = count;
        this.seed = seed;
        this.rand = seed != 0 ? new Random(seed) : new Random();
        this.colors = new List<Dictionary<string, string>>();
    }

    private string HslToHex(int h, int s, int l) {
        // Преобразование HSL в RGB
        float hue = h / 360f;
        float sat = s / 100f;
        float lig = l / 100f;
        Color color = FromHsl(hue, sat, lig);
        return $"#{color.R:X2}{color.G:X2}{color.B:X2}";
    }

    private Color FromHsl(float h, float s, float l) {
        // Упрощённая конвертация
        double r = 0, g = 0, b = 0;
        if (s == 0) {
            r = g = b = l;
        } else {
            var v2 = l < 0.5 ? l * (1 + s) : (l + s) - (s * l);
            var v1 = 2 * l - v2;
            r = HueToRgb(v1, v2, h + (1.0f/3.0f));
            g = HueToRgb(v1, v2, h);
            b = HueToRgb(v1, v2, h - (1.0f/3.0f));
        }
        return Color.FromArgb((int)(r * 255), (int)(g * 255), (int)(b * 255));
    }

    private double HueToRgb(double v1, double v2, double vH) {
        if (vH < 0) vH += 1;
        if (vH > 1) vH -= 1;
        if (6 * vH < 1) return v1 + (v2 - v1) * 6 * vH;
        if (2 * vH < 1) return v2;
        if (3 * vH < 2) return v1 + (v2 - v1) * (2.0/3.0 - vH) * 6;
        return v1;
    }

    public void Generate() {
        colors.Clear();
        for (int i = 0; i < count; i++) {
            int h = rand.Next(360);
            int s = rand.Next(20) + 20;
            int l = rand.Next(20) + 75;
            string hex = HslToHex(h, s, l);
            var color = new Dictionary<string, string> {
                ["hsl"] = $"hsl({h}, {s}%, {l}%)",
                ["hex"] = hex,
                ["rgb"] = $"rgb({(int)(h/360f*255)}, {(int)(s/100f*255)}, {(int)(l/100f*255)})"
            };
            colors.Add(color);
        }
    }

    public void PrintPalette() {
        Console.WriteLine($"\n📊 Палитра из {colors.Count} цветов" + (seed != 0 ? $" (seed: {seed})" : ""));
        Console.WriteLine();

        foreach (var c in colors) {
            string hex = c["hex"].Substring(1);
            int r = Convert.ToInt32(hex.Substring(0, 2), 16);
            int g = Convert.ToInt32(hex.Substring(2, 2), 16);
            int b = Convert.ToInt32(hex.Substring(4, 2), 16);
            Console.Write($"\u001B[48;2;{r};{g};{b}m██████████\u001B[0m ");
        }
        Console.WriteLine();
        foreach (var c in colors) {
            Console.Write($"  {c["hex"],-8}");
        }
        Console.WriteLine();
    }

    public void SaveJSON(string filename) {
        var data = new { seed = this.seed, colors = this.colors };
        string json = JsonSerializer.Serialize(data, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText(filename, json);
        Console.WriteLine($"💾 Сохранено JSON: {filename}");
    }

    public void SaveHTML(string filename) {
        var html = $@"<!DOCTYPE html>
<html>
<head><meta charset=""UTF-8""><title>Pastel Palette</title>
<style>
body {{ font-family: monospace; background: #f5f5f5; padding: 20px; }}
.palette {{ display: flex; gap: 20px; justify-content: center; margin-top: 30px; flex-wrap: wrap; }}
.color {{ width: 120px; height: 150px; border-radius: 12px; display: flex; flex-direction: column; align-items: center; justify-content: flex-end; padding: 10px; color: #333; font-weight: bold; box-shadow: 0 4px 8px rgba(0,0,0,0.1); }}
.hex {{ background: rgba(255,255,255,0.7); padding: 4px 8px; border-radius: 4px; }}
</style>
</head>
<body>
<h1 style=""text-align:center;"">🎨 Pastel Palette</h1>
<div class=""palette"">";
        foreach (var c in colors) {
            html += $@"<div class=""color"" style=""background: {c["hex"]};"">
    <span class=""hex"">{c["hex"]}</span>
    <span style=""font-size:12px;"">{c["hsl"]}</span>
</div>";
        }
        html += @"
</div>
</body>
</html>";
        File.WriteAllText(filename, html);
        Console.WriteLine($"💾 Сохранено HTML: {filename}");
    }

    public static void Main(string[] args) {
        int count = 5;
        int seed = 0;
        string jsonFile = "palette.json";
        string htmlFile = "palette.html";

        for (int i = 0; i < args.Length; i++) {
            if (args[i] == "--count" || args[i] == "-c") count = int.Parse(args[++i]);
            else if (args[i] == "--seed" || args[i] == "-s") seed = int.Parse(args[++i]);
            else if (args[i] == "--output-json" || args[i] == "-j") jsonFile = args[++i];
            else if (args[i] == "--output-html" || args[i] == "-h") htmlFile = args[++i];
        }

        Console.WriteLine("🎨 Pastel Palette Generator (C#)");
        var gen = new PastelPalette(count, seed);
        gen.Generate();
        gen.PrintPalette();
        gen.SaveJSON(jsonFile);
        gen.SaveHTML(htmlFile);
    }
}

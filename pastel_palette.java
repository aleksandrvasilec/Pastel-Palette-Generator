// pastel_palette.java — Java версия

import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.awt.Color;

public class pastel_palette {
    private int count;
    private long seed;
    private List<Map<String, String>> colors;
    private Random rand;

    public pastel_palette(int count, long seed) {
        this.count = count;
        this.seed = seed;
        this.rand = seed != 0 ? new Random(seed) : new Random();
        this.colors = new ArrayList<>();
    }

    private String hslToHex(int h, int s, int l) {
        float hue = h / 360.0f;
        float sat = s / 100.0f;
        float lig = l / 100.0f;
        Color color = Color.getHSBColor(hue, sat, lig);
        return String.format("#%02x%02x%02x", color.getRed(), color.getGreen(), color.getBlue());
    }

    public void generate() {
        colors.clear();
        for (int i = 0; i < count; i++) {
            int h = rand.nextInt(360);
            int s = rand.nextInt(20) + 20;
            int l = rand.nextInt(20) + 75;
            String hex = hslToHex(h, s, l);
            Map<String, String> color = new LinkedHashMap<>();
            color.put("hsl", String.format("hsl(%d, %d%%, %d%%)", h, s, l));
            color.put("hex", hex);
            color.put("rgb", String.format("rgb(%d, %d, %d)", (int)(h/360.0*255), (int)(s/100.0*255), (int)(l/100.0*255)));
            colors.add(color);
        }
    }

    public void printPalette() {
        System.out.println("\n📊 Палитра из " + colors.size() + " цветов" + (seed != 0 ? " (seed: " + seed + ")" : ""));
        System.out.println();

        // Цветные блоки
        for (Map<String, String> c : colors) {
            String hex = c.get("hex").substring(1);
            int r = Integer.parseInt(hex.substring(0, 2), 16);
            int g = Integer.parseInt(hex.substring(2, 4), 16);
            int b = Integer.parseInt(hex.substring(4, 6), 16);
            System.out.print("\u001B[48;2;" + r + ";" + g + ";" + b + "m██████████\u001B[0m ");
        }
        System.out.println();
        for (Map<String, String> c : colors) {
            System.out.printf("  %-8s", c.get("hex"));
        }
        System.out.println();
    }

    public void saveJSON(String filename) throws IOException {
        Map<String, Object> data = new LinkedHashMap<>();
        data.put("seed", seed);
        data.put("colors", colors);
        String json = new com.google.gson.GsonBuilder().setPrettyPrinting().create().toJson(data);
        Files.write(Paths.get(filename), json.getBytes());
        System.out.println("💾 Сохранено JSON: " + filename);
    }

    public void saveHTML(String filename) throws IOException {
        StringBuilder html = new StringBuilder();
        html.append("<!DOCTYPE html>\n<html>\n<head><meta charset=\"UTF-8\"><title>Pastel Palette</title>\n");
        html.append("<style>\nbody { font-family: monospace; background: #f5f5f5; padding: 20px; }\n");
        html.append(".palette { display: flex; gap: 20px; justify-content: center; margin-top: 30px; flex-wrap: wrap; }\n");
        html.append(".color { width: 120px; height: 150px; border-radius: 12px; display: flex; flex-direction: column; align-items: center; justify-content: flex-end; padding: 10px; color: #333; font-weight: bold; box-shadow: 0 4px 8px rgba(0,0,0,0.1); }\n");
        html.append(".hex { background: rgba(255,255,255,0.7); padding: 4px 8px; border-radius: 4px; }\n");
        html.append("</style>\n</head>\n<body>\n<h1 style=\"text-align:center;\">🎨 Pastel Palette</h1>\n<div class=\"palette\">\n");
        for (Map<String, String> c : colors) {
            html.append("<div class=\"color\" style=\"background: ").append(c.get("hex")).append(";\">");
            html.append("<span class=\"hex\">").append(c.get("hex")).append("</span>");
            html.append("<span style=\"font-size:12px;\">").append(c.get("hsl")).append("</span>");
            html.append("</div>\n");
        }
        html.append("</div>\n</body>\n</html>");
        Files.write(Paths.get(filename), html.toString().getBytes());
        System.out.println("💾 Сохранено HTML: " + filename);
    }

    public static void main(String[] args) throws Exception {
        int count = 5;
        long seed = 0;
        String jsonFile = "palette.json";
        String htmlFile = "palette.html";

        for (int i = 0; i < args.length; i++) {
            if (args[i].equals("--count") || args[i].equals("-c")) {
                count = Integer.parseInt(args[++i]);
            } else if (args[i].equals("--seed") || args[i].equals("-s")) {
                seed = Long.parseLong(args[++i]);
            } else if (args[i].equals("--output-json") || args[i].equals("-j")) {
                jsonFile = args[++i];
            } else if (args[i].equals("--output-html") || args[i].equals("-h")) {
                htmlFile = args[++i];
            }
        }

        System.out.println("🎨 Pastel Palette Generator (Java)");
        pastel_palette gen = new pastel_palette(count, seed);
        gen.generate();
        gen.printPalette();
        gen.saveJSON(jsonFile);
        gen.saveHTML(htmlFile);
    }
}

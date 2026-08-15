<?php
// pastel_palette.php — PHP версия

class PastelPalette {
    private $count;
    private $seed;
    private $colors = [];
    private $rng;

    public function __construct($count = 5, $seed = null) {
        $this->count = $count;
        $this->seed = $seed;
        if ($seed !== null) {
            mt_srand($seed);
        }
        $this->rng = function() {
            return mt_rand() / mt_getrandmax();
        };
    }

    private function hslToHex($h, $s, $l) {
        // Конвертация HSL в HEX
        $h = $h / 360;
        $s = $s / 100;
        $l = $l / 100;
        list($r, $g, $b) = $this->hslToRgb($h, $s, $l);
        return sprintf("#%02x%02x%02x", round($r * 255), round($g * 255), round($b * 255));
    }

    private function hslToRgb($h, $s, $l) {
        if ($s == 0) {
            return [$l, $l, $l];
        }
        $c = (1 - abs(2 * $l - 1)) * $s;
        $x = $c * (1 - abs(fmod($h / (1/6), 2) - 1));
        $m = $l - $c / 2;
        if ($h < 1/6) {
            list($r, $g, $b) = [$c, $x, 0];
        } elseif ($h < 2/6) {
            list($r, $g, $b) = [$x, $c, 0];
        } elseif ($h < 3/6) {
            list($r, $g, $b) = [0, $c, $x];
        } elseif ($h < 4/6) {
            list($r, $g, $b) = [0, $x, $c];
        } elseif ($h < 5/6) {
            list($r, $g, $b) = [$x, 0, $c];
        } else {
            list($r, $g, $b) = [$c, 0, $x];
        }
        return [$r + $m, $g + $m, $b + $m];
    }

    public function generate() {
        $this->colors = [];
        for ($i = 0; $i < $this->count; $i++) {
            $h = mt_rand(0, 359);
            $s = mt_rand(20, 40);
            $l = mt_rand(75, 95);
            $hex = $this->hslToHex($h, $s, $l);
            $this->colors[] = [
                'hsl' => "hsl($h, {$s}%, {$l}%)",
                'hex' => $hex,
                'rgb' => "rgb(" . round($h/360*255) . ", " . round($s/100*255) . ", " . round($l/100*255) . ")"
            ];
        }
        return $this->colors;
    }

    public function printPalette() {
        echo "\n📊 Палитра из " . count($this->colors) . " цветов" . ($this->seed !== null ? " (seed: {$this->seed})" : "") . "\n\n";

        // Цветные блоки
        foreach ($this->colors as $c) {
            $hex = substr($c['hex'], 1);
            $r = hexdec(substr($hex, 0, 2));
            $g = hexdec(substr($hex, 2, 2));
            $b = hexdec(substr($hex, 4, 2));
            echo "\033[48;2;{$r};{$g};{$b}m██████████\033[0m ";
        }
        echo "\n";
        foreach ($this->colors as $c) {
            printf("  %-8s", $c['hex']);
        }
        echo "\n";
    }

    public function saveJSON($filename = 'palette.json') {
        $data = ['seed' => $this->seed, 'colors' => $this->colors];
        file_put_contents($filename, json_encode($data, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
        echo "💾 Сохранено JSON: $filename\n";
    }

    public function saveHTML($filename = 'palette.html') {
        $html = '<!DOCTYPE html>
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
';
        foreach ($this->colors as $c) {
            $html .= "<div class=\"color\" style=\"background: {$c['hex']};\">";
            $html .= "<span class=\"hex\">{$c['hex']}</span>";
            $html .= "<span style=\"font-size:12px;\">{$c['hsl']}</span>";
            $html .= "</div>";
        }
        $html .= '
</div>
</body>
</html>';
        file_put_contents($filename, $html);
        echo "💾 Сохранено HTML: $filename\n";
    }
}

function main($argv) {
    $count = 5;
    $seed = null;
    $jsonFile = 'palette.json';
    $htmlFile = 'palette.html';

    for ($i = 1; $i < count($argv); $i++) {
        if ($argv[$i] == '--count' || $argv[$i] == '-c') {
            $count = (int)$argv[++$i];
        } elseif ($argv[$i] == '--seed' || $argv[$i] == '-s') {
            $seed = (int)$argv[++$i];
        } elseif ($argv[$i] == '--output-json' || $argv[$i] == '-j') {
            $jsonFile = $argv[++$i];
        } elseif ($argv[$i] == '--output-html' || $argv[$i] == '-h') {
            $htmlFile = $argv[++$i];
        }
    }

    echo "🎨 Pastel Palette Generator (PHP)\n";
    $gen = new PastelPalette($count, $seed);
    $gen->generate();
    $gen->printPalette();
    $gen->saveJSON($jsonFile);
    $gen->saveHTML($htmlFile);
}

$argc = $_SERVER['argc'] ?? 0;
$argv = $_SERVER['argv'] ?? [];
main($argv);
?>

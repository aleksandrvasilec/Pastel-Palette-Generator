# pastel_palette.rb — Ruby версия

require 'json'
require 'optparse'

class PastelPalette
  def initialize(count: 5, seed: nil)
    @count = count
    @seed = seed
    @rng = seed ? Random.new(seed) : Random.new
    @colors = []
  end

  def hsl_to_hex(h, s, l)
    # Конвертация HSL в HEX
    h = h / 360.0
    s = s / 100.0
    l = l / 100.0
    r, g, b = hsl_to_rgb(h, s, l)
    "##{ (r * 255).round.to_s(16).rjust(2, '0') }#{ (g * 255).round.to_s(16).rjust(2, '0') }#{ (b * 255).round.to_s(16).rjust(2, '0') }"
  end

  def hsl_to_rgb(h, s, l)
    if s == 0
      return l, l, l
    end
    c = (1 - (2 * l - 1).abs) * s
    x = c * (1 - ((h / (1.0/6.0)) % 2 - 1).abs)
    m = l - c / 2
    if h < 1.0/6.0
      r, g, b = c, x, 0
    elsif h < 2.0/6.0
      r, g, b = x, c, 0
    elsif h < 3.0/6.0
      r, g, b = 0, c, x
    elsif h < 4.0/6.0
      r, g, b = 0, x, c
    elsif h < 5.0/6.0
      r, g, b = x, 0, c
    else
      r, g, b = c, 0, x
    end
    [r + m, g + m, b + m]
  end

  def generate
    @colors = []
    @count.times do
      h = @rng.rand(360)
      s = @rng.rand(20) + 20
      l = @rng.rand(20) + 75
      hex = hsl_to_hex(h, s, l)
      @colors << {
        'hsl' => "hsl(#{h}, #{s}%, #{l}%)",
        'hex' => hex,
        'rgb' => "rgb(#{(h/360.0*255).round}, #{(s/100.0*255).round}, #{(l/100.0*255).round})"
      }
    end
    @colors
  end

  def print_palette
    puts "\n📊 Палитра из #{@colors.size} цветов" + (@seed ? " (seed: #{@seed})" : "")
    puts

    # Цветные блоки
    @colors.each do |c|
      hex = c['hex'][1..-1]
      r, g, b = [hex[0..1], hex[2..3], hex[4..5]].map { |x| x.to_i(16) }
      print "\e[48;2;#{r};#{g};#{b}m██████████\e[0m "
    end
    puts
    @colors.each { |c| print "  #{c['hex'].ljust(8)}" }
    puts
  end

  def save_json(filename = 'palette.json')
    data = { seed: @seed, colors: @colors }
    File.write(filename, JSON.pretty_generate(data))
    puts "💾 Сохранено JSON: #{filename}"
  end

  def save_html(filename = 'palette.html')
    html = <<~HTML
      <!DOCTYPE html>
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
    HTML
    @colors.each do |c|
      html << "<div class=\"color\" style=\"background: #{c['hex']};\">"
      html << "<span class=\"hex\">#{c['hex']}</span>"
      html << "<span style=\"font-size:12px;\">#{c['hsl']}</span>"
      html << "</div>"
    end
    html << <<~HTML
      </div>
      </body>
      </html>
    HTML
    File.write(filename, html)
    puts "💾 Сохранено HTML: #{filename}"
  end
end

def main
  options = { count: 5, seed: nil }
  parser = OptionParser.new do |opts|
    opts.banner = "Usage: ruby pastel_palette.rb [options]"
    opts.on("--count N", "-c", Integer, "Количество цветов") { |v| options[:count] = v }
    opts.on("--seed N", "-s", Integer, "Seed") { |v| options[:seed] = v }
    opts.on("--output-json FILE", "-j", "Сохранить в JSON") { |v| options[:json] = v }
    opts.on("--output-html FILE", "-h", "Сохранить в HTML") { |v| options[:html] = v }
  end.parse!

  puts "🎨 Pastel Palette Generator (Ruby)"
  gen = PastelPalette.new(count: options[:count], seed: options[:seed])
  gen.generate
  gen.print_palette
  gen.save_json(options[:json] || 'palette.json')
  gen.save_html(options[:html] || 'palette.html')
end

main if __FILE__ == $0

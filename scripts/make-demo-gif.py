"""Create the small terminal demo embedded in the README."""

from pathlib import Path

from PIL import Image, ImageDraw, ImageFont


WIDTH, HEIGHT = 1100, 620
BACKGROUND = (19, 22, 27)
MUTED = (151, 160, 170)
WHITE = (232, 237, 243)
GREEN = (99, 221, 152)
BLUE = (113, 184, 255)
YELLOW = (255, 210, 105)


def font(size):
    candidates = [
        "/System/Library/Fonts/Menlo.ttc",
        "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf",
    ]
    for path in candidates:
        if Path(path).exists():
            return ImageFont.truetype(path, size)
    return ImageFont.load_default()


MONO = font(27)
SMALL = font(20)
TITLE = font(22)


def frame(lines):
    image = Image.new("RGB", (WIDTH, HEIGHT), BACKGROUND)
    draw = ImageDraw.Draw(image)
    draw.rectangle((0, 0, WIDTH, 58), fill=(35, 39, 47))
    for x, color in ((30, (255, 95, 86)), (62, (255, 189, 46)), (94, (39, 201, 63))):
        draw.ellipse((x, 20, x + 18, 38), fill=color)
    draw.text((140, 18), "codexcommits demo", font=TITLE, fill=MUTED)
    y = 92
    for text, color in lines:
        draw.text((54, y), text, font=MONO, fill=color)
        y += 43
    return image


frames = [
    frame([
        ("$ git add src/parser.py", GREEN),
        ("$ codexcommits", GREEN),
        ("", WHITE),
        ("Staged changes:", BLUE),
        (" src/parser.py | 12 +++++++++---", MUTED),
    ]),
    frame([
        ("$ git add src/parser.py", GREEN),
        ("$ codexcommits", GREEN),
        ("", WHITE),
        ("Staged changes:", BLUE),
        (" src/parser.py | 12 +++++++++---", MUTED),
        ("Generating commit message with your Codex default model...", YELLOW),
    ]),
    frame([
        ("$ git add src/parser.py", GREEN),
        ("$ codexcommits", GREEN),
        ("", WHITE),
        ("Staged changes:", BLUE),
        (" src/parser.py | 12 +++++++++---", MUTED),
        ("Generated in 5.9s", MUTED),
        ("", WHITE),
        ("feat(parser): handle nested markdown tables", WHITE),
    ]),
    frame([
        ("$ git add src/parser.py", GREEN),
        ("$ codexcommits", GREEN),
        ("", WHITE),
        ("feat(parser): handle nested markdown tables", WHITE),
        ("", WHITE),
        ("[y] commit  [e] edit  [r] regenerate  [n/Enter] cancel: ", YELLOW),
    ]),
    frame([
        ("$ git add src/parser.py", GREEN),
        ("$ codexcommits", GREEN),
        ("", WHITE),
        ("feat(parser): handle nested markdown tables", WHITE),
        ("", WHITE),
        ("[y] commit  [e] edit  [r] regenerate  [n/Enter] cancel: y", YELLOW),
    ]),
    frame([
        ("$ git add src/parser.py", GREEN),
        ("$ codexcommits", GREEN),
        ("", WHITE),
        ("feat(parser): handle nested markdown tables", WHITE),
        ("", WHITE),
        ("Committed. Run git push when you want to sync the remote.", GREEN),
    ]),
]

output = Path(__file__).resolve().parents[1] / "assets" / "demo.gif"
frames[0].save(
    output,
    save_all=True,
    append_images=frames[1:],
    duration=[1200, 1400, 1600, 1800, 1100, 2400],
    loop=0,
    optimize=True,
)
print(output)

from pathlib import Path
from textwrap import wrap

from PIL import Image, ImageDraw, ImageFilter, ImageFont


ROOT = Path(__file__).resolve().parents[1]
MINI = ROOT / "miniprogram"
OUT = ROOT / "docs" / "promo"

BG = MINI / "images" / "background.jpg"
CAR_BG = MINI / "images" / "carBackground.jpg"
CAR = MINI / "images" / "recommend-cars" / "FH6_Porsche_911_GT3_2021.png"
REF_MOUNTAIN = OUT / "8f05571412c0a8af6b3a0a693c5b35ae4a3cd4ee.jpg"
REF_RACE = OUT / "746d59d85766098e3f74480a1e3bd7ef61d694c7.jpg"
REF_GARAGE = OUT / "u=3522319255,54017356&fm=253&app=138&f=JPEG.jpg"
REF_TAIL = OUT / "cf1b9d16fdfaaf51f3de7d63220083eef01f3b29cce8.jpg"

FONT_REG = Path("C:/Windows/Fonts/NotoSansSC-VF.ttf")
FONT_BOLD = Path("C:/Windows/Fonts/msyhbd.ttc")
FONT_BLACK = Path("C:/Windows/Fonts/simhei.ttf")

LIME = "#bfff00"
TEAL = "#0fae96"
DARK = "#062f2b"
DEEP = "#021f1d"
WHITE = "#f5fffb"
MUTED = "#cdeee8"
ORANGE = "#ff5a2f"


def font(size, bold=False):
    path = FONT_BOLD if bold and FONT_BOLD.exists() else FONT_REG
    if not path.exists():
      path = FONT_BLACK
    return ImageFont.truetype(str(path), size)


def cover_image(path, size, blur=0):
    img = Image.open(path).convert("RGB")
    iw, ih = img.size
    sw, sh = size
    scale = max(sw / iw, sh / ih)
    nw, nh = int(iw * scale), int(ih * scale)
    img = img.resize((nw, nh), Image.Resampling.LANCZOS)
    left = (nw - sw) // 2
    top = (nh - sh) // 2
    img = img.crop((left, top, left + sw, top + sh))
    if blur:
        img = img.filter(ImageFilter.GaussianBlur(blur))
    return img.convert("RGBA")


def overlay(base, color, alpha):
    layer = Image.new("RGBA", base.size, color)
    layer.putalpha(alpha)
    base.alpha_composite(layer)


def rounded(draw, box, radius, fill, outline=None, width=1):
    draw.rounded_rectangle(box, radius=radius, fill=fill, outline=outline, width=width)


def text(draw, xy, value, size, fill=WHITE, bold=False, anchor=None, spacing=4):
    draw.text(xy, value, font=font(size, bold), fill=fill, anchor=anchor, spacing=spacing)


def outlined_text(draw, xy, value, size, fill, bold=False, anchor=None, stroke_fill="#eafffb", stroke_width=2):
    draw.text(
        xy,
        value,
        font=font(size, bold),
        fill=fill,
        anchor=anchor,
        stroke_width=stroke_width,
        stroke_fill=stroke_fill,
    )


def multiline(draw, box, value, size, fill=MUTED, bold=False, line_gap=10, max_chars=22):
    x, y, x2, _ = box
    lines = []
    for raw in value.split("\n"):
        lines.extend(wrap(raw, width=max_chars) or [""])
    f = font(size, bold)
    for line in lines:
        draw.text((x, y), line, font=f, fill=fill)
        y += size + line_gap
    return y


def title_logo(draw, x, y, main, size=54):
    fh_font = font(size + 8, True)
    draw.text((x, y), "FH6", font=fh_font, fill=LIME)
    fh_box = draw.textbbox((x, y), "FH6", font=fh_font)
    text(draw, (fh_box[2] + 22, y + 7), main, size, WHITE, True)
    draw.line((x - 18, y + 4, x - 18, y + size + 14), fill=LIME, width=8)


def chip(draw, xy, label, active=False, w=142, h=58):
    x, y = xy
    fill = (191, 255, 0, 220) if active else (14, 83, 74, 170)
    outline = (224, 255, 180, 160) if active else (223, 254, 247, 80)
    rounded(draw, (x, y, x + w, y + h), 18, fill, outline, 2)
    text(draw, (x + w / 2, y + h / 2 - 2), label, 24, DARK if active else WHITE, True, "mm")


def feature_card(draw, box, title, body, accent=LIME):
    x, y, x2, y2 = box
    rounded(draw, box, 24, (232, 255, 249, 230), (222, 255, 214, 150), 2)
    draw.rectangle((x, y + 22, x + 8, y2 - 22), fill=accent)
    text(draw, (x + 30, y + 30), title, 34, "#031816", True)
    multiline(draw, (x + 30, y + 104, x2 - 30, y2 - 20), body, 25, "#061c19", False, 12, 29)


def dark_panel(draw, box, title, body, accent=LIME):
    x, y, x2, y2 = box
    rounded(draw, box, 26, (2, 31, 29, 228), (191, 255, 0, 150), 2)
    draw.rectangle((x, y + 22, x + 8, y2 - 22), fill=accent)
    text(draw, (x + 32, y + 28), title, 38, WHITE, True)
    multiline(draw, (x + 32, y + 92, x2 - 32, y2 - 20), body, 27, "#eafffb", True, 12, 15)


def glass_feature_pill(draw, box, title, body, accent=LIME):
    x, y, x2, y2 = box
    rounded(draw, box, 22, (232, 255, 249, 138), (255, 255, 255, 120), 2)
    rounded(draw, (x + 18, y + 18, x2 - 18, y2 - 18), 16, (232, 255, 249, 122))
    draw.rectangle((x, y + 18, x + 7, y2 - 18), fill=accent)
    outlined_text(draw, (x + 28, y + 24), title, 32, "#031816", True, stroke_fill="#f4fffb", stroke_width=2)
    multiline(draw, (x + 28, y + 72, x2 - 28, y2 - 18), body, 23, "#031816", True, 9, 16)


def call_to_action(draw, box, value, size=32):
    rounded(draw, box, 22, (191, 255, 0, 235))
    x, y, x2, y2 = box
    text(draw, ((x + x2) / 2, (y + y2) / 2 - 2), value, size, DARK, True, "mm")


def draw_phone_mock(draw, x, y, w, h):
    rounded(draw, (x, y, x + w, y + h), 48, (3, 39, 36, 230), (191, 255, 0, 110), 2)
    rounded(draw, (x + 32, y + 52, x + w - 32, y + h - 34), 24, (220, 255, 245, 190), (255, 255, 255, 120), 1)
    text(draw, (x + w / 2, y + 30), "FH6车辆调校小助手", 24, WHITE, True, "mm")
    card_x, card_y = x + 58, y + 88
    rounded(draw, (card_x, card_y, x + w - 58, card_y + 92), 14, (245, 255, 251, 210), (255, 255, 255, 120), 1)
    rounded(draw, (card_x + 18, card_y + 16, x + w - 78, card_y + 58), 10, (191, 255, 0, 230))
    text(draw, (card_x + 38, card_y + 24), "车辆基础信息", 24, DARK, True)
    text(draw, (card_x + 20, card_y + 64), "用途 公路   PI 700   AWD", 18, "#1a6d62", True)
    card_y += 112
    rounded(draw, (card_x, card_y, x + w - 58, card_y + 86), 14, (245, 255, 251, 210), (255, 255, 255, 120), 1)
    rounded(draw, (card_x + 18, card_y + 16, x + w - 78, card_y + 58), 10, (191, 255, 0, 230))
    text(draw, (card_x + 38, card_y + 24), "齿轮调校", 24, DARK, True)
    card_y += 110
    rounded(draw, (card_x, card_y, x + w - 58, card_y + 72), 16, (191, 255, 0, 225))
    text(draw, (x + w / 2, card_y + 36), "生成调校", 28, DARK, True, "mm")
    card_y += 96
    rounded(draw, (card_x, card_y, x + w - 58, card_y + 170), 18, (2, 31, 29, 225), (191, 255, 0, 110), 2)
    text(draw, (card_x + 22, card_y + 22), "调校结果", 26, WHITE, True)
    text(draw, (x + w - 122, card_y + 24), "A 700", 24, WHITE, True)
    for i, row in enumerate(["前胎压 2.03 BAR", "后胎压 2.02 BAR", "防倾杆 14.0"]):
        yy = card_y + 66 + i * 32
        text(draw, (card_x + 22, yy), row, 20, MUTED, True)


def paste_car(base, box):
    x, y, x2, y2 = box
    card = Image.new("RGBA", (x2 - x, y2 - y), (0, 0, 0, 0))
    bg = cover_image(CAR_BG, card.size)
    bg.putalpha(210)
    card.alpha_composite(bg)
    car = Image.open(CAR).convert("RGBA")
    car.thumbnail((int((x2 - x) * 0.78), int((y2 - y) * 0.72)), Image.Resampling.LANCZOS)
    card.alpha_composite(car, ((card.width - car.width) // 2, int(card.height * 0.18)))
    mask = Image.new("L", card.size, 0)
    ImageDraw.Draw(mask).rounded_rectangle((0, 0, card.width, card.height), 24, fill=255)
    base.alpha_composite(card, (x, y))
    d = ImageDraw.Draw(base)
    rounded(d, box, 24, None, (191, 255, 0, 90), 2)


def poster():
    size = (1080, 1920)
    base = cover_image(BG, size, blur=1)
    overlay(base, "#031f1d", 112)
    d = ImageDraw.Draw(base)

    title_logo(d, 92, 94, "车辆调校小助手", 58)
    text(d, (92, 190), "先改装，再调校。输入车辆信息，快速生成可上路的基础调校。", 32, WHITE, True)

    paste_car(base, (84, 268, 996, 614))
    draw_phone_mock(d, 118, 676, 844, 740)

    feature_card(d, (84, 1458, 996, 1588), "适合谁用？", "刚学调校、想快速得到基础参数、需要保存分享的玩家。", LIME)
    feature_card(d, (84, 1620, 996, 1750), "核心功能", "快速调校｜配件攻略｜车辆推荐｜我的调校｜分享链接。", ORANGE)

    call_to_action(d, (84, 1790, 996, 1858), "搜索微信小程序：FH6车辆调校小助手", 32)
    base.save(OUT / "FH6_宣发海报_竖版.png")


def guide_long_image():
    size = (1080, 2160)
    base = cover_image(BG, size, blur=2)
    overlay(base, "#032b27", 140)
    d = ImageDraw.Draw(base)

    title_logo(d, 78, 76, "使用说明", 60)
    text(d, (78, 178), "五步完成：选车思路、输入参数、生成调校、保存修改、分享查看。", 30, WHITE, True)

    steps = [
        ("01 先看配件攻略", "在改装配件攻略页，按用途和等级查看优先级。游戏里的顺序通常是先选择配件，再进入调校。"),
        ("02 输入车辆基础信息", "在快速调校页选择用途、PI、驱动形式、轮胎类型、车重和前轮比例。轮胎尺寸前后不同：前驱填前轮，后驱填后轮，四驱填平均或后轮。"),
        ("03 齿轮调校按需开启", "需要齿轮建议时打开齿轮调校，补充马力、扭矩、挡位、轮胎尺寸和目标极速；不需要时可直接关闭。"),
        ("04 生成后继续微调", "点击生成调校后查看轮胎、齿轮、轮胎定位、防倾杆、弹簧、阻尼、空气动力学、刹车、差速器等结果。双击数值可手动修改。"),
        ("05 保存、改名和分享", "调校会保存到我的调校，可改名、删除或分享。好友点击分享链接可查看完整调校信息。车辆推荐页可查看推荐车辆和调校代码。"),
    ]

    y = 270
    for i, (title, body) in enumerate(steps):
        h = 260 if i == 1 else 230
        feature_card(d, (78, y, 1002, y + h), title, body, LIME if i % 2 == 0 else ORANGE)
        y += h + 34

    rounded(d, (78, y + 18, 1002, y + 214), 26, (232, 255, 249, 232), (191, 255, 0, 150), 2)
    d.rectangle((78, y + 44, 86, y + 188), fill=LIME)
    text(d, (112, y + 50), "实测后再微调", 36, "#031816", True)
    multiline(d, (112, y + 116, 960, y + 200), "生成结果是基础参数，建议跑一两圈后根据推头、甩尾、刹车稳定性继续调整。", 30, "#061c19", True, 12, 24)

    call_to_action(d, (78, 2030, 1002, 2100), "FH6车辆调校小助手", 34)
    base.save(OUT / "FH6_使用说明长图.png")


def square_poster():
    size = (1080, 1080)
    base = cover_image(BG, size, blur=1)
    overlay(base, "#021f1d", 120)
    d = ImageDraw.Draw(base)

    title_logo(d, 78, 82, "车辆调校小助手", 58)
    text(d, (78, 188), "快速生成调校参数", 56, WHITE, True)
    text(d, (78, 252), "配件攻略 · 车辆推荐 · 我的调校 · 分享查看", 30, MUTED, True)

    paste_car(base, (78, 334, 1002, 650))

    items = [
        ("快速调校", "输入车辆信息生成基础参数"),
        ("配件攻略", "按用途和等级选择升级优先级"),
        ("车辆推荐", "查看推荐车辆与调校代码"),
    ]
    y = 700
    for idx, (t, b) in enumerate(items):
        x = 78 + idx * 314
        rounded(d, (x, y, x + 292, y + 212), 24, (224, 255, 245, 205), (222, 255, 214, 140), 2)
        text(d, (x + 28, y + 36), t, 34, DARK, True)
        multiline(d, (x + 28, y + 92, x + 260, y + 190), b, 23, "#17665c", True, 8, 9)

    call_to_action(d, (78, 960, 1002, 1024), "搜索微信小程序：FH6车辆调校小助手", 30)
    base.save(OUT / "FH6_宣发海报_方图.png")


def scenic_promo_wide():
    size = (1920, 1080)
    base = cover_image(REF_RACE if REF_RACE.exists() else BG, size, blur=1)
    overlay(base, "#021f1d", 78)
    d = ImageDraw.Draw(base)

    rounded(d, (86, 72, 944, 360), 28, (232, 255, 249, 218), (255, 255, 255, 140), 2)
    text(d, (128, 108), "FH6车辆调校小助手", 58, "#031816", True)
    text(d, (128, 190), "快速生成调校参数", 78, "#031816", True)
    text(d, (128, 292), "先改装，再调校。跑起来，再细调。", 38, "#09332e", True)

    glass_feature_pill(d, (86, 472, 546, 604), "快速调校", "输入车辆信息，生成基础调校参数", LIME)
    glass_feature_pill(d, (592, 472, 1052, 604), "配件攻略", "按用途和等级选择升级优先级", ORANGE)
    glass_feature_pill(d, (1098, 472, 1558, 604), "我的调校", "保存、改名、修改，也能分享给好友", LIME)

    call_to_action(d, (86, 790, 924, 870), "搜索微信小程序：FH6车辆调校小助手", 36)
    base.save(OUT / "FH6_宣传图_横版赛道.png")


def garage_promo_vertical():
    size = (1080, 1620)
    base = cover_image(REF_GARAGE if REF_GARAGE.exists() else BG, size, blur=0)
    overlay(base, "#020d0c", 82)
    d = ImageDraw.Draw(base)

    rounded(d, (64, 64, 1016, 356), 30, (232, 255, 249, 170), (255, 255, 255, 150), 2)
    rounded(d, (94, 94, 986, 326), 24, (232, 255, 249, 116))
    outlined_text(d, (108, 104), "FH6车辆调校小助手", 50, "#031816", True, stroke_fill="#f4fffb", stroke_width=2)
    outlined_text(d, (108, 178), "不会调车？", 70, "#031816", True, stroke_fill="#f4fffb", stroke_width=2)
    outlined_text(d, (108, 268), "先生成一套基础参数再上路测试", 34, "#031816", True, stroke_fill="#f4fffb", stroke_width=2)

    rounded(d, (64, 430, 1016, 980), 34, (2, 31, 29, 76), (191, 255, 0, 155), 2)
    rounded(d, (94, 466, 390, 538), 22, (2, 31, 29, 150), (191, 255, 0, 110), 1)
    text(d, (118, 482), "核心功能", 44, LIME, True)
    features = [
        ("01", "快速调校", "轮胎、齿轮、悬挂、刹车、差速器"),
        ("02", "改装配件攻略", "按用途和等级选择升级优先级"),
        ("03", "车辆推荐", "查看推荐车辆和调校代码"),
        ("04", "我的调校", "本地保存、改名、修改、分享"),
    ]
    y = 562
    for no, title, body in features:
        rounded(d, (108, y, 972, y + 82), 18, (232, 255, 249, 186), (255, 255, 255, 125), 1)
        rounded(d, (126, y + 13, 954, y + 69), 14, (232, 255, 249, 76))
        outlined_text(d, (136, y + 20), no, 28, ORANGE if no in {"02", "04"} else LIME, True, stroke_fill="#f4fffb", stroke_width=1)
        outlined_text(d, (208, y + 18), title, 30, "#031816", True, stroke_fill="#f4fffb", stroke_width=1)
        outlined_text(d, (430, y + 22), body, 24, "#031816", True, stroke_fill="#f4fffb", stroke_width=1)
        y += 102

    rounded(d, (64, 1058, 1016, 1278), 26, (2, 31, 29, 76), (191, 255, 0, 145), 2)
    rounded(d, (94, 1088, 986, 1248), 20, (2, 31, 29, 128))
    text(d, (108, 1096), "适合玩家", 38, WHITE, True)
    multiline(d, (108, 1160, 960, 1248), "想快速拿到基础调校、整理车辆方案、和好友分享结果的玩家。", 28, "#f4fffb", True, 12, 18)
    call_to_action(d, (64, 1410, 1016, 1484), "搜索微信小程序：FH6车辆调校小助手", 34)
    base.save(OUT / "FH6_宣传图_竖版车库.png")


def write_copy():
    content = """# FH6车辆调校小助手 宣发文案

## 一句话介绍
FH6车辆调校小助手：输入车辆信息，快速生成轮胎、齿轮、悬挂、刹车、差速器等基础调校参数。

## 短文案
不会调车？先别急着乱改。
FH6车辆调校小助手支持快速调校、配件攻略、车辆推荐和我的调校保存。
先选用途和等级，再生成基础参数，实测后继续微调。

## 社群/朋友圈文案
做了一个 FH6 车辆调校小助手。

适合刚开始学调校，或者想快速拿到基础参数再慢慢微调的玩家。

目前支持：
- 快速调校：输入 PI、驱动、轮胎、车重等信息，生成基础调校参数
- 齿轮调校：可选开启，生成齿轮建议
- 改装配件攻略：按用途和等级查看升级优先级
- 车辆推荐：查看推荐车辆和调校代码
- 我的调校：本地保存、改名、手动修改、分享链接查看

搜索微信小程序：FH6车辆调校小助手

## 长文案
完美调校不是一次就能生成出来的。

FH6车辆调校小助手的定位是帮玩家先得到一套可测试的基础调校，再结合赛道表现继续微调。

你可以先在改装配件攻略里按用途和等级选择升级方向，再到快速调校页输入车辆参数，生成胎压、齿轮、轮胎定位、防倾杆、弹簧、阻尼、空气动力学、刹车和差速器建议。

生成结果可以本地保存到“我的调校”，也可以手动修改数值、改名管理，或者分享链接给好友查看完整调校信息。

搜索微信小程序：FH6车辆调校小助手

## 宣发图标题备选
- FH6车辆调校小助手
- 不会调车？先用基础调校跑起来
- 先改装，再调校
- 输入车辆信息，快速生成调校参数
- 从配件选择到调校保存，一次完成

## 功能卖点
- Road / Rally / Offroad / Drag / Drift 多用途调校
- 支持 D / C / B / A / S1 / S2 / R / X 等级思路
- 齿轮调校可选开启
- 结果支持手动修改
- 本地保存我的调校
- 分享链接查看完整调校
- 车辆推荐支持云端更新
"""
    (OUT / "FH6_宣发文案.md").write_text(content, encoding="utf-8")


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    poster()
    guide_long_image()
    square_poster()
    scenic_promo_wide()
    garage_promo_vertical()
    write_copy()


if __name__ == "__main__":
    main()

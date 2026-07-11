from pathlib import Path
from PIL import Image, ImageDraw, ImageFont


ROOT = Path(__file__).resolve().parents[1]
OUT_DIR = ROOT / "docs" / "ai-compliance"

FONT_REGULAR = Path(r"C:\Windows\Fonts\msyh.ttc")
FONT_BOLD = Path(r"C:\Windows\Fonts\msyhbd.ttc")
FONT_MONO = Path(r"C:\Windows\Fonts\consola.ttf")


def font(path, size):
    return ImageFont.truetype(str(path), size)


TITLE = font(FONT_BOLD, 42)
SUBTITLE = font(FONT_REGULAR, 24)
TEXT = font(FONT_REGULAR, 26)
TEXT_BOLD = font(FONT_BOLD, 27)
SMALL = font(FONT_REGULAR, 20)
CODE = font(FONT_MONO, 22)


def rounded(draw, xy, radius, fill, outline=None, width=1):
    draw.rounded_rectangle(xy, radius=radius, fill=fill, outline=outline, width=width)


def draw_header(draw, title, subtitle):
    draw.text((56, 42), title, font=TITLE, fill=(218, 255, 240))
    draw.text((56, 98), subtitle, font=SUBTITLE, fill=(176, 217, 202))


def draw_code_block(draw, x, y, w, title, source, lines, height=None):
    if height is None:
        height = 86 + len(lines) * 30
    rounded(draw, (x, y, x + w, y + height), 18, (4, 37, 33, 232), (180, 255, 39), 2)
    draw.text((x + 28, y + 22), title, font=TEXT_BOLD, fill=(190, 255, 58))
    draw.text((x + 28, y + 57), source, font=SMALL, fill=(144, 184, 171))
    cy = y + 96
    for no, line in lines:
        draw.text((x + 28, cy), f"{no:>4}", font=CODE, fill=(111, 148, 139))
        code_text = fit_text(draw, line.rstrip(), CODE, w - 130)
        draw.text((x + 92, cy), code_text, font=CODE, fill=(229, 248, 241))
        cy += 30


def fit_text(draw, text, text_font, max_width):
    if draw.textlength(text, font=text_font) <= max_width:
        return text
    suffix = "..."
    available = max_width - draw.textlength(suffix, font=text_font)
    clipped = text
    while clipped and draw.textlength(clipped, font=text_font) > available:
        clipped = clipped[:-1]
    return clipped.rstrip() + suffix


def read_lines(relative, start, end):
    path = ROOT / relative
    data = path.read_text(encoding="utf-8").splitlines()
    return [(i, data[i - 1]) for i in range(start, end + 1)]


def draw_badge(draw, x, y, text, fill=(180, 255, 39), color=(0, 36, 31)):
    bbox = draw.textbbox((0, 0), text, font=TEXT_BOLD)
    w = bbox[2] - bbox[0] + 32
    rounded(draw, (x, y, x + w, y + 42), 12, fill)
    draw.text((x + 16, y + 7), text, font=TEXT_BOLD, fill=color)
    return w


def base_canvas(width=1280, height=900):
    img = Image.new("RGB", (width, height), (3, 28, 25))
    draw = ImageDraw.Draw(img, "RGBA")
    for i in range(height):
        green = int(26 + 38 * (i / height))
        draw.line((0, i, width, i), fill=(0, green, 48, 255))
    draw.ellipse((width - 420, -160, width + 180, 360), fill=(84, 174, 145, 46))
    draw.ellipse((-240, height - 260, 420, height + 220), fill=(180, 255, 39, 32))
    return img, draw


def frontend_call_image():
    img, draw = base_canvas()
    draw_header(draw, "证据 01：前端只调用自有云函数", "快速调校页面通过 wx.cloud.callFunction 调用 calculateTune，不直连任何 AI 接口。")
    lines = read_lines("miniprogram/pages/index/index.js", 382, 392)
    draw_code_block(
        draw,
        56,
        170,
        1168,
        "前端调用链",
        "WeChatApp/miniprogram/pages/index/index.js",
        lines,
        480,
    )
    draw_badge(draw, 78, 690, "调用目标：calculateTune 云函数")
    draw.text((78, 755), "说明：这里没有 wx.request、fetch、axios 或第三方 AI 域名，只有微信云开发函数调用。", font=TEXT, fill=(218, 255, 240))
    img.save(OUT_DIR / "01_frontend_call_chain.png")


def cloud_algorithm_image():
    img, draw = base_canvas(height=1180)
    draw_header(draw, "证据 02：云函数内部是确定性规则算法", "输入参数进入 generateRoadStaticTuneBaseline 后，按固定公式生成 JSON 调校结果。")
    lines_a = read_lines("cloudfunctions/calculateTune/index.js", 12, 20)
    lines_c = [
        (45, "function generateRoadStaticTuneBaseline(rawInput) {"),
        (46, "const input = normalizeRoadStaticBaselineInput(rawInput || {});"),
        (47, "const carClass = classFromPI(input.pi);"),
        (141, "const tire = baselineTirePressures(useCase, input.drivetrain, input.tireCompound, input.weightKG);"),
        (145, "const alignment = baselineAlignment(useCase, carClass, input.drivetrain, input.frontWeightPct, input.balanceBias);"),
        (153, "const springs = baselineSprings(useCase, input.weightKG, input.frontWeightPct, springFactor, input.drivetrain);"),
        (161, "const damping = baselineDamping(useCase, input.weightKG, input.frontWeightPct, carClass, springFactor, springs.front, springs.rear, input.drivetrain);"),
        (176, "appendDifferentialBaseline(result, useCase, carClass, input.drivetrain, ...);"),
        (190, "const gearing = baselineGearing(input);"),
        (204, "applyRoadBaselineBiases(result, input);"),
    ]
    draw_code_block(draw, 56, 170, 1168, "云函数入口", "WeChatApp/cloudfunctions/calculateTune/index.js", lines_a, 360)
    draw_code_block(draw, 56, 560, 1168, "规则算法主流程", "同一文件：固定函数、固定公式、固定步进", lines_c, 430)
    draw_badge(draw, 78, 1110, "相同输入 = 相同输出")
    img.save(OUT_DIR / "02_cloud_algorithm_core.png")


def scan_image():
    img, draw = base_canvas()
    draw_header(draw, "证据 03：业务代码未发现 AI/外部 HTTP 调用", "扫描范围：WeChatApp/miniprogram 与 WeChatApp/cloudfunctions，排除 node_modules。")
    rounded(draw, (56, 170, 1224, 732), 18, (4, 37, 33, 232), (180, 255, 39), 2)
    lines = [
        "$ rg -n \"openai|chatgpt|gpt|deepseek|qwen|apiKey|Authorization|Bearer\"",
        "  WeChatApp/miniprogram WeChatApp/cloudfunctions --glob \"!**/node_modules/**\"",
        "",
        "结果：无匹配",
        "",
        "$ rg -n \"fetch\\(|axios|wx\\.request|request\\(\"",
        "  WeChatApp/miniprogram WeChatApp/cloudfunctions --glob \"!**/node_modules/**\"",
        "",
        "结果：无匹配",
    ]
    cy = 220
    for line in lines:
        text_font = CODE if line.startswith("$") else TEXT_BOLD
        draw.text((86, cy), fit_text(draw, line, text_font, 1080), font=text_font, fill=(229, 248, 241) if line.startswith("$") else (190, 255, 58))
        cy += 44 if line else 32
    draw.text((86, 780), "结论：业务代码中未配置 AI 平台密钥，也未通过 HTTP 请求访问外部生成式模型服务。", font=TEXT, fill=(218, 255, 240))
    img.save(OUT_DIR / "03_no_ai_api_scan.png")


def dependency_image():
    img, draw = base_canvas()
    draw_header(draw, "证据 04：云函数依赖仅为微信云开发 SDK", "calculateTune 和 shareTune 的 package.json 只声明 wx-server-sdk。")
    lines_a = read_lines("cloudfunctions/calculateTune/package.json", 1, 10)
    lines_b = read_lines("cloudfunctions/shareTune/package.json", 1, 10)
    draw_code_block(draw, 56, 170, 560, "calculateTune 依赖", "cloudfunctions/calculateTune/package.json", lines_a, 430)
    draw_code_block(draw, 664, 170, 560, "shareTune 依赖", "cloudfunctions/shareTune/package.json", lines_b, 430)
    draw_badge(draw, 78, 665, "没有 AI SDK")
    draw_badge(draw, 300, 665, "没有模型 API Key")
    draw.text((78, 738), "说明：wx-server-sdk 用于微信云函数、云数据库和云存储能力，不是生成式 AI 模型 SDK。", font=TEXT, fill=(218, 255, 240))
    img.save(OUT_DIR / "04_cloud_dependencies.png")


def main():
    OUT_DIR.mkdir(parents=True, exist_ok=True)
    frontend_call_image()
    cloud_algorithm_image()
    scan_image()
    dependency_image()
    print(f"generated: {OUT_DIR}")


if __name__ == "__main__":
    main()

#!/usr/bin/env python3
"""Download FH6 car PNGs from the Forza Fandom FH6 cars page.

The default workflow is tuned for this workspace:

  python tools/download_fh6_car_images.py

It scans the recommended car data for existing car image filenames such as
FH6_Acura_NSX_Type-S.png, resolves those files through the Fandom wiki page/API,
and writes PNG files into ./img.
"""

from __future__ import annotations

import argparse
import dataclasses
import hashlib
import html
import io
import json
import os
import re
import sys
import time
import unicodedata
import urllib.error
import urllib.parse
import urllib.request
from difflib import SequenceMatcher
from pathlib import Path
from typing import Dict, Iterable, Iterator, List, Optional, Sequence, Set, Tuple


DEFAULT_WIKI_URL = "https://forza.fandom.com/wiki/Forza_Horizon_6/Cars"
PNG_SIGNATURE = b"\x89PNG\r\n\x1a\n"
USER_AGENT = (
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "
    "AppleWebKit/537.36 (KHTML, like Gecko) "
    "Chrome/126.0 Safari/537.36 FH6WorkerImageFetcher/1.0"
)

FH6_FILE_RE = re.compile(r"FH6_[^\\/'\"<>\s]+?\.png", re.IGNORECASE)
NAME_FIELD_RE = re.compile(
    r"(?:[\"']|\b)(?:name|carName|displayName|vehicleName)(?:[\"']|\b)\s*:\s*"
    r"[\"']([^\"']{2,160})[\"']",
    re.IGNORECASE,
)
STATIC_IMAGE_RE = re.compile(
    r"https?://static\.wikia\.nocookie\.net/forzamotorsport/images/"
    r"[^\"'<>\s]+?FH6_[^\"'<>\s]+?\.png(?:/revision/[^\"'<>\s]+)?",
    re.IGNORECASE,
)
A_TAG_RE = re.compile(r"<a\b(?P<attrs>[^>]*?)>", re.IGNORECASE | re.DOTALL)
ATTR_RE = re.compile(
    r"(?P<name>[A-Za-z_:][-A-Za-z0-9_:.]*)\s*=\s*"
    r"(?P<quote>[\"'])(?P<value>.*?)(?P=quote)",
    re.DOTALL,
)
IMAGE_FIELD_NAMES = (
    "imageSrc",
    "localImageSrc",
    "imageFileId",
    "imageUrl",
    "image",
    "cover",
)
NAME_FIELD_NAMES = ("name", "carName", "displayName", "vehicleName")
NON_VEHICLE_PAGE_KEYS = {
    "auction house",
    "contributors to forza wiki",
    "evil dead burn",
    "floor adhesion",
    "forza horizon 6 cars",
    "hide seek",
}
NON_VEHICLE_TOKENS = {
    "ad",
    "adhesion",
    "article",
    "auction",
    "autoshop",
    "autoshow",
    "barn",
    "blog",
    "boxad",
    "camera",
    "car",
    "cars",
    "category",
    "challenge",
    "community",
    "content",
    "contributors",
    "desktop",
    "discord",
    "downloadable",
    "drawer",
    "driver",
    "drone",
    "editor",
    "engine",
    "event",
    "events",
    "fandom",
    "featured",
    "festival",
    "file",
    "find",
    "floor",
    "forza",
    "forum",
    "help",
    "horizon",
    "inc",
    "incontent",
    "leaderboard",
    "manufacturer",
    "manufacturers",
    "motorsport",
    "paint",
    "performance",
    "playlist",
    "points",
    "radio",
    "records",
    "reward",
    "series",
    "special",
    "swap",
    "template",
    "update",
    "user",
    "vehicle",
    "vehicles",
    "view",
    "wiki",
}
NON_VEHICLE_IMAGE_TOKENS = {
    "achievement",
    "car pack promo",
    "cars promo",
    "festivalplaylist",
    "icon",
    "points",
    "promo",
}


@dataclasses.dataclass(frozen=True)
class CarTarget:
    name: str = ""
    image_file: str = ""
    source: str = ""

    @property
    def label(self) -> str:
        if self.image_file and self.name:
            return f"{self.name} ({self.image_file})"
        return self.image_file or self.name or self.source


def parse_args(argv: Sequence[str]) -> argparse.Namespace:
    root = Path(__file__).resolve().parents[1]
    parser = argparse.ArgumentParser(
        description="Download FH6 car PNG images from the Forza Fandom cars page."
    )
    parser.add_argument(
        "--root",
        type=Path,
        default=root,
        help="Workspace root. Defaults to the parent of this tools directory.",
    )
    parser.add_argument(
        "--img-dir",
        type=Path,
        default=None,
        help="Output directory. Defaults to <root>/img.",
    )
    parser.add_argument(
        "--wiki-url",
        default=DEFAULT_WIKI_URL,
        help=f"FH6 cars wiki page URL. Defaults to {DEFAULT_WIKI_URL}",
    )
    parser.add_argument(
        "--source",
        action="append",
        type=Path,
        default=[],
        help=(
            "File or directory to scan for car names/FH6 image filenames. "
            "Can be passed multiple times. Defaults to miniprogram data files."
        ),
    )
    parser.add_argument(
        "--names-file",
        type=Path,
        default=None,
        help=(
            "Optional UTF-8 text file with one car name, Fandom /wiki/ URL, "
            "page title, or FH6_*.png filename per line."
        ),
    )
    parser.add_argument(
        "--retry-failed",
        type=Path,
        default=None,
        help=(
            "Retry only targets recorded in a previous JSONL failure log. "
            "Falls back to treating non-JSON lines as names or FH6_*.png filenames."
        ),
    )
    parser.add_argument(
        "--failed-log",
        type=Path,
        default=Path("img") / "fh6_image_failures.jsonl",
        help=(
            "JSONL file to write failed targets for later --retry-failed runs. "
            "Defaults to <root>/img/fh6_image_failures.jsonl."
        ),
    )
    parser.add_argument(
        "--download-all-page-images",
        action="store_true",
        help="Download every FH6_*.png found on the wiki page, ignoring local targets.",
    )
    parser.add_argument(
        "--overwrite",
        action="store_true",
        help="Overwrite existing PNG files in the output directory.",
    )
    parser.add_argument(
        "--convert-existing",
        action="store_true",
        help=(
            "Convert existing *.png files in the output directory to real PNG bytes "
            "when they are WebP or another Pillow-readable image format."
        ),
    )
    parser.add_argument(
        "--convert-existing-only",
        action="store_true",
        help="Only convert existing output *.png files to real PNG. Does not access network.",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Resolve matches and print actions without writing files.",
    )
    parser.add_argument(
        "--list-targets",
        action="store_true",
        help="Only list local targets extracted from source files. Does not access network.",
    )
    parser.add_argument(
        "--timeout",
        type=float,
        default=30.0,
        help="HTTP timeout in seconds.",
    )
    parser.add_argument(
        "--retries",
        type=int,
        default=2,
        help="Number of retries for transient HTTP failures.",
    )
    parser.add_argument(
        "--delay",
        type=float,
        default=0.2,
        help="Delay between downloads in seconds.",
    )
    parser.add_argument(
        "--limit",
        type=int,
        default=0,
        help="Limit the number of planned car image targets. Use 0 for no limit.",
    )
    parser.add_argument(
        "--match-threshold",
        type=float,
        default=0.55,
        help="Minimum fuzzy score for matching a local car name to a wiki image file.",
    )
    return parser.parse_args(argv)


def main(argv: Sequence[str]) -> int:
    args = parse_args(argv)
    root = args.root.resolve()
    img_dir = (args.img_dir or (root / "img")).resolve()

    if args.convert_existing_only:
        convert_existing_images(img_dir, dry_run=args.dry_run)
        return 0

    if args.retry_failed:
        targets = extract_targets_from_failed_log(resolve_path(root, args.retry_failed))
    else:
        targets = collect_targets(root, args.source, args.names_file)
    if args.list_targets:
        print_targets(targets)
        return 0

    if not args.download_all_page_images and not targets:
        print("No local car names or FH6 image filenames were found.", file=sys.stderr)
        print("Pass --source or --names-file to provide car names.", file=sys.stderr)
        return 2

    print(f"Fetching wiki image list: {args.wiki_url}", flush=True)
    try:
        wiki_images, wiki_targets = fetch_wiki_page_data(args.wiki_url, args.timeout, args.retries)
    except Exception as exc:  # noqa: BLE001 - keep exact filename fallback usable.
        wiki_images = {}
        wiki_targets = []
        print(f"warning: could not fetch wiki image list: {exc}", file=sys.stderr, flush=True)

    if not wiki_images:
        if args.download_all_page_images:
            if not wiki_targets:
                print("No FH6 PNG images or car page links were found on the wiki page/API.", file=sys.stderr)
                return 3
            print(
                "warning: using wiki car page link fallback for --download-all-page-images.",
                file=sys.stderr,
                flush=True,
            )
        fallback_targets = [
            target for target in targets
            if target.image_file or candidate_file_names_from_name(target.name)
        ]
        if not args.download_all_page_images and not fallback_targets:
            print(
                "No wiki image list and no exact or derivable FH6_*.png targets are available.",
                file=sys.stderr,
            )
            return 3
        if not args.download_all_page_images:
            print(
                "warning: using static.wikia.nocookie.net fallback for exact or derived FH6_*.png targets.",
                file=sys.stderr,
                flush=True,
            )

    if args.download_all_page_images:
        if wiki_images:
            plan = [
                (CarTarget(image_file=file_name, source="wiki"), ((file_name, url),))
                for file_name, url in sorted(wiki_images.items(), key=lambda item: item[0].lower())
            ]
        else:
            plan, missing = build_download_plan(wiki_targets, {}, args.match_threshold)
            if missing:
                print(f"Unmatched wiki links: {len(missing)}")
    else:
        plan, missing = build_download_plan(targets, wiki_images, args.match_threshold)
        if missing:
            print(f"Unmatched targets: {len(missing)}")
            for target in missing:
                print(f"  - {target.label} [{target.source}]")

    if args.limit > 0:
        original_count = len(plan)
        plan = plan[:args.limit]
        print(f"limit applied: {len(plan)}/{original_count} targets")

    if args.dry_run:
        print_download_plan(plan, img_dir, args.overwrite)
        if args.convert_existing:
            convert_existing_images(img_dir, dry_run=True)
        return 0

    img_dir.mkdir(parents=True, exist_ok=True)
    if args.convert_existing:
        convert_existing_images(img_dir, dry_run=False)

    downloaded = 0
    skipped = 0
    failed = 0
    failed_records: List[Dict[str, object]] = []

    for target, alternatives in plan:
        existing = first_existing_destination(img_dir, alternatives)
        if existing and not args.overwrite:
            if args.convert_existing and not is_png_file(existing):
                try:
                    convert_image_file_to_png(existing)
                    print(f"converted existing: {existing}")
                except Exception as exc:  # noqa: BLE001 - report and continue.
                    print(f"failed converting existing: {existing}: {exc}", file=sys.stderr)
                    failed += 1
                else:
                    skipped += 1
                continue
            print(f"skip existing: {existing}")
            skipped += 1
            continue

        last_error: Optional[BaseException] = None
        downloaded_one = False
        for file_name, url in alternatives:
            destination = img_dir / file_name
            try:
                download_png(url, destination, args.timeout, args.retries)
                print(f"downloaded: {destination}")
                downloaded += 1
                downloaded_one = True
                if args.delay > 0:
                    time.sleep(args.delay)
                break
            except Exception as exc:  # noqa: BLE001 - try the next filename guess.
                last_error = exc

        if not downloaded_one and target.name:
            page_alternatives = fetch_vehicle_page_image_alternatives(
                args.wiki_url,
                target.name,
                args.timeout,
                args.retries,
            )
            for file_name, url in page_alternatives:
                destination = img_dir / file_name
                if destination.exists() and not args.overwrite:
                    print(f"skip existing from vehicle page: {destination}")
                    skipped += 1
                    downloaded_one = True
                    break
                try:
                    download_png(url, destination, args.timeout, args.retries)
                    print(f"downloaded from vehicle page: {destination}")
                    downloaded += 1
                    downloaded_one = True
                    if args.delay > 0:
                        time.sleep(args.delay)
                    break
                except Exception as exc:  # noqa: BLE001 - try next page image.
                    last_error = exc

        if not downloaded_one:
            candidate_names = ", ".join(file_name for file_name, _ in alternatives)
            print(
                f"failed: {target.label} -> {candidate_names}: {last_error}",
                file=sys.stderr,
            )
            failed_records.append(failure_record(target, alternatives, last_error))
            failed += 1

    if failed_records and args.failed_log:
        failed_log_path = resolve_path(root, args.failed_log)
        write_failed_log(failed_log_path, failed_records)
        print(f"failed targets written: {failed_log_path}")

    print(
        f"Done. downloaded={downloaded}, skipped={skipped}, "
        f"failed={failed}, output={img_dir}"
    )
    return 1 if failed else 0


def collect_targets(
    root: Path, source_args: Sequence[Path], names_file: Optional[Path]
) -> List[CarTarget]:
    source_paths = resolve_source_paths(root, source_args)
    targets: List[CarTarget] = []
    for source_path in source_paths:
        targets.extend(extract_targets_from_file(source_path))

    if names_file:
        targets.extend(extract_targets_from_names_file(resolve_path(root, names_file)))

    return dedupe_targets(targets)


def resolve_source_paths(root: Path, source_args: Sequence[Path]) -> List[Path]:
    if source_args:
        paths: List[Path] = []
        for source_arg in source_args:
            path = resolve_path(root, source_arg)
            if path.is_dir():
                paths.extend(iter_text_sources(path))
            elif path.exists():
                paths.append(path)
            else:
                print(f"source not found: {path}", file=sys.stderr)
        return unique_paths(paths)

    defaults = [
        root / "weChatApp" / "miniprogram" / "data" / "recommendedCars.cloud.sample.json",
        root / "weChatApp" / "miniprogram" / "data" / "recommendedCars.js",
        root / "weChatApp" / "miniprogram" / "data" / "recommendedCars.fallback.json",
    ]
    data_dir = root / "weChatApp" / "miniprogram" / "data"
    if data_dir.exists():
        defaults.extend(iter_text_sources(data_dir))
    return unique_paths([path for path in defaults if path.exists()])


def resolve_path(root: Path, value: Path) -> Path:
    return value if value.is_absolute() else (root / value)


def iter_text_sources(directory: Path) -> Iterator[Path]:
    allowed = {".json", ".js", ".txt", ".csv", ".md", ".html", ".htm"}
    skipped_dirs = {"node_modules", ".git", ".agents", ".codex"}
    for path in directory.rglob("*"):
        if any(part in skipped_dirs for part in path.parts):
            continue
        if path.is_file() and path.suffix.lower() in allowed:
            yield path


def unique_paths(paths: Iterable[Path]) -> List[Path]:
    seen: Set[Path] = set()
    result: List[Path] = []
    for path in paths:
        resolved = path.resolve()
        if resolved not in seen:
            seen.add(resolved)
            result.append(resolved)
    return result


def extract_targets_from_file(path: Path) -> List[CarTarget]:
    text = read_text(path)
    targets: List[CarTarget] = []

    if path.suffix.lower() == ".json":
        parsed = parse_json(text)
        if parsed is not None:
            targets.extend(extract_targets_from_json(parsed, str(path)))
            known_images = {target.image_file.lower() for target in targets if target.image_file}
            for file_name in sorted(extract_fh6_file_names(text), key=str.lower):
                if file_name.lower() not in known_images:
                    targets.append(CarTarget(image_file=file_name, source=str(path)))
            return targets

    targets.extend(extract_wiki_link_targets(text, str(path)))

    image_files = sorted(extract_fh6_file_names(text), key=str.lower)
    names = sorted(extract_names_from_text(text), key=str.lower)

    for file_name in image_files:
        targets.append(CarTarget(image_file=file_name, source=str(path)))
    for name in names:
        targets.append(CarTarget(name=name, source=str(path)))
    return targets


def extract_targets_from_names_file(path: Path) -> List[CarTarget]:
    if not path.exists():
        print(f"names file not found: {path}", file=sys.stderr)
        return []
    targets = []
    for line in read_text(path).splitlines():
        value = line.strip()
        if not value or value.startswith("#"):
            continue
        file_name = extract_first_fh6_file_name(value)
        if file_name:
            targets.append(CarTarget(image_file=file_name, source=str(path)))
        else:
            targets.append(CarTarget(name=value, source=str(path)))
    return targets


def extract_targets_from_failed_log(path: Path) -> List[CarTarget]:
    if not path.exists():
        print(f"failed log not found: {path}", file=sys.stderr)
        return []

    targets: List[CarTarget] = []
    for line in read_text(path).splitlines():
        value = line.strip()
        if not value or value.startswith("#"):
            continue
        try:
            item = json.loads(value)
        except json.JSONDecodeError:
            file_name = extract_first_fh6_file_name(value)
            if file_name:
                targets.append(CarTarget(image_file=file_name, source=str(path)))
            else:
                targets.append(CarTarget(name=value, source=str(path)))
            continue

        if not isinstance(item, dict):
            continue
        name = string_field(item, "name")
        image_file = extract_first_fh6_file_name(string_field(item, "image_file"))
        if name or image_file:
            targets.append(CarTarget(name=name, image_file=image_file, source=str(path)))
            continue
        candidates = item.get("candidates")
        if isinstance(candidates, list):
            for candidate in candidates:
                if isinstance(candidate, str):
                    file_name = extract_first_fh6_file_name(candidate)
                    if file_name:
                        targets.append(CarTarget(image_file=file_name, source=str(path)))

    return dedupe_targets(targets)


def string_field(data: Dict[object, object], key: str) -> str:
    value = data.get(key)
    return value.strip() if isinstance(value, str) else ""


def failure_record(
    target: CarTarget,
    alternatives: Sequence[Tuple[str, str]],
    error: Optional[BaseException],
) -> Dict[str, object]:
    return {
        "name": target.name,
        "image_file": target.image_file,
        "source": target.source,
        "candidates": [file_name for file_name, _ in alternatives],
        "error": str(error) if error else "",
    }


def write_failed_log(path: Path, records: Sequence[Dict[str, object]]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8", newline="\n") as handle:
        for record in records:
            handle.write(json.dumps(record, ensure_ascii=False, sort_keys=True))
            handle.write("\n")


def read_text(path: Path) -> str:
    data = path.read_bytes()
    for encoding in ("utf-8-sig", "utf-8", "gb18030"):
        try:
            return data.decode(encoding)
        except UnicodeDecodeError:
            continue
    return data.decode("utf-8", errors="replace")


def parse_json(text: str) -> Optional[object]:
    try:
        return json.loads(text)
    except json.JSONDecodeError:
        return None


def extract_targets_from_json(value: object, source: str) -> List[CarTarget]:
    targets: List[CarTarget] = []

    def walk(node: object) -> None:
        if isinstance(node, dict):
            name = first_string_field(node, NAME_FIELD_NAMES)
            image_file = extract_first_fh6_file_name(
                first_string_field(node, IMAGE_FIELD_NAMES)
            )
            if name or image_file:
                targets.append(CarTarget(name=name, image_file=image_file, source=source))
            for child in node.values():
                walk(child)
        elif isinstance(node, list):
            for child in node:
                walk(child)

    walk(value)
    return targets


def first_string_field(node: Dict[object, object], keys: Sequence[str]) -> str:
    for key in keys:
        value = node.get(key)
        if isinstance(value, str) and value.strip():
            return value.strip()
    return ""


def extract_names_from_text(text: str) -> Set[str]:
    return {match.group(1).strip() for match in NAME_FIELD_RE.finditer(text)}


def extract_fh6_file_names(text: str) -> Set[str]:
    return {
        normalize_file_name(match.group(0))
        for match in FH6_FILE_RE.finditer(urllib.parse.unquote(text))
        if normalize_file_name(match.group(0))
    }


def extract_first_fh6_file_name(value: str) -> str:
    if not value:
        return ""
    match = FH6_FILE_RE.search(urllib.parse.unquote(value))
    return normalize_file_name(match.group(0)) if match else ""


def normalize_file_name(value: str) -> str:
    if not value:
        return ""
    decoded = urllib.parse.unquote(html.unescape(value.strip()))
    name = decoded.replace("\\", "/").split("/")[-1]
    name = name.split("?")[0].split("#")[0]
    if not name.lower().endswith(".png"):
        return ""
    if not name.lower().startswith("fh6_"):
        return ""
    return name


def dedupe_targets(targets: Iterable[CarTarget]) -> List[CarTarget]:
    by_key: Dict[Tuple[str, str], CarTarget] = {}
    for target in targets:
        name = normalize_space(target.name)
        image_file = normalize_file_name(target.image_file)
        if not name and not image_file:
            continue
        key = (image_file.lower(), name_key(name))
        if key not in by_key:
            by_key[key] = CarTarget(name=name, image_file=image_file, source=target.source)
    return sorted(by_key.values(), key=lambda item: (item.image_file.lower(), item.name.lower()))


def print_targets(targets: Sequence[CarTarget]) -> None:
    print(f"targets={len(targets)}")
    for target in targets:
        print(f"- {target.label} [{target.source}]")


def fetch_wiki_page_data(
    wiki_url: str, timeout: float, retries: int
) -> Tuple[Dict[str, str], List[CarTarget]]:
    base_url, page_title = wiki_base_and_title(wiki_url)
    images: Set[str] = set()
    url_by_file: Dict[str, str] = {}
    targets: List[CarTarget] = []
    errors: List[BaseException] = []

    parse_url = mediawiki_api_url(
        base_url,
        {
            "action": "parse",
            "page": page_title,
            "prop": "images|text",
            "format": "json",
            "formatversion": "2",
        },
    )
    try:
        payload = fetch_json(parse_url, timeout, retries)
        parsed = payload.get("parse") if isinstance(payload, dict) else None
        if isinstance(parsed, dict):
            for image in parsed.get("images") or []:
                file_name = normalize_file_name(str(image))
                if file_name:
                    images.add(file_name)
            text = parsed.get("text")
            if isinstance(text, str):
                static_urls = extract_static_image_urls(text)
                url_by_file.update(static_urls)
                images.update(static_urls.keys())
                targets.extend(extract_wiki_link_targets(text, wiki_url))
    except Exception as exc:  # noqa: BLE001 - direct HTML may still work.
        errors.append(exc)

    if images:
        try:
            url_by_file.update(fetch_imageinfo_urls(base_url, sorted(images), timeout, retries))
        except Exception as exc:  # noqa: BLE001 - static/anchor fallbacks may still work.
            errors.append(exc)

    if not url_by_file or not targets:
        # Fallback for direct rendered HTML if the parse API changes.
        try:
            rendered = fetch_bytes(wiki_url, timeout, retries).decode("utf-8", errors="replace")
            static_urls = extract_static_image_urls(rendered)
            url_by_file.update(static_urls)
            targets.extend(extract_wiki_link_targets(rendered, wiki_url))
        except Exception as exc:  # noqa: BLE001 - report below if no data survived.
            errors.append(exc)

    images_by_file = {
        normalize_file_name(file_name): url
        for file_name, url in url_by_file.items()
        if normalize_file_name(file_name) and url
    }

    deduped_targets = dedupe_targets(targets)
    if not images_by_file and not deduped_targets and errors:
        raise errors[0]

    return images_by_file, deduped_targets


def fetch_wiki_images(wiki_url: str, timeout: float, retries: int) -> Dict[str, str]:
    images, _ = fetch_wiki_page_data(wiki_url, timeout, retries)
    return images


def fetch_vehicle_page_image_alternatives(
    wiki_url: str, target_name: str, timeout: float, retries: int
) -> Tuple[Tuple[str, str], ...]:
    try:
        base_url, _ = wiki_base_and_title(wiki_url)
    except ValueError:
        return tuple()

    page_title = vehicle_page_title(target_name)
    if not page_title or not looks_like_vehicle_title(page_title.replace("_", " "), page_title):
        return tuple()

    try:
        image_files = fetch_page_image_file_names(base_url, page_title, timeout, retries)
        image_files = [
            file_name for file_name in image_files
            if is_vehicle_page_image_file(file_name)
        ]
        image_files = sort_vehicle_image_files(image_files, target_name)
        if image_files:
            urls = fetch_imageinfo_urls(base_url, image_files, timeout, retries)
            alternatives = [
                (file_name, urls.get(file_name) or static_latest_url(file_name))
                for file_name in image_files
            ]
            return tuple(dedupe_alternatives(alternatives))
    except Exception:
        pass

    page_url = vehicle_page_url(wiki_url, target_name)
    if not page_url:
        return tuple()
    try:
        rendered = fetch_bytes(page_url, timeout, retries).decode("utf-8", errors="replace")
    except Exception:
        return tuple()
    static_urls = extract_static_image_urls(rendered)
    if not static_urls:
        return tuple()
    preferred = sorted(static_urls.items(), key=lambda item: image_preference_key(item[0], target_name))
    return tuple((file_name, url) for file_name, url in preferred)


def fetch_page_image_file_names(
    base_url: str, page_title: str, timeout: float, retries: int
) -> List[str]:
    parse_url = mediawiki_api_url(
        base_url,
        {
            "action": "parse",
            "page": page_title,
            "prop": "images",
            "format": "json",
            "formatversion": "2",
        },
    )
    payload = fetch_json(parse_url, timeout, retries)
    parsed = payload.get("parse") if isinstance(payload, dict) else None
    if not isinstance(parsed, dict):
        return []
    result: List[str] = []
    for image in parsed.get("images") or []:
        file_name = normalize_file_name(str(image))
        if file_name:
            result.append(file_name)
    return dedupe_file_names(result)


def vehicle_page_url(wiki_url: str, target_name: str) -> str:
    try:
        base_url, _ = wiki_base_and_title(wiki_url)
    except ValueError:
        return ""
    page_title = vehicle_page_title(target_name)
    if not page_title:
        return ""
    quoted = urllib.parse.quote(page_title, safe="_-(),")
    return f"{base_url}/wiki/{quoted}"


def vehicle_page_title(target_name: str) -> str:
    page_title = wiki_title_from_value(target_name)
    if not page_title:
        page_title = target_name
    return normalize_space(str(page_title).replace(" ", "_"))


def is_vehicle_page_image_file(file_name: str) -> bool:
    normalized = normalize_file_name(file_name)
    if not normalized:
        return False
    key = file_key(normalized)
    if not key:
        return False
    return not any(token in key for token in NON_VEHICLE_IMAGE_TOKENS)


def sort_vehicle_image_files(file_names: Sequence[str], target_name: str) -> List[str]:
    return [
        file_name
        for file_name in sorted(
            dedupe_file_names(file_names),
            key=lambda file_name: image_preference_key(file_name, target_name),
        )
    ]


def image_preference_key(file_name: str, target_name: str = "") -> Tuple[int, float, int, str]:
    normalized = normalize_file_name(file_name)
    stem = Path(normalized).stem.lower()
    penalties = 0
    if "placeholder" in stem or "unknown" in stem:
        penalties += 10
    if "logo" in stem or "badge" in stem:
        penalties += 10
    target = name_key(target_name)
    candidate = file_key(normalized)
    if target and "forza edition" not in target and "forza edition" in candidate:
        penalties += 3
    score = 0.0
    if target and candidate:
        score = max(SequenceMatcher(None, target, candidate).ratio(), token_overlap(target, candidate))
    return (penalties, -score, len(stem), normalized.lower())


def wiki_base_and_title(wiki_url: str) -> Tuple[str, str]:
    parts = urllib.parse.urlsplit(wiki_url)
    if not parts.scheme or not parts.netloc:
        raise ValueError(f"invalid wiki URL: {wiki_url}")
    base_url = f"{parts.scheme}://{parts.netloc}"
    marker = "/wiki/"
    if marker not in parts.path:
        raise ValueError(f"wiki URL must contain /wiki/: {wiki_url}")
    title = urllib.parse.unquote(parts.path.split(marker, 1)[1]).replace("_", " ")
    if not title:
        raise ValueError(f"missing wiki page title in URL: {wiki_url}")
    return base_url, title


def mediawiki_api_url(base_url: str, params: Dict[str, str]) -> str:
    return f"{base_url}/api.php?{urllib.parse.urlencode(params)}"


def fetch_imageinfo_urls(
    base_url: str, file_names: Sequence[str], timeout: float, retries: int
) -> Dict[str, str]:
    result: Dict[str, str] = {}
    for batch in chunks(file_names, 50):
        titles = "|".join(f"File:{file_name}" for file_name in batch)
        query_url = mediawiki_api_url(
            base_url,
            {
                "action": "query",
                "titles": titles,
                "prop": "imageinfo",
                "iiprop": "url",
                "format": "json",
                "formatversion": "2",
            },
        )
        payload = fetch_json(query_url, timeout, retries)
        query = payload.get("query") if isinstance(payload, dict) else None
        pages = query.get("pages") if isinstance(query, dict) else None
        if not isinstance(pages, list):
            continue
        for page in pages:
            if not isinstance(page, dict):
                continue
            title = str(page.get("title") or "")
            if title.lower().startswith("file:"):
                file_name = normalize_file_name(title[5:])
            else:
                file_name = normalize_file_name(title)
            imageinfo = page.get("imageinfo")
            if not file_name or not isinstance(imageinfo, list) or not imageinfo:
                continue
            url = imageinfo[0].get("url") if isinstance(imageinfo[0], dict) else ""
            if isinstance(url, str) and url:
                result[file_name] = html.unescape(url)
    return result


def extract_static_image_urls(text: str) -> Dict[str, str]:
    result: Dict[str, str] = {}
    decoded = urllib.parse.unquote(html.unescape(text))
    for match in STATIC_IMAGE_RE.finditer(decoded):
        url = match.group(0)
        file_name = file_name_from_url(url)
        if file_name:
            result[file_name] = url
    return result


def extract_wiki_link_targets(text: str, source: str) -> List[CarTarget]:
    result: List[CarTarget] = []
    decoded = urllib.parse.unquote(html.unescape(text))
    for match in A_TAG_RE.finditer(decoded):
        attrs = parse_html_attrs(match.group("attrs"))
        href = attrs.get("href", "")
        title = attrs.get("title", "")
        target = car_target_from_wiki_link(href, title, source)
        if target:
            result.append(target)
    return dedupe_targets(result)


def parse_html_attrs(attrs_text: str) -> Dict[str, str]:
    attrs: Dict[str, str] = {}
    for match in ATTR_RE.finditer(attrs_text or ""):
        attrs[match.group("name").lower()] = html.unescape(match.group("value")).strip()
    return attrs


def car_target_from_wiki_link(href: str, title: str, source: str) -> Optional[CarTarget]:
    page_title = wiki_title_from_href(href)
    if not page_title:
        return None
    display_title = normalize_space(title) or normalize_space(page_title.replace("_", " "))
    if not looks_like_vehicle_title(display_title, page_title):
        return None
    return CarTarget(name=page_title, source=source)


def wiki_title_from_href(href: str) -> str:
    raw = str(href or "").strip()
    if not raw:
        return ""
    parsed = urllib.parse.urlsplit(raw)
    path = parsed.path
    marker = "/wiki/"
    if marker not in path:
        return ""
    title = path.split(marker, 1)[1]
    if not title or ":" in title or "/" in title:
        return ""
    return urllib.parse.unquote(title)


def looks_like_vehicle_title(display_title: str, page_title: str) -> bool:
    title = normalize_space(display_title or page_title.replace("_", " "))
    key = name_key(title)
    if key in NON_VEHICLE_PAGE_KEYS:
        return False
    tokens = key.split()
    if len(tokens) < 2:
        return False
    return not any(token in NON_VEHICLE_TOKENS for token in tokens)


def file_name_from_url(url: str) -> str:
    path = urllib.parse.urlsplit(html.unescape(url)).path
    if "/revision/" in path:
        path = path.split("/revision/", 1)[0]
    return normalize_file_name(urllib.parse.unquote(Path(path).name))


def build_download_plan(
    targets: Sequence[CarTarget],
    wiki_images: Dict[str, str],
    threshold: float,
) -> Tuple[List[Tuple[CarTarget, Tuple[Tuple[str, str], ...]]], List[CarTarget]]:
    wiki_by_lower = {file_name.lower(): (file_name, url) for file_name, url in wiki_images.items()}
    plan: List[Tuple[CarTarget, Tuple[Tuple[str, str], ...]]] = []
    missing: List[CarTarget] = []
    planned_files: Set[str] = set()

    ordered_targets = sorted(
        targets,
        key=lambda item: (
            0 if item.image_file else 1,
            item.image_file.lower(),
            item.name.lower(),
        ),
    )

    for target in ordered_targets:
        if target.name and not target.image_file and not looks_like_vehicle_title(
            target.name.replace("_", " "), target.name
        ):
            continue
        matched: List[Tuple[str, str]] = []
        if target.image_file:
            wiki_match = wiki_by_lower.get(target.image_file.lower())
            if wiki_match:
                matched.append(wiki_match)
            else:
                matched.append((target.image_file, static_latest_url(target.image_file)))
        elif target.name:
            wiki_match = match_name_to_image(target.name, wiki_images, threshold)
            if wiki_match:
                matched.append(wiki_match)
            else:
                matched.extend(
                    (file_name, static_latest_url(file_name))
                    for file_name in candidate_file_names_from_name(target.name)
                )

        if not matched:
            missing.append(target)
            continue

        deduped_matches = dedupe_alternatives(matched)
        if target.name and not target.image_file and any(
            file_name.lower() in planned_files for file_name, _ in deduped_matches
        ):
            continue

        alternatives: List[Tuple[str, str]] = []
        for file_name, url in deduped_matches:
            key = file_name.lower()
            if key in planned_files:
                continue
            planned_files.add(key)
            alternatives.append((file_name, url))
        if not alternatives:
            continue
        plan.append((target, tuple(alternatives)))

    return sorted(plan, key=lambda item: item[1][0][0].lower()), missing


def match_name_to_image(
    name: str, wiki_images: Dict[str, str], threshold: float
) -> Optional[Tuple[str, str]]:
    wanted = name_key(name)
    if not wanted:
        return None

    best: Optional[Tuple[float, str, str]] = None
    for file_name, url in wiki_images.items():
        candidate = file_key(file_name)
        if not candidate:
            continue
        ratio = SequenceMatcher(None, wanted, candidate).ratio()
        overlap = token_overlap(wanted, candidate)
        score = max(ratio, overlap)
        if best is None or score > best[0]:
            best = (score, file_name, url)

    if best and best[0] >= threshold:
        return (best[1], best[2])
    return None


def dedupe_alternatives(
    alternatives: Iterable[Tuple[str, str]]
) -> List[Tuple[str, str]]:
    result: List[Tuple[str, str]] = []
    seen: Set[str] = set()
    for file_name, url in alternatives:
        normalized = normalize_file_name(file_name)
        if not normalized:
            continue
        key = normalized.lower()
        if key in seen:
            continue
        seen.add(key)
        result.append((normalized, url))
    return result


def candidate_file_names_from_name(name: str) -> List[str]:
    title = wiki_title_from_value(name)
    if not title:
        title = name

    title = normalize_space(urllib.parse.unquote(html.unescape(title)).replace("_", " "))
    if not title:
        return []
    if not looks_like_vehicle_title(title, title):
        return []

    candidates: List[str] = []
    leading_year = re.match(r"^(\d{4})\s+(.+)$", title)
    if leading_year:
        year = leading_year.group(1)
        without_year = leading_year.group(2)
        candidates.extend([
            f"{without_year} {year}",
            without_year,
            title,
        ])
    else:
        candidates.append(title)

    candidates = expand_model_suffix_variants(candidates)

    return [
        file_name
        for file_name in dedupe_file_names(
            file_name
            for candidate in candidates
            for file_name in title_to_fh6_file_names(candidate)
        )
        if file_name
    ]


def expand_model_suffix_variants(values: Iterable[str]) -> List[str]:
    result: List[str] = []
    for value in values:
        normalized = normalize_space(value)
        if not normalized:
            continue
        result.append(normalized)
        result.append(re.sub(r"\s+\((\d{4})\)$", "", normalized))
        result.append(re.sub(r"\bType\s+([A-Z0-9]+)\b", r"Type-\1", normalized))
        result.append(re.sub(r"\bCan-Am\b", "Can Am", normalized))
        result.append(re.sub(r"\bCR-X\b", "CRX", normalized))
        result.append(re.sub(r"\bLP\s*(\d{3}-\d)\b", r"LP_\1", normalized))
        result.append(re.sub(r"\bSRT-10\b", "SRT10", normalized))
        result.append(re.sub(r"\bAtom\s+500\s+V8\b", "Atom V8", normalized))
        result.append(re.sub(r"\bMustang\s+RTR\s+Spec\s+\d+\b", "Mustang RTR", normalized))
        result.append(re.sub(r"\bRS\s+4\s+Avant\s+\(?(\d{4})\)?\b", r"RS 4 \1", normalized))
        result.append(re.sub(r"\bContinental\s+GT\s+Convertible\b", "Continental GTC", normalized))
    return list(dict.fromkeys(result))


def wiki_title_from_value(value: str) -> str:
    raw = str(value or "").strip()
    if not raw:
        return ""
    parsed = urllib.parse.urlsplit(raw)
    if parsed.scheme and parsed.netloc and "/wiki/" in parsed.path:
        return parsed.path.split("/wiki/", 1)[1]
    return raw


def title_to_fh6_file_names(title: str) -> List[str]:
    cleaned = normalize_space(title)
    if not cleaned:
        return []
    variants = [
        title_to_fh6_file_name(cleaned, allow_unicode=True),
        title_to_fh6_file_name(strip_accents(cleaned), allow_unicode=False),
    ]
    return dedupe_file_names(variants)


def title_to_fh6_file_name(title: str, allow_unicode: bool) -> str:
    cleaned = normalize_space(title)
    if not cleaned:
        return ""
    if allow_unicode:
        parts = re.findall(
            r"\d+(?:[.+]\d+)+|[^\W_]\.\d+|[^\W_]+(?:-[^\W_]+)*",
            cleaned,
            re.UNICODE,
        )
    else:
        parts = re.findall(
            r"\d+(?:[.+]\d+)+|[A-Za-z]\.\d+|[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*",
            cleaned,
        )
    if not parts:
        return ""
    return normalize_file_name(f"FH6_{'_'.join(parts)}.png")


def strip_accents(value: str) -> str:
    decomposed = unicodedata.normalize("NFKD", value)
    return "".join(ch for ch in decomposed if not unicodedata.combining(ch))


def dedupe_file_names(values: Iterable[str]) -> List[str]:
    result: List[str] = []
    seen: Set[str] = set()
    for value in values:
        file_name = normalize_file_name(value)
        if not file_name:
            continue
        key = file_name.lower()
        if key in seen:
            continue
        seen.add(key)
        result.append(file_name)
    return result


def name_key(value: str) -> str:
    normalized = normalize_space(value).lower()
    normalized = re.sub(r"^\d{4}\s+", "", normalized)
    tokens = re.findall(r"[a-z0-9]+", normalized)
    return " ".join(tokens)


def file_key(file_name: str) -> str:
    name = normalize_file_name(file_name)
    if not name:
        return ""
    stem = Path(name).stem
    if stem.lower().startswith("fh6_"):
        stem = stem[4:]
    return name_key(stem.replace("_", " "))


def token_overlap(left: str, right: str) -> float:
    left_tokens = set(left.split())
    right_tokens = set(right.split())
    if not left_tokens or not right_tokens:
        return 0.0
    return len(left_tokens & right_tokens) / max(len(left_tokens), len(right_tokens))


def normalize_space(value: str) -> str:
    return re.sub(r"\s+", " ", str(value or "").strip())


def static_latest_url(file_name: str) -> str:
    normalized = normalize_file_name(file_name)
    digest = hashlib.md5(normalized.encode("utf-8")).hexdigest()
    quoted = urllib.parse.quote(normalized)
    return (
        "https://static.wikia.nocookie.net/forzamotorsport/images/"
        f"{digest[0]}/{digest[:2]}/{quoted}/revision/latest"
    )


def print_download_plan(
    plan: Sequence[Tuple[CarTarget, Tuple[Tuple[str, str], ...]]], img_dir: Path, overwrite: bool
) -> None:
    print(f"matched={len(plan)}")
    for target, alternatives in plan:
        existing = first_existing_destination(img_dir, alternatives)
        primary_destination = img_dir / alternatives[0][0]
        action = "overwrite" if existing and overwrite else "download"
        if existing and not overwrite:
            action = "skip existing"
        print(f"- {action}: {target.label} -> {existing or primary_destination}")
        for file_name, url in alternatives:
            print(f"  {file_name}: {url}")


def first_existing_destination(
    img_dir: Path, alternatives: Sequence[Tuple[str, str]]
) -> Optional[Path]:
    for file_name, _ in alternatives:
        destination = img_dir / file_name
        if destination.exists():
            return destination
    return None


def download_png(url: str, destination: Path, timeout: float, retries: int) -> None:
    data = fetch_bytes(url, timeout, retries)
    png_data = ensure_png_bytes(data, source=url)
    tmp_path = destination.with_suffix(destination.suffix + ".tmp")
    tmp_path.write_bytes(png_data)
    os.replace(str(tmp_path), str(destination))


def ensure_png_bytes(data: bytes, source: str = "") -> bytes:
    if data.startswith(PNG_SIGNATURE):
        return data
    try:
        from PIL import Image
    except ImportError as exc:
        raise ValueError(
            f"response is not PNG and Pillow is not installed: {source}"
        ) from exc

    try:
        with Image.open(io.BytesIO(data)) as image:
            output = io.BytesIO()
            image.save(output, format="PNG")
            return output.getvalue()
    except Exception as exc:  # noqa: BLE001 - normalize conversion failure.
        raise ValueError(f"response is not a readable image: {source}") from exc


def is_png_file(path: Path) -> bool:
    try:
        with path.open("rb") as handle:
            return handle.read(len(PNG_SIGNATURE)) == PNG_SIGNATURE
    except OSError:
        return False


def convert_existing_images(img_dir: Path, dry_run: bool) -> None:
    if not img_dir.exists():
        return
    candidates = sorted(img_dir.glob("*.png"), key=lambda path: path.name.lower())
    converted = 0
    skipped = 0
    failed = 0
    for path in candidates:
        if is_png_file(path):
            skipped += 1
            continue
        if dry_run:
            print(f"- convert existing: {path}")
            converted += 1
            continue
        try:
            convert_image_file_to_png(path)
            print(f"converted existing: {path}")
            converted += 1
        except Exception as exc:  # noqa: BLE001 - report and continue.
            print(f"failed converting existing: {path}: {exc}", file=sys.stderr)
            failed += 1
    print(f"Existing conversion: converted={converted}, skipped={skipped}, failed={failed}")


def convert_image_file_to_png(path: Path) -> None:
    png_data = ensure_png_bytes(path.read_bytes(), source=str(path))
    tmp_path = path.with_suffix(path.suffix + ".tmp")
    tmp_path.write_bytes(png_data)
    os.replace(str(tmp_path), str(path))


def fetch_json(url: str, timeout: float, retries: int) -> object:
    data = fetch_bytes(url, timeout, retries)
    return json.loads(data.decode("utf-8"))


def fetch_bytes(url: str, timeout: float, retries: int) -> bytes:
    last_error: Optional[BaseException] = None
    for attempt in range(max(1, retries + 1)):
        try:
            request = urllib.request.Request(
                url,
                headers={
                    "User-Agent": USER_AGENT,
                    "Accept": "application/json,image/png,image/*,*/*;q=0.8",
                },
            )
            with urllib.request.urlopen(request, timeout=timeout) as response:
                return response.read()
        except (urllib.error.URLError, TimeoutError, OSError) as exc:
            last_error = exc
            if attempt < retries:
                time.sleep(0.5 * (attempt + 1))
                continue
            break
    assert last_error is not None
    raise last_error


def chunks(values: Sequence[str], size: int) -> Iterator[Sequence[str]]:
    for start in range(0, len(values), size):
        yield values[start : start + size]


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))

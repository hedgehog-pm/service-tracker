from pathlib import Path

from mkdocs.structure.files import File


def on_files(files, config):
    root = Path(config.config_file_path).parent
    for spec in sorted((root / "api/server").rglob("openapi.yaml")):
        files.append(
            File(
                path=spec.relative_to(root).as_posix(),
                src_dir=str(root),
                dest_dir=config.site_dir,
                use_directory_urls=config.use_directory_urls,
            )
        )
    return files

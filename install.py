#! /usr/bin/env python3

import shutil
from pathlib import Path
import os

install_path = Path("~/bin").expanduser()


def run_install():
    # ensure path exist
    install_path.mkdir(exist_ok=True)

    start_path = Path(".").joinpath("build")
    with os.scandir("build") as entries:
        for entry in entries:
            tool_name = entry.name
            src = start_path.joinpath(tool_name)
            dst = install_path.joinpath(f".{tool_name}")
            # shutil.copy(src, dst)
            print(f" copy {src} to {dst}")
    print(f"finish installation")


def main():
    run_install()


if __name__ == "__main__":
    main()

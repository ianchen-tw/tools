#! /usr/bin/env python3

import subprocess as sp
import argparse
import os
import shutil
from pathlib import Path

# target programs
targets = ['update-all']

parser = argparse.ArgumentParser(description='Install/Uninstall Tools')
parser.add_argument("--remove", help="remove the installation", action="store_true")


install_path = Path('~/bin').expanduser()

def run_install():
    outcome = sp.run('make', cwd='.')
    if outcome.returncode != 0:
        class MakeNotSuccess(Exception):
            pass
        raise MakeNotSuccess

    # ensure path exist
    install_path.mkdir(exist_ok=True)

    start_path = Path ('.').joinpath('build')
    for cmd_name in targets:
        src = start_path.joinpath(cmd_name)
        dst = install_path.joinpath(f".{cmd_name}")
        shutil.copy(src, dst)
        print(f" copy {src} to {dst}")
    print(f'finish installation')

def run_remove():
    # remove program in bin folder
    for cmd_name in targets:
        target_file_path = install_path.joinpath(f'.{cmd_name}')
        if target_file_path.exists():
            print(f" remove {target_file_path}")
            os.remove(target_file_path)

    # cleanup build folder
    print(f" remove build folder")
    shutil.rmtree('./build', ignore_errors=True)

    # cleanup cargo tmp files
    for cmd_name in targets:
        sp.run(['cargo', 'clean'], cwd=f'./{cmd_name}')
    print('finish remove files')

def main():
    args = parser.parse_args()
    if args.remove:
        run_remove()
    else:
        run_install()

if __name__ == "__main__":
    main()
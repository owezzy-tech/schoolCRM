"""Prepare the mounted file directory, then run the server without root privileges."""

import os
import sys
from pathlib import Path


def run(command: list[str]) -> None:
    if not command:
        raise SystemExit("A container command is required")
    directory = Path(os.environ.get("RAG_FILE_STORAGE_DIR", "/service/var/files")).resolve()
    if not directory.is_relative_to(Path("/service/var")):
        raise SystemExit("Container file storage must be inside /service/var")
    directory.mkdir(parents=True, exist_ok=True)
    if os.getuid() == 0:
        os.chown(directory, 1000, 1000)
        os.setgroups([])
        os.setgid(1000)
        os.setuid(1000)
        os.environ["HOME"] = "/service"
    if not os.access(directory, os.W_OK):
        raise SystemExit("Container file storage is not writable")
    os.execvp(command[0], command)


if __name__ == "__main__":
    run(sys.argv[1:])

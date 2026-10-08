from unittest.mock import Mock

import pytest

from infrastructure import container_entrypoint as entrypoint


def test_root_prepares_storage_then_drops_privileges_before_exec(monkeypatch, capsys):
    operations = []
    monkeypatch.setenv("RAG_FILE_STORAGE_DIR", "/service/var/files")
    monkeypatch.setenv("HOME", "/root")
    monkeypatch.setattr(entrypoint.Path, "mkdir", Mock())
    monkeypatch.setattr(entrypoint.os, "getuid", Mock(side_effect=[0, 1000]))
    monkeypatch.setattr(entrypoint.os, "chown", lambda *args: operations.append("chown"))
    monkeypatch.setattr(entrypoint.os, "setgroups", lambda *args: operations.append("groups"))
    monkeypatch.setattr(entrypoint.os, "setgid", lambda *args: operations.append("gid"))
    monkeypatch.setattr(entrypoint.os, "setuid", lambda *args: operations.append("uid"))
    monkeypatch.setattr(entrypoint.os, "access", lambda *args: True)
    monkeypatch.setattr(entrypoint.os, "execvp", lambda *args: operations.append("exec"))
    entrypoint.run(["uvicorn", "main:app"])
    assert operations == ["chown", "groups", "gid", "uid", "exec"]
    assert entrypoint.os.environ["HOME"] == "/service"
    assert capsys.readouterr().out == "rag-container uid=1000 storage=writable\n"


def test_nonroot_requires_writable_storage_and_does_not_chown(monkeypatch):
    monkeypatch.setenv("RAG_FILE_STORAGE_DIR", "/service/var/files")
    monkeypatch.setattr(entrypoint.Path, "mkdir", Mock())
    monkeypatch.setattr(entrypoint.os, "getuid", lambda: 1000)
    chown = Mock()
    monkeypatch.setattr(entrypoint.os, "chown", chown)
    monkeypatch.setattr(entrypoint.os, "access", lambda *args: False)
    with pytest.raises(SystemExit, match="not writable"):
        entrypoint.run(["uvicorn", "main:app"])
    chown.assert_not_called()


@pytest.mark.parametrize("directory", ["/", "/etc", "/service/var/../../etc"])
def test_rejects_storage_outside_owned_directory_before_mutation(monkeypatch, directory):
    monkeypatch.setenv("RAG_FILE_STORAGE_DIR", directory)
    mkdir = Mock()
    monkeypatch.setattr(entrypoint.Path, "mkdir", mkdir)
    with pytest.raises(SystemExit, match="inside /service/var"):
        entrypoint.run(["uvicorn", "main:app"])
    mkdir.assert_not_called()


def test_missing_command_fails_before_storage_mutation(monkeypatch):
    mkdir = Mock()
    monkeypatch.setattr(entrypoint.Path, "mkdir", mkdir)
    with pytest.raises(SystemExit, match="command is required"):
        entrypoint.run([])
    mkdir.assert_not_called()

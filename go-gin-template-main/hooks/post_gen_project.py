from __future__ import print_function
import os
import shutil
from subprocess import Popen

PROJECT_DIRECTORY = os.path.realpath(os.path.curdir)

def init_git():
    """
    Initialises git on the new project folder
    """
    GIT_COMMANDS = [
        ["git", "init"],
        ["git", "add", "."],
        ["git", "commit", "-a", "-m", "Initial Commit."]
    ]

    for command in GIT_COMMANDS:
        git = Popen(command, cwd=PROJECT_DIRECTORY)
        git.wait()


def update_deps():
    """Update private deps to latest version"""

    deps = [
        "git-ffd.kz/fmobile/ferr",
        "git-ffd.kz/pkg/goerr",
        "git-ffd.kz/pkg/golog",
        "git-ffd.kz/pkg/gotags",
        "git-ffd.kz/pkg/gotransaction",
        "git-ffd.kz/pkg/gobackoff",
        "git-ffd.kz/pkg/gosentry",
        "git-ffd.kz/pkg/goauth",
        "git-ffd.kz/pkg/gowatermill",
        "git-ffd.kz/fmobile/events/goevents",
        "git-ffd.kz/pkg/clientrip",
        "git-ffd.kz/pkg/requestid",
    ]

    for dep in deps:
        gomod = Popen(["go", "get", dep], cwd=PROJECT_DIRECTORY)
        gomod.wait()

    gomodtidy = Popen(["go", "mod", "tidy"], cwd=PROJECT_DIRECTORY)
    gomodtidy.wait()


if __name__ == "__main__":
    init_git()
    update_deps()

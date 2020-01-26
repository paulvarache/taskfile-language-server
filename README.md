# Taskfile Language Server

[![License: GPL v2](https://img.shields.io/badge/License-GPL%20v2-blue.svg)](https://www.gnu.org/licenses/old-licenses/gpl-2.0.en.html)

This project is an LSP implementation for https://taskfile.dev

## Features

Here is the list of supported features of the Language Server Protocol

### Completion

The server supports compleion for expression in values

## Custom method

One custom method is supported: `extension/getTasks`. It returns a list of tasks for a given Taskfile.

Request must implement the following structure: 

```json
{
    "fsPath": "<absolute-path-to-the-taskfile>"
}
```

The response will return a list of tasks implementing the following structure:

```json
{
    "scope": "<absolute-path-to-the-taskfile>",
    "task": {
        "value": "<name-of-the-task>",
        "startLine": "<line-no-starting-the-task>",
        "startCol": "<column-no-starting-the-task>",
        "endLine": "<line-no-ending-the-task>",
        "endCol": "<column-no-ending-the-task>",
    }
}
```


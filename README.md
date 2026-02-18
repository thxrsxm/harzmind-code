# HarzMind Code

HarzMind Code is a powerful command-line interface tool that enables you to have a conversation with a Large Language Model that has full context of your entire project's codebase.

It works by packaging your project files, respecting ignore patterns, and sending them along with your prompts to a compatible LLM API. This allows you to ask complex questions, request refactoring, generate documentation, or debug issues with the AI having a complete picture of your code.

## Installation

### Windows

You can easily install HarzMind Code on Windows by running the following command in PowerShell. It will download the latest release, place it in a user-specific directory, and add it to your `PATH`.

```powershell
irm https://raw.githubusercontent.com/thxrsxm/harzmind-code/main/install.ps1 | iex
```

> **Note:** You may need to restart your terminal for the `PATH` changes to take effect.

### macOS

> COMING SOON

### Linux

> COMING SOON

## Getting Started: First-Time Setup

To start using HarzMind Code in your project, follow these steps:

1. **Navigate to your project's root directory** in your terminal.

   ```bash
   cd /path/to/your/project
   ```

2. **Initialize HarzMind Code:**
   Run the `init` command using the `-i` flag.

   ```bash
   hzmind -i
   ```

   This creates a new directory named `hzmind/` in your project root. It contains two essential files:

   *   `HZMIND.md`: Your main system prompt file.
   *   `.hzmignore`: A file to specify which files and directories to exclude from the context.

3. **Configure the System Prompt:**
   Open `hzmind/HZMIND.md` and write the instructions for the AI. Define its role, its goal, and any constraints. This is the most important step for getting good results.

4. **Start the application:**

   ```bash
   hzmind
   ```

5. **Add your API Account:**
   The first time you run the tool, you won't have an account configured. Use the `/acc new` command to add one. You will be prompted for a name, your API URL, and your API Key.

   **OpenAI example:**

   ```
   > /acc new
   Create account
   Name: openai
   API Url: https://api.openai.com/v1/chat/completions
   API Token: [your token will be hidden]
   Model (optional): gpt-4-turbo
   ```

   **Ollama example:**

   ```
   > /acc new
   Create account
   Name: ollama
   API Url: http://localhost:11434/v1/chat/completions
   API Token: ollama
   Model (optional): llama3.2:3b
   ```

6. **Login to your account** and start chatting!

   ```
   > /acc login my-openai
   Successfully logged in to 'my-openai'
   
   > What is the purpose of the `main.go` file?
   ```

## Core Concepts

### The `HZMIND.md` File

This file is the heart of HarzMind Code. Its content is sent to the LLM as the **system prompt** with every message you send. It defines the AI's persona, task, and rules. The entire codebase is appended to this prompt, giving the AI the context it needs.

**Example `HZMIND.md`:**

```markdown
You are an expert Go developer and a helpful programming assistant. Your task is to analyze the provided codebase and answer questions accurately.

When I ask for a code change, you must provide only the complete, updated content of the file(s) that need to be changed. Do not add explanations unless I explicitly ask for them.
```

### The `.hzmignore` File

This file works exactly like a `.gitignore` file. It tells HarzMind Code which files and directories to ignore when bundling the codebase for the LLM. This is crucial for keeping the context clean, focused, and within token limits.

The `.hzmignore` file is located at `hzmind/.hzmignore` inside your project.

By default, the following patterns are **always ignored**:

- `.git`
- `.idea`
- `.vscode`
- `node_modules`
- `vendor`
- `*.exe`
- `config.xml`
- `HZMIND.md` (the prompt file itself)
- The `hzmind/` directory

You can add any other patterns (e.g., `build/`, `*.log`, `*.tmp`) to your `hzmind/.hzmignore` file.

### Configuration & API Keys

All account configurations, including your API credentials, are stored in a `config.json` file.

*   **Location:** The file is stored in a platform-specific data directory:
    *   **macOS:** `~/Library/Application Support/hzmind/`
    *   **Linux:** `~/.config/hzmind/`
    *   **Windows:** The same directory as the `hzmind.exe` executable. If you used the PowerShell installer, this will be `%LOCALAPPDATA%\HarzMindCode\`.
*   **Security:** Your API keys are stored in plain text in this file. Ensure that this directory is secure and not synced to public repositories.
*   **Management:** You should not edit this file manually. Use the `/acc` commands within the application to manage your accounts safely.

## Usage

### Command-Line Flags

Flags are used when launching the application.

| Flag | Description                                                  |
| :--- | :----------------------------------------------------------- |
| `-h` | Display the help message with all available flags.           |
| `-i` | **Init Project**: Creates the `hzmind` directory and its files. |
| `-v` | Show the application's version and build date.               |
| `-o` | **Output**: Write the entire conversation to a timestamped Markdown file in the `hzmind/out/` directory. |

### REPL Commands

Commands are used inside the application's interactive prompt and start with a `/`.

| Command                        | Description                                                  |
| :----------------------------- | :----------------------------------------------------------- |
| `/help`                        | List all available REPL commands.                            |
| `/exit`                        | Quit the application.                                        |
| `/init`                        | Initializes the project (same as the `-i` flag).             |
| `/clear`                       | Clears the current chat history, starting a fresh conversation (but keeps the system prompt and codebase). |
| `/info`                        | Show application info, version, and author.                  |
| `/session`                     | Show current session info including account, model, directory, and token count. |
| `/tree`                        | Display the project's file structure as a tree, respecting ignore patterns. |
| `/models`                      | List all available models from the currently logged-in account's API. |
| `/model <model_name>`          | Change the LLM model for the current session (e.g., `/model gpt-3.5-turbo`). |
| `/bash <command>`              | Execute a shell command and display the output (e.g., `/bash ls -l`). |
| `/editor <editor_name> [file]` | Open a file in a terminal-based editor (e.g., `/editor nano internal/api/api.go`). |
| `/acc`                         | List all configured accounts.                                |
| `/acc new`                     | Start the wizard to create a new account (prompts for name, URL, key, model). |
| `/acc login <account_name>`    | Log in to a specific account to make it active.              |
| `/acc logout`                  | Log out of the current account.                              |
| `/acc remove <account_name>`   | Delete a configured account.                                 |
| `/acc info <account_name>`     | Show details for a specific account (without the API key).   |

## Example Workflow

1. **Project Goal:** Refactor a function `GetCodeBase` in `internal/codebase/codebase.go` to be more efficient.

2. **Setup:**

   *   `cd /my/go/project`
   *   `hzmind -i`
   *   Edit `hzmind/HZMIND.md` to instruct the AI to act as a senior Go performance expert.
   *   Run `hzmind` and log in with `/acc login <my_account>`.

3. **Interaction:**

   ```
   > Analyze the function `GetCodeBase` in `internal/codebase/codebase.go`. Is there any way to improve its performance or make it more idiomatic?
   
   [LLM responds with analysis and suggestions...]
   
   > Great. Please provide the complete, refactored version of the `internal/codebase/codebase.go` file with your suggested improvements.
   
   [LLM provides the full file content, which you can copy and paste]
   ```

## License

MIT License

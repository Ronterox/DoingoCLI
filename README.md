# DoingCLI - A Fast Go CLI for Task & Time Tracking

DoingCLI is a command-line utility written in Go designed for efficient task and time tracking. It serves as a faster, more robust replacement for existing tools, focusing on seamless task management and performance.

## Features

-   **Quick Task Annotation:** Mark tasks as "now" or "later".
-   **Effortless Completion:** Easily mark tasks as "done".
-   **Activity Overview:** View recent tasks and your last active task.
-   **Local Storage:** All tasks are stored locally in `what_was_I_doing.txt`.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

-   Go (version 1.18 or higher recommended)

### Installation & Building

Clone the repository:
```bash
git clone https://github.com/your-username/DoingCLI.git # Replace with actual repo URL
cd DoingCLI
```

Build the application using the provided `start` script:
```bash
./start build
```
This will create an executable named `doingcli` in the project root.

### Running

To run the application:
```bash
./start run
```
Or, after building, you can run the executable directly:
```bash
./doingcli [command]
```

For development with live-reloading:
```bash
./start dev
```

## Usage

DoingCLI operates with simple, intuitive commands.

### `now <task description>`
Records a new task that you are currently working on.
Example:
```bash
./doingcli now Working on the new feature
```

### `later <task description>`
Records a task that you plan to do later.
Example:
```bash
./doingcli later Research Go concurrency patterns
```

### `done`
Marks your last recorded task as complete.
Example:
```bash
./doingcli done
```

### `recent`
Displays a list of your recently completed or active tasks.
Example:
```bash
./doingcli recent
```

### `last`
Shows the very last task you recorded.
Example:
```bash
./doingcli last
```

### `random`
Displays a random task from your recent unfinished tasks.
Example:
```bash
./doingcli random
```

## Contributing

Contributions are welcome!

## Planned Features

The following features are planned for future development:

-   `edit <task-id> <new-description>`: Modify an existing task.
-   `archive`: Move completed tasks to an archive.
-   `undo`: Revert the last action.
-   `help`: Display help information for commands.
-   Improved handling of "later" tasks, possibly separating them into a dedicated view.
-   Performance metrics and tracking for tasks.

---

#### What makes it finished?

The project will be considered feature-complete when it offers robust task management, including the ability to annotate tasks (now/later), mark as done, view recent/last activities, and manage distractions. Key aspects include being super fast, showing task duration, and providing a seamless cycle for completing tasks.

#### Notes:

> This project is driven by the need for a highly performant and user-friendly CLI task tracker. It aims to address the shortcomings of existing tools, such as slow performance, unreliable commands, and difficulties in managing and completing old tasks efficiently. This is an exciting opportunity to learn and implement new concepts in Go.

---

#### Deadline:

Will be finished in five days by Tue 25 Nov 20:23:12 CET 2025, else it won't be finished.
*(Note: This deadline is a personal goal for initial feature set completion.)*

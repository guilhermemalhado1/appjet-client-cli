package main

import (
	"appjet-client-cli/configurations"
	"appjet-client-cli/requests"
	"bytes"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

type CommandHandler func(args []string)

var commands = make(map[string]CommandHandler)

func init() {
	commands["help"] = help
	commands["start"] = start
	commands["ls"] = logout
	commands["login"] = login
	commands["logout"] = logout
	commands["inspect"] = inspect
	commands["show"] = show
	commands["rename"] = rename
	commands["delete"] = delete
	commands["create"] = create

	configurations.LoadConfig()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing command argument")
	}

	command := os.Args[1]
	args := os.Args[2:]
	handler, availableCommand := commands[command]

	if !availableCommand {
		log.Fatalf("Invalid command: %s", command)
	}

	handler(args)
}

func help(args []string) {
	fmt.Println("Available commands: prefix: appjet ")
	for command, _ := range commands {
		commandUsage := getCommandUsage(command)
		fmt.Println("  " + command + commandUsage)
	}
}

func getCommandUsage(command string) string {
	switch command {
	case "delete":
		return "   <app-name> <token>"
	case "create":
		return "   <token>"
	case "inspect":
		return "  <app-name> <token>"
	case "rename":
		return "   <old-app-name> <new-app-name> <token>"
	case "show":
		return "     <app-name> <token>"
	case "login":
		return "    <username>"
	case "start":
		return "    <token>"
	case "ls":
		return "       <token>"
	case "logout":
		return "   <token>"
	default:
		return ""
	}
}

func start(args []string) {
	fmt.Println("Starting app...")
	if len(args) > 0 {
		fmt.Println("App name:", args[0])
	}
}

func config(args []string) {
	fmt.Println("Listing configurations...")
}

func login(args []string) {
	if len(args) != 1 {
		fmt.Println("Invalid arguments for login command")
		return
	}

	username := args[0]
	fmt.Println("Trying to log in appjet-project.com ...")
	fmt.Println("Username:", username)

	fmt.Print("Enter password: ")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}
	fmt.Println()

	requests.DoLogin(username, string(password))
}

func logout(args []string) {

	fmt.Println("Logging out...")

	deleteLoginState()
}

func inspect(args []string) {
	if len(args) != 1 {
		fmt.Println("Invalid arguments for inspect command")
		return
	}

	appName := args[0]
	fmt.Println("Inspecting app:", appName)
}

func show(args []string) {

	if len(args) != 1 {
		fmt.Println("Invalid arguments for show command")
		return
	}

	appName := args[0]
	fmt.Println("Showing app details:", appName)
}

func rename(args []string) {

	if len(args) != 2 {
		fmt.Println("Invalid arguments for rename command")
		return
	}

	oldName := args[0]
	newName := args[1]
	fmt.Println("Renaming app from", oldName, "to", newName)
}

func delete(args []string) {

	if len(args) != 1 {
		fmt.Println("Invalid arguments for delete command")
		return
	}

	appName := args[0]
	fmt.Println("Deleting app:", appName)
}

func requireLogin(handler CommandHandler) CommandHandler {
	return func(args []string) {
		handler(args)
	}
}

func deleteLoginState() {

	requests.DeleteLoginState()
}

func create(args []string) {
	fmt.Print("Enter username: ")
	username := getUserInput()

	fmt.Print("Enter email: ")
	email := getUserInput()

	password := getPasswordInput()

	requests.SignupUser(username, password, email)
}

func getUserInput() string {
	var input string
	fmt.Scanln(&input)
	return input
}

func getPasswordInput() string {
	fmt.Print("Enter password (hidden): ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		fmt.Println("Error reading password:", err)
		return ""
	}
	fmt.Println()
	return string(bytes.TrimSpace(password))
}

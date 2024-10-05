package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/BurntSushi/toml"
)

var setting Setting

func load() error {
	rootDir := filepath.Join(os.Getenv("USERPROFILE"), "AppWrapper")
	file := filepath.Join(rootDir, "config.toml")
	_, err := os.Stat(file)
	if err == nil {
		// ファイルが存在している場合
		_, err := toml.DecodeFile(file, &setting)
		if err != nil {
			return err
		}
		return nil
	}
	if err := os.MkdirAll(rootDir, 0700); err != nil {
		return fmt.Errorf("設定ディレクトリを作成できませんでした")
	}
	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return fmt.Errorf("設定ファイルを作成しました。\nファイル：%s", file)
}

func createMessageWindow(errm error) {
	dll, err := syscall.LoadDLL("user32.dll")
	defer dll.Release()

	proc, err := dll.FindProc("MessageBoxW")
	if err != nil {
		log.Fatal(err)
	}
	message, err := syscall.UTF16PtrFromString(fmt.Sprint(errm))
	if err != nil {
		log.Fatal(err)
	}
	title, err := syscall.UTF16PtrFromString("AppWrapper")
	if err != nil {
		log.Fatal(err)
	}
	proc.Call(
		0,
		uintptr(unsafe.Pointer(message)),
		uintptr(unsafe.Pointer(title)),
		0)
}

func main() {
	err := load()
	if err != nil {
		createMessageWindow(err)
	}

	filename := filepath.Base(os.Args[0][:len(os.Args[0]) - len(filepath.Ext(os.Args[0]))])
	for cmd, command := range setting.Cmd {
		if cmd != filename {
			continue
		}
		cmd := exec.Command(command.Command)
		env := []string{}
		for key, value := range command.Env {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = append(os.Environ(), env...)
		cmd.Run()
	}
}

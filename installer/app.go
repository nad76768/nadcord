/*
 * SPDX-License-Identifier: GPL-3.0
 * Nadcord Installer - web based UI (Wails)
 * Based on Vencord Installer, Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend
var frontendFiles embed.FS

// Errors raised deep inside the shared installer code are collected here so
// they can be returned to the UI instead of being shown in a native dialog.
var (
	lastDownloadErr error
	scuffedInstall  bool
)

func InstallLatestBuilds() error {
	if IsDevInstall {
		return nil
	}
	lastDownloadErr = installLatestBuilds()
	return lastDownloadErr
}

func HandleScuffedInstall() {
	scuffedInstall = true
}

// DiscordInfo is what the UI receives for each Discord install.
type DiscordInfo struct {
	Path       string `json:"path"`
	Branch     string `json:"branch"`
	IsPatched  bool   `json:"isPatched"`
	IsOpenAsar bool   `json:"isOpenAsar"`
}

type State struct {
	LatestHash    string        `json:"latestHash"`
	InstalledHash string        `json:"installedHash"`
	GithubError   string        `json:"githubError"`
	FilesDirError string        `json:"filesDirError"`
	Discords      []DiscordInfo `json:"discords"`
}

// App holds the methods the UI can call (window.go.main.App.*).
type App struct {
	ctx         context.Context
	githubOnce  sync.Once
	busy        sync.Mutex
	customPaths []string
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) waitForGithub() {
	a.githubOnce.Do(func() { <-GithubDoneChan })
}

func toInfo(di *DiscordInstall) DiscordInfo {
	return DiscordInfo{
		Path:       di.path,
		Branch:     di.branch,
		IsPatched:  di.isPatched,
		IsOpenAsar: di.IsOpenAsar(),
	}
}

func parseAny(p string) *DiscordInstall {
	if di := ParseDiscord(p, ""); di != nil {
		return di
	}
	return ParseDiscordNew(p, "", strings.Contains(p, "com.discordapp"))
}

func (a *App) installs() []*DiscordInstall {
	var list []*DiscordInstall
	seen := map[string]bool{}
	for _, d := range FindDiscords() {
		di := d.(*DiscordInstall)
		list = append(list, di)
		seen[di.path] = true
	}
	sort.Slice(list, func(i, j int) bool { return list[i].branch < list[j].branch })
	for _, p := range a.customPaths {
		if seen[p] {
			continue
		}
		if di := parseAny(p); di != nil {
			list = append(list, di)
			seen[p] = true
		}
	}
	return list
}

func (a *App) find(p string) (*DiscordInstall, error) {
	for _, di := range a.installs() {
		if di.path == p {
			return di, nil
		}
	}
	if di := parseAny(p); di != nil {
		return di, nil
	}
	return nil, errors.New("this folder is not a Discord install anymore. Pick it again")
}

func friendly(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, os.ErrPermission) {
		return errors.New("permission denied. Fully close Discord (also from the system tray) and try again")
	}
	return err
}

// GetState returns the latest build info and every Discord install found.
func (a *App) GetState() State {
	a.waitForGithub()
	s := State{
		LatestHash:    LatestHash,
		InstalledHash: InstalledHash,
		Discords:      []DiscordInfo{},
	}
	if GithubError != nil {
		s.GithubError = GithubError.Error()
	}
	if FilesDirErr != nil {
		s.FilesDirError = FilesDirErr.Error()
	}
	for _, di := range a.installs() {
		s.Discords = append(s.Discords, toInfo(di))
	}
	return s
}

func (a *App) patch(p string, forceDownload bool) error {
	a.busy.Lock()
	defer a.busy.Unlock()

	di, err := a.find(p)
	if err != nil {
		return err
	}

	scuffedInstall = false
	if CheckScuffedInstall() || scuffedInstall {
		return errors.New("your Discord install is broken (Discord installed itself in the wrong place). Reinstall Discord, then try again")
	}

	if forceDownload {
		if err := InstallLatestBuilds(); err != nil {
			return errors.New("couldn't download the latest Nadcord build: " + err.Error())
		}
	}

	lastDownloadErr = nil
	err = di.patch()
	if lastDownloadErr != nil {
		return errors.New("couldn't download the latest Nadcord build: " + lastDownloadErr.Error())
	}
	return friendly(err)
}

// Install installs Nadcord into the Discord at path p.
func (a *App) Install(p string) error {
	return a.patch(p, false)
}

// Repair downloads the latest build again and reinstalls it.
func (a *App) Repair(p string) error {
	return a.patch(p, true)
}

// Uninstall removes Nadcord from the Discord at path p.
func (a *App) Uninstall(p string) error {
	a.busy.Lock()
	defer a.busy.Unlock()

	di, err := a.find(p)
	if err != nil {
		return err
	}
	return friendly(di.unpatch())
}

// ToggleOpenAsar installs OpenAsar, or removes it if it's already there.
func (a *App) ToggleOpenAsar(p string) error {
	a.busy.Lock()
	defer a.busy.Unlock()

	di, err := a.find(p)
	if err != nil {
		return err
	}
	if di.IsOpenAsar() {
		return friendly(di.UninstallOpenAsar())
	}
	return friendly(di.InstallOpenAsar())
}

// PickFolder lets the user choose a Discord folder that wasn't found automatically.
func (a *App) PickFolder() (*DiscordInfo, error) {
	dir, err := wruntime.OpenDirectoryDialog(a.ctx, wruntime.OpenDialogOptions{
		Title: "Choose your Discord folder",
	})
	if err != nil || dir == "" {
		return nil, err
	}
	di := parseAny(dir)
	if di == nil {
		return nil, errors.New("that folder isn't a Discord install. Choose the folder that contains the app-… folders (for example AppData\\Local\\Discord)")
	}
	a.customPaths = append(a.customPaths, di.path)
	info := toInfo(di)
	return &info, nil
}

// OpenURL opens a link in the default browser.
func (a *App) OpenURL(url string) {
	wruntime.BrowserOpenURL(a.ctx, url)
}

func main() {
	InitGithubDownloader()

	assets, err := fs.Sub(frontendFiles, "frontend")
	if err != nil {
		panic(err)
	}

	app := &App{}
	err = wails.Run(&options.App{
		Title:            "Nadcord Installer",
		Width:            920,
		Height:           600,
		DisableResize:    true,
		Frameless:        true,
		BackgroundColour: &options.RGBA{R: 0x15, G: 0x1b, B: 0x0f, A: 0xff},
		AssetServer:      &assetserver.Options{Assets: assets},
		OnStartup:        app.startup,
		Bind:             []any{app},
		Windows: &windows.Options{
			Theme: windows.Dark,
		},
	})
	if err != nil {
		Log.Error("Failed to start the installer:", err)
	}
}

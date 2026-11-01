//go:build !cli

/*
 * SPDX-License-Identifier: GPL-3.0
 * Vencord Installer, a cross platform gui/cli app for installing Vencord
 * Copyright (c) 2023 Vendicated and Vencord contributors
 */

package main

import (
	"bytes"
	_ "embed"
	"errors"
	"image"
	"image/color"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"

	// png decoder for icon
	_ "image/png"
	"os"
	path "path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var (
	discords        []any
	radioIdx        int
	customChoiceIdx int

	customDir              string
	autoCompleteDir        string
	autoCompleteFile       string
	autoCompleteCandidates []string
	autoCompleteIdx        int
	lastAutoComplete       string
	didAutoComplete        bool

	modalId          = 0
	modalTitle       = "Oh No :("
	modalMessage     = "You should never see this"
	modalExtra       = ""
	modalSupportLink = ""

	acceptedOpenAsar   bool
	showedUpdatePrompt bool
	showCustomLocation bool

	win *g.MasterWindow
)

//go:embed winres/icon.png
var iconBytes []byte

func init() {
	LogLevel = LevelDebug
}

func main() {
	InitGithubDownloader()
	discords = FindDiscords()
	customChoiceIdx = len(discords)

	go func() {
		<-GithubDoneChan
		g.Update()
	}()

	go func() {
		<-SelfUpdateCheckDoneChan
		g.Update()
	}()

	win = g.NewMasterWindow("Nadcord Installer", 1200, 800, 0)

	icon, _, err := image.Decode(bytes.NewReader(iconBytes))
	if err != nil {
		Log.Warn("Failed to load application icon", err)
		Log.Debug(iconBytes, len(iconBytes))
	} else {
		win.SetIcon([]image.Image{icon})
	}
	win.SetBgColor(ThemeBg)

	win.Run(loop)
}

type CondWidget struct {
	predicate  bool
	ifWidget   func() g.Widget
	elseWidget func() g.Widget
}

func (w *CondWidget) Build() {
	if w.predicate {
		w.ifWidget().Build()
	} else if w.elseWidget != nil {
		w.elseWidget().Build()
	}
}

func getChosenInstall() *DiscordInstall {
	if radioIdx != customChoiceIdx {
		return discords[radioIdx].(*DiscordInstall)
	}

	if discord := ParseDiscord(customDir, ""); discord != nil {
		return discord
	}

	if discord := ParseDiscordNew(customDir, "", strings.Contains(customDir, "com.discordapp")); discord != nil {
		return discord
	}

	g.OpenPopup("#invalid-custom-location")
	return nil
}

func InstallLatestBuilds() (err error) {
	if IsDevInstall {
		return
	}

	err = installLatestBuilds()
	if err != nil {
		ShowErrorModal("Failed to install the latest Nadcord builds from GitHub", err.Error(), "https://github.com/nad76768/nadcord/issues")
	}
	return
}

func handlePatch() {
	choice := getChosenInstall()
	if choice != nil {
		choice.Patch()
	}
}

func handleUnpatch() {
	choice := getChosenInstall()
	if choice != nil {
		choice.Unpatch()
	}
}

func handleOpenAsar() {
	if acceptedOpenAsar || getChosenInstall().IsOpenAsar() {
		handleOpenAsarConfirmed()
		return
	}

	g.OpenPopup("#openasar-confirm")
}

func handleOpenAsarConfirmed() {
	choice := getChosenInstall()
	if choice != nil {
		if choice.IsOpenAsar() {
			if err := choice.UninstallOpenAsar(); err != nil {
				handleErr(choice, err, "uninstall OpenAsar from")
			} else {
				g.OpenPopup("#openasar-unpatched")
				g.Update()
			}
		} else {
			if err := choice.InstallOpenAsar(); err != nil {
				handleErr(choice, err, "install OpenAsar on")
			} else {
				g.OpenPopup("#openasar-patched")
				g.Update()
			}
		}
	}
}

func handleErr(di *DiscordInstall, err error, action string) {
	if errors.Is(err, os.ErrPermission) {
		switch runtime.GOOS {
		case "windows":
			err = errors.New("Permission denied. Make sure your Discord is fully closed (from the tray)!")
		case "darwin":
			HandleInsufficientPermissions()
			return
		default:
			err = errors.New("Permission denied. Maybe try running me as Administrator?")
		}
	}

	ShowErrorModal("Failed to "+action+" this Install.", err.Error(), "https://github.com/nad76768/nadcord/issues")
}

func HandleScuffedInstall() {
	g.OpenPopup("#scuffed-install")
}

func HandleInsufficientPermissions() {
	g.OpenPopup("#insufficient-permissions")
}

func (di *DiscordInstall) Patch() {
	if CheckScuffedInstall() {
		return
	}
	if err := di.patch(); err != nil {
		handleErr(di, err, "patch")
	} else {
		g.OpenPopup("#patched")
	}
}

func (di *DiscordInstall) Unpatch() {
	if err := di.unpatch(); err != nil {
		handleErr(di, err, "unpatch")
	} else {
		g.OpenPopup("#unpatched")
	}
}

func onCustomInputChanged() {
	p := customDir
	if len(p) != 0 {
		// Select the custom option for people
		radioIdx = customChoiceIdx
	}

	dir := path.Dir(p)

	isNewDir := strings.HasSuffix(p, "/")
	wentUpADir := !isNewDir && dir != autoCompleteDir

	if isNewDir || wentUpADir {
		autoCompleteDir = dir
		// reset all the funnies
		autoCompleteIdx = 0
		lastAutoComplete = ""
		autoCompleteFile = ""
		autoCompleteCandidates = nil

		// Generate autocomplete items
		files, err := os.ReadDir(dir)
		if err == nil {
			for _, file := range files {
				autoCompleteCandidates = append(autoCompleteCandidates, file.Name())
			}
		}
	} else if !didAutoComplete {
		// reset auto complete and update our file
		autoCompleteFile = path.Base(p)
		lastAutoComplete = ""
	}

	if wentUpADir {
		autoCompleteFile = path.Base(p)
	}

	didAutoComplete = false
}

// go can you give me []any?
// to pass to giu RangeBuilder?
// yeeeeees
// actually returns []string like a boss
func makeAutoComplete() []any {
	input := strings.ToLower(autoCompleteFile)

	var candidates []any
	for _, e := range autoCompleteCandidates {
		file := strings.ToLower(e)
		if autoCompleteFile == "" || strings.HasPrefix(file, input) {
			candidates = append(candidates, e)
		}
	}
	return candidates
}

func makeRadioOnChange(i int) func() {
	return func() {
		radioIdx = i
	}
}

func renderFilesDirErr() g.Widget {
	return g.Layout{
		g.Dummy(0, 50),
		g.Style().
			SetColor(g.StyleColorText, DiscordRed).
			SetFontSize(30).
			To(
				g.Align(g.AlignCenter).To(
					g.Label("Error: Failed to create: "+FilesDirErr.Error()),
					g.Label("Resolve this error, then restart me!"),
				),
			),
	}
}

func Tooltip(label string) g.Widget {
	return g.Style().
		SetColor(g.StyleColorText, ThemeText).
		SetColor(g.StyleColorPopupBg, ThemeSurface).
		SetColor(g.StyleColorBorder, NeonLimeDim).
		SetStyle(g.StyleVarWindowPadding, 10, 8).
		SetStyleFloat(g.StyleVarWindowRounding, 8).
		To(
			g.Tooltip(label),
		)
}

func InfoModal(id, title, description string) g.Widget {
	return RawInfoModal(id, title, description, "", false)
}

func InfoModalExtra(id, title, description, extra string) g.Widget {
	return RawInfoModal(id, title, description, extra, false)
}

func ErrorModal(id, title, error, supportLink string) g.Widget {
	return g.Style().
		SetStyle(g.StyleVarWindowPadding, 30, 30).
		SetStyleFloat(g.StyleVarWindowRounding, 12).
		To(
			g.PopupModal(id).
				Flags(g.WindowFlagsNoTitleBar|g.WindowFlagsAlwaysAutoResize).
				Layout(
					g.Style().SetFontSize(35).SetColor(g.StyleColorText, DangerText).To(g.Label(title)),
					g.Dummy(0, 5),
					g.Style().SetFontSize(25).To(
						g.Row(
							g.Label("If this issue persists, visit"),
							g.Style().
								SetColor(g.StyleColorButton, DiscordBlue).
								SetStyle(g.StyleVarFramePadding, 4, 4).
								To(
									g.Button("our support page").OnClick(func() {
										g.OpenURL(supportLink)
									}),
								),
						),
					),
					g.Dummy(0, 20),
					g.Style().SetFontSize(20).To(g.Label("Error Information")),
					g.Style().SetFontSize(18).To(g.Label(error)),
					g.Dummy(0, 20),
					g.Button("Ok").
						OnClick(func() { g.CloseCurrentPopup() }).
						Size(100, 30),
				),
		)
}

func RawInfoModal(id, title, description, extra string, isOpenAsar bool) g.Widget {
	isDynamic := strings.HasPrefix(id, "#modal") && extra != ""
	return g.Style().
		SetStyle(g.StyleVarWindowPadding, 30, 30).
		SetStyleFloat(g.StyleVarWindowRounding, 12).
		To(
			g.PopupModal(id).
				Flags(g.WindowFlagsNoTitleBar|Ternary(isDynamic, g.WindowFlagsAlwaysAutoResize, 0)).
				Layout(
					g.Style().SetFontSize(30).SetColor(g.StyleColorText, NeonLime).To(
						g.Label(title),
					),
					g.Dummy(0, 5),
					g.Style().SetFontSize(20).To(
						g.Label(description),
					),
					&CondWidget{extra != "", func() g.Widget {
						return g.Column(
							g.Dummy(0, 10),
							g.Style().SetFontSize(20).To(
								g.Label(extra).Wrapped(true),
							),
						)
					}, nil},
					&CondWidget{id == "#scuffed-install", func() g.Widget {
						return g.Column(
							g.Dummy(0, 10),
							g.Button("Take me there!").OnClick(func() {
								// this issue only exists on windows so using Windows specific path is oki
								username := os.Getenv("USERNAME")
								programData := os.Getenv("PROGRAMDATA")
								g.OpenURL("file://" + path.Join(programData, username))
							}).Size(200, 30),
						)
					}, nil},
					g.Dummy(0, 20),
					&CondWidget{id == "#insufficient-permissions", func() g.Widget {
						return g.Column(
							g.Dummy(0, 10),
							g.Button("Open Settings").OnClick(func() {
								// "App Management" permissions doesnt exist on Monterey, just use full disk access for now
								g.OpenURL("x-apple.systempreferences:com.apple.preference.security?Privacy_AllFiles")
							}).Size(200, 30),
						)
					}, nil},
					&CondWidget{isOpenAsar,
						func() g.Widget {
							return g.Row(
								g.Button("Accept").
									OnClick(func() {
										acceptedOpenAsar = true
										g.CloseCurrentPopup()
									}).
									Size(100, 30),
								g.Button("Cancel").
									OnClick(func() {
										g.CloseCurrentPopup()
									}).
									Size(100, 30),
							)
						},
						func() g.Widget {
							return g.Button("Ok").
								OnClick(func() {
									g.CloseCurrentPopup()
								}).
								Size(100, 30)
						},
					},
				),
		)
}

func UpdateModal() g.Widget {
	return g.Style().
		SetStyle(g.StyleVarWindowPadding, 30, 30).
		SetStyleFloat(g.StyleVarWindowRounding, 12).
		To(
			g.PopupModal("#update-prompt").
				Flags(g.WindowFlagsNoTitleBar | g.WindowFlagsAlwaysAutoResize).
				Layout(
					g.Align(g.AlignCenter).To(
						g.Style().SetFontSize(30).To(
							g.Label("Your Installer is outdated!"),
						),
						g.Style().SetFontSize(20).To(
							g.Label(
								"Would you like to update now?\n\n"+
									"Once you press Update Now, the new installer will automatically be downloaded.\n"+
									"The installer will temporarily seem unresponsive. Just wait!\n"+
									"Once the update is done, the Installer will automatically reopen.\n\n"+
									"On MacOs, Auto updates are not supported, so it will instead open in browser.",
							),
						),
						g.Row(
							g.Button("Update Now").
								OnClick(func() {
									if runtime.GOOS == "darwin" {
										g.CloseCurrentPopup()
										g.OpenURL(GetInstallerDownloadLink())
										return
									}

									err := UpdateSelf()
									g.CloseCurrentPopup()

									if err != nil {
										ShowModal("Failed to update self!", "Please manually download the latest Installer.", err.Error())
									} else {
										if err = RelaunchSelf(); err != nil {
											ShowModal("Failed to restart self!", "Please manually restart the Installer.", err.Error())
										}
									}
								}).
								Size(100, 30),
							g.Button("Later").
								OnClick(func() {
									g.CloseCurrentPopup()
								}).
								Size(100, 30),
						),
					),
				),
		)
}

func ShowModal(title, desc, extra string) {
	modalTitle = title
	modalMessage = desc
	modalExtra = extra
	modalId++
	g.OpenPopup("#modal" + strconv.Itoa(modalId))
}

func ShowErrorModal(title, extra, supportLink string) {
	modalTitle = title
	modalExtra = extra
	modalSupportLink = supportLink
	modalId++
	g.OpenPopup("#error-modal" + strconv.Itoa(modalId))
}

func renderInstaller() g.Widget {
	candidates := makeAutoComplete()

	wi, _ := win.GetSize()
	w := float32(wi) - 96

	var currentDiscord *DiscordInstall
	if radioIdx != customChoiceIdx {
		currentDiscord = discords[radioIdx].(*DiscordInstall)
	}
	var isOpenAsar = currentDiscord != nil && currentDiscord.IsOpenAsar()

	if CanUpdateSelf() && !showedUpdatePrompt {
		showedUpdatePrompt = true
		g.OpenPopup("#update-prompt")
	}

	layout := g.Layout{
		g.Style().SetFontSize(20).To(
			renderErrorCard(
				ThemeSurfaceLime,
				ThemeText,
				"**github.com/nad76768/nadcord** is the only official place to get Nadcord. Any other site claiming to be us is malicious.\n"+
					"If you downloaded from any other source, you should delete / uninstall everything immediately, run a malware scan and change your Discord password.",
				90,
			),
		),

		g.Dummy(0, 20),

		g.Style().SetFontSize(28).SetColor(g.StyleColorText, NeonLime).To(
			g.Label("Select your Discord"),
		),
		g.Style().SetFontSize(18).SetColor(g.StyleColorText, ThemeTextDim).To(
			g.Label("Fully close Discord before installing, then pick the install you want to patch."),
		),
		g.Dummy(0, 10),

		&CondWidget{len(discords) == 0, func() g.Widget {
			s := "No Discord installs found. You first need to install Discord."
			if runtime.GOOS == "linux" {
				s += " snap is not supported."
			}

			return &CondWidget{!showCustomLocation, func() g.Widget {
				return g.Column(
					g.Style().SetFontSize(25).To(g.Label(s)),
					g.Dummy(0, 10),
					g.Checkbox("I am an advanced user and have Discord installed at a different location", &showCustomLocation),
				)
			}, nil}
		}, nil},

		g.Style().SetFontSize(20).To(
			g.RangeBuilder("Discords", discords, func(i int, v any) g.Widget {
				d := v.(*DiscordInstall)
				var text string
				switch d.branch {
				case "ptb":
					text = "Discord PTB"
				case "canary":
					text = "Discord Canary"
				case "development":
					text = "Discord Development"
				default:
					text = "Discord"
				}

				if d.isPatched {
					text += " (Nadcord Installed)"
				}

				return g.Row(
					g.RadioButton(text, radioIdx == i).OnChange(makeRadioOnChange(i)),
					g.Style().SetColor(g.StyleColorText, ThemeTextDim).To(g.Label(" "+d.path)),
				)
			}),
			&CondWidget{showCustomLocation, func() g.Widget {
				return g.RadioButton("Custom Install Location", radioIdx == customChoiceIdx).OnChange(makeRadioOnChange(customChoiceIdx))
			}, nil},
		),

		&CondWidget{showCustomLocation, func() g.Widget {
			return g.Column(
				g.Dummy(0, 5),
				g.Style().
					SetStyle(g.StyleVarFramePadding, 16, 16).
					SetFontSize(20).
					To(
						g.InputText(&customDir).Hint("The custom location").
							Size(w-16).
							Flags(g.InputTextFlagsCallbackCompletion).
							OnChange(onCustomInputChanged).
							// this library has its own autocomplete but it's broken
							Callback(
								func(data imgui.InputTextCallbackData) int32 {
									if len(candidates) == 0 {
										return 0
									}
									// just wrap around
									if autoCompleteIdx >= len(candidates) {
										autoCompleteIdx = 0
									}

									// used by change handler
									didAutoComplete = true

									start := len(customDir)
									// Delete previous auto complete
									if lastAutoComplete != "" {
										start -= len(lastAutoComplete)
										data.DeleteBytes(start, len(lastAutoComplete))
									} else if autoCompleteFile != "" { // delete partial input
										start -= len(autoCompleteFile)
										data.DeleteBytes(start, len(autoCompleteFile))
									}

									// Insert auto complete
									lastAutoComplete = candidates[autoCompleteIdx].(string)
									data.InsertBytes(start, []byte(lastAutoComplete))
									autoCompleteIdx++

									return 0
								},
							),
					),
				g.RangeBuilder("AutoComplete", candidates, func(i int, v any) g.Widget {
					dir := v.(string)
					return g.Label(dir)
				}),
			)
		}, nil},

		g.Dummy(0, 20),

		g.Style().SetFontSize(20).To(
			g.Row(
				g.Style().
					SetColor(g.StyleColorButton, NeonLime).
					SetColor(g.StyleColorButtonHovered, NeonLimeHovered).
					SetColor(g.StyleColorButtonActive, NeonLimeActive).
					SetColor(g.StyleColorText, ThemeBlack).
					SetStyle(g.StyleVarFrameRounding, 12, 12).
					SetDisabled(GithubError != nil).
					To(
						g.Button("Install").
							OnClick(handlePatch).
							Size((w-40)/4, 50),
						Tooltip("Install Nadcord into the selected Discord"),
					),
				g.Style().
					SetColor(g.StyleColorButton, ThemeSurfaceLime).
					SetColor(g.StyleColorButtonHovered, ThemeSurfaceHover).
					SetColor(g.StyleColorButtonActive, ThemeSurface).
					SetColor(g.StyleColorText, NeonLime).
					SetColor(g.StyleColorBorder, NeonLime).
					SetStyleFloat(g.StyleVarFrameBorderSize, 2).
					SetStyle(g.StyleVarFrameRounding, 12, 12).
					SetDisabled(GithubError != nil).
					To(
						g.Button("Reinstall / Repair").
							OnClick(func() {
								if IsDevInstall {
									handlePatch()
								} else {
									err := InstallLatestBuilds()
									if err == nil {
										handlePatch()
									}
								}
							}).
							Size((w-40)/4, 50),
						Tooltip("Reinstall & Update Nadcord"),
					),
				g.Style().
					SetColor(g.StyleColorButton, DangerBg).
					SetColor(g.StyleColorButtonHovered, DangerBgHover).
					SetColor(g.StyleColorButtonActive, DangerBg).
					SetColor(g.StyleColorText, DangerText).
					SetColor(g.StyleColorBorder, DangerText).
					SetStyleFloat(g.StyleVarFrameBorderSize, 1).
					SetStyle(g.StyleVarFrameRounding, 12, 12).
					To(
						g.Button("Uninstall").
							OnClick(handleUnpatch).
							Size((w-40)/4, 50),
						Tooltip("Remove Nadcord from the selected Discord"),
					),
				g.Style().
					SetColor(g.StyleColorButton, ThemeSurface).
					SetColor(g.StyleColorButtonHovered, ThemeSurfaceHover).
					SetColor(g.StyleColorButtonActive, ThemeSurface).
					SetColor(g.StyleColorText, ThemeTextDim).
					SetColor(g.StyleColorBorder, NeonLimeFaint).
					SetStyleFloat(g.StyleVarFrameBorderSize, 1).
					SetStyle(g.StyleVarFrameRounding, 12, 12).
					To(
						g.Button(Ternary(isOpenAsar, "Uninstall OpenAsar", Ternary(currentDiscord != nil, "Install OpenAsar", "(Un-)Install OpenAsar"))).
							OnClick(handleOpenAsar).
							Size((w-40)/4, 50),
						Tooltip("Manage OpenAsar"),
					),
			),
		),

		g.Dummy(0, 30),
		g.Style().SetColor(g.StyleColorSeparator, NeonLimeFaint).To(g.Separator()),
		g.Dummy(0, 8),
		g.Align(g.AlignCenter).To(
			g.Style().SetFontSize(16).SetColor(g.StyleColorText, ThemeTextDim).To(
				g.Label("Nadcord  •  based on Vencord Installer  •  GPL-3.0"),
			),
		),

		InfoModal("#patched", "Installed!", "Nadcord was successfully installed! Restart Discord."),
		InfoModal("#unpatched", "Uninstalled", "Nadcord has been uninstalled!"),
		InfoModal("#scuffed-install", "Hold On!", "You have a broken Discord Install.\n"+
			"Sometimes Discord decides to install to the wrong location for some reason!\n"+
			"You need to fix this before patching, otherwise Nadcord will likely not work.\n\n"+
			"Use the below button to jump there and delete any folder called Discord or Squirrel.\n"+
			"If the folder is now empty, feel free to go back a step and delete that folder too.\n"+
			"Then see if Discord still starts. If not, reinstall it"),
		RawInfoModal("#openasar-confirm", "OpenAsar", "OpenAsar is an open-source alternative of Discord desktop's app.asar.\n"+
			"Nadcord is in no way affiliated with OpenAsar.\n"+
			"You're installing OpenAsar at your own risk. If you run into issues with OpenAsar,\n"+
			"no support will be provided, join the OpenAsar Server instead!\n\n"+
			"To install OpenAsar, press Accept and click 'Install OpenAsar' again.", "", true),
		InfoModal("#insufficient-permissions", "Insufficient Permissions", "Permission denied. Please grant the installer permissions in the settings."),
		InfoModal("#openasar-patched", "Successfully Installed OpenAsar", "If Discord is still open, fully close it first. Then start it again and verify OpenAsar installed successfully!"),
		InfoModal("#openasar-unpatched", "Successfully Uninstalled OpenAsar", "If Discord is still open, fully close it first. Then start it again and it should be back to stock!"),
		InfoModal("#invalid-custom-location", "Invalid Location", "The specified location is not a valid Discord install.\nMake sure you select the base folder.\n\nHint: Discord snap is not supported. use flatpak or .deb"),
		InfoModalExtra("#modal"+strconv.Itoa(modalId), modalTitle, modalMessage, modalExtra),
		ErrorModal("#error-modal"+strconv.Itoa(modalId), modalTitle, modalExtra, modalSupportLink),
		UpdateModal(),
	}

	return layout
}

func renderErrorCard(bgColor color.Color, textColor color.Color, message string, height float32) g.Widget {
	return g.Style().
		SetColor(g.StyleColorChildBg, bgColor).
		SetStyleFloat(g.StyleVarAlpha, 0.9).
		SetStyle(g.StyleVarWindowPadding, 10, 10).
		SetStyleFloat(g.StyleVarChildRounding, 12).
		To(
			g.Child().
				Size(g.Auto, height).
				Layout(
					g.Row(
						g.Style().SetColor(g.StyleColorText, textColor).To(
							g.Markdown(&message),
						),
					),
				),
		)
}

func loop() {
	g.PushWindowPadding(48, 48)
	g.PushColorWindowBg(ThemeBg)

	g.SingleWindow().
		RegisterKeyboardShortcuts(
			g.WindowShortcut{Key: g.KeyUp, Callback: func() {
				if radioIdx > 0 {
					radioIdx--
				}
			}},
			g.WindowShortcut{Key: g.KeyDown, Callback: func() {
				if radioIdx < customChoiceIdx {
					radioIdx++
				}
			}},
		).
		Layout(g.Style().
			SetColor(g.StyleColorText, ThemeText).
			SetColor(g.StyleColorTextDisabled, ThemeTextDim).
			SetColor(g.StyleColorPopupBg, ThemeSurface).
			SetColor(g.StyleColorChildBg, ThemeSurface).
			SetColor(g.StyleColorBorder, NeonLimeDim).
			SetColor(g.StyleColorFrameBg, ThemeSurface).
			SetColor(g.StyleColorFrameBgHovered, ThemeSurfaceHover).
			SetColor(g.StyleColorFrameBgActive, ThemeSurfaceHover).
			SetColor(g.StyleColorCheckMark, NeonLime).
			SetColor(g.StyleColorButton, ThemeSurfaceLime).
			SetColor(g.StyleColorButtonHovered, ThemeSurfaceHover).
			SetColor(g.StyleColorButtonActive, NeonLimeActive).
			SetColor(g.StyleColorHeader, NeonLimeFaint).
			SetColor(g.StyleColorHeaderHovered, NeonLimeDim).
			SetColor(g.StyleColorHeaderActive, NeonLimeDim).
			SetColor(g.StyleColorTextSelectedBg, NeonLimeDim).
			SetColor(g.StyleColorScrollbarBg, ThemeBg).
			SetColor(g.StyleColorScrollbarGrab, ThemeSurfaceHover).
			SetColor(g.StyleColorScrollbarGrabHovered, NeonLimeDim).
			SetColor(g.StyleColorScrollbarGrabActive, NeonLime).
			SetColor(g.StyleColorSeparator, NeonLimeDim).
			SetColor(g.StyleColorModalWindowDimBg, color.RGBA{0, 0, 0, 0xb0}).
			SetStyleFloat(g.StyleVarFrameRounding, 10).
			SetStyleFloat(g.StyleVarPopupRounding, 12).
			SetStyleFloat(g.StyleVarPopupBorderSize, 1).
			SetStyleFloat(g.StyleVarScrollbarRounding, 8).
			SetStyleFloat(g.StyleVarGrabRounding, 8).
			SetStyle(g.StyleVarFramePadding, 10, 8).
			To(
				g.Align(g.AlignCenter).To(
					g.Style().SetFontSize(56).SetColor(g.StyleColorText, NeonLime).To(
						g.Label("NADCORD"),
					),
				),
				g.Align(g.AlignCenter).To(
					g.Style().SetFontSize(20).SetColor(g.StyleColorText, ThemeTextDim).To(
						g.Label("I N S T A L L E R"),
					),
				),
				g.Dummy(0, 14),
				g.Style().SetColor(g.StyleColorSeparator, NeonLime).To(g.Separator()),
				g.Dummy(0, 30),

				&CondWidget{
					GithubError != nil,
					func() g.Widget {
						return g.Style().SetFontSize(20).To(renderErrorCard(DiscordRed, color.White, "Failed to fetch Info from GitHub. Check your connection and try again.", 40))
					},
					nil,
				},

				&CondWidget{
					predicate:  FilesDirErr != nil,
					ifWidget:   renderFilesDirErr,
					elseWidget: renderInstaller,
				},
			))

	g.PopStyleColor()
	g.PopStyle()
}

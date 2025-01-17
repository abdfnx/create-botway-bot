package templates

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/abdfnx/botway/constants"
	"github.com/abdfnx/resto/core/api"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func BotSecrets(platform string) string {
	if strings.Contains(platform, "discord") {
		return "DISCORD_TOKEN DISCORD_CLIENT_ID\n# You can add guild ids of your servers by adding ARG SERVER_NAME_GUILD_ID"
	} else if strings.Contains(platform, "telegram") {
		return "TELEGRAM_TOKEN"
	} else if strings.Contains(platform, "slack") {
		return "SLACK_TOKEN SLACK_APP_TOKEN SLACK_SIGNING_SECRET"
	} else if strings.Contains(platform, "twitch") {
		return "TWITCH_OAUTH_TOKEN TWITCH_CLIENT_ID TWITCH CLIENT_SECRET"
	}

	return "" + platform
}

func Content(arg, templateName, botName, platform string) string {
	org := "botwayorg"

	if templateName == "botway" {
		org = "abdfnx"
	}

	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s", org, templateName, arg)
	respone, status, _, err := api.BasicGet(url, "GET", "", "", "", "", false, 0, nil)

	if err != nil {
		fmt.Println(err.Error())
	}

	if status == "404" || status == "401" || strings.Contains(respone, "404") {
		fmt.Println("404: " + url)
		os.Exit(0)
	}

	// TODO: fix botway c++ telegram template
	if strings.Contains(respone, "#include <{{.BotName}}/{{.BotName}}.h>") && strings.Contains(platform, "telegram") {
		respone = strings.ReplaceAll(respone, "#include <{{.BotName}}/{{.BotName}}.h>", "")
	} else if strings.Contains(respone, `#include "botway/botway.hpp"`) && strings.Contains(platform, "telegram") {
		respone = strings.ReplaceAll(respone, `#include "botway/botway.hpp"`, `#include "botway.hpp"`)
	} else if strings.Contains(arg, "pubspec.yaml") {
		respone = strings.ReplaceAll(respone, "{{.BotName}}", strings.ReplaceAll(botName, "-", ""))
	}

	respone = strings.ReplaceAll(respone, "{{.BotName}}", botName)

	author := gjson.Get(string(constants.BotwayConfig), "github.username").String()

	if author == "" {
		author = "botway"
	}

	respone = strings.ReplaceAll(respone, "{{.Author}}", author)

	respone = strings.ReplaceAll(respone, "{{.BotSecrets}}", BotSecrets(platform))

	return respone
}

func CheckProject(botName, botType string) {
	if _, err := os.Stat(botName); !os.IsNotExist(err) {
		fmt.Print(constants.SUCCESS_BACKGROUND.Render("SUCCESS"))
		fmt.Println(constants.SUCCESS_FOREGROUND.Render(" " + botName + " Created successfully 🎉"))
	}
}

func startCmd(botName, lang, pm string) string {
	nodeCmd := pm + " start"
	denoCmd := "deno task run"

	if lang == "python" && pm == "pip" {
		if runtime.GOOS == "windows" {
			return `py .\src\main.py`
		} else {
			return `python3 ./src/main.py`
		}
	} else if lang == "python" && pm == "pipenv" {
		if runtime.GOOS == "windows" {
			return `pipenv run py .\src\main.py`
		} else {
			return `pipenv run python3 ./src/main.py`
		}
	} else if lang == "python" && pm == "poetry" {
		if runtime.GOOS == "windows" {
			return `poetry run .\src\main.py`
		} else {
			return `poetry run ./src/main.py`
		}
	} else if lang == "go" {
		return "go run src/main.go"
	} else if lang == "nodejs" || lang == "typescript" {
		return nodeCmd
	} else if lang == "deno" {
		return denoCmd
	} else if lang == "ruby" {
		return "bundle exec ruby src/main.rb"
	} else if lang == "java" || lang == "kotlin" {
		if runtime.GOOS == "windows" {
			return `.\gradlew.bat run`
		} else {
			return "./gradlew run"
		}
	} else if lang == "csharp" {
		return "dotnet run"
	} else if lang == "dart" {
		return "dart run src/main.dart"
	} else if lang == "php" {
		return "php src/main.php"
	} else if lang == "cpp" {
		if runtime.GOOS == "windows" {
			return `.\run.ps1`
		} else {
			return "cd build; make -j; ./" + botName
		}
	} else if lang == "swift" {
		return "swift run"
	} else if lang == "c" {
		if runtime.GOOS == "windows" {
			return `.\run.ps1`
		} else {
			return "gcc src/main.c -o bot -pthread -ldiscord -lcurl; ./bot"
		}
	} else if lang == "nim" {
		return "nim c -r src/main.nim"
	} else if lang == "crystal" {
		return "crystal run src/main.cr"
	} else if lang == "rust" {
		return "cargo run src/main.rs"
	}

	return "# Write your start command here"
}

func csharpGitIgnore() string {
	return `## Ignore Visual Studio temporary files, build results, and
## files generated by popular Visual Studio add-ons.

# User-specific files
*.suo
*.user
*.userosscache
*.sln.docstates

# User-specific files (MonoDevelop/Xamarin Studio)
*.userprefs

# Build results
[Dd]ebug/
[Dd]ebugPublic/
[Rr]elease/
[Rr]eleases/
build/
bld/
[Bb]in/
[Oo]bj/

# Visual Studo 2015 cache/options directory
.vs/

# MSTest test Results
[Tt]est[Rr]esult*/
[Bb]uild[Ll]og.*

# NUNIT
*.VisualState.xml
TestResult.xml

# Build Results of an ATL Project
[Dd]ebugPS/
[Rr]eleasePS/
dlldata.c

*_i.c
*_p.c
*_i.h
*.ilk
*.meta
*.obj
*.pch
*.pdb
*.pgc
*.pgd
*.rsp
*.sbr
*.tlb
*.tli
*.tlh
*.tmp
*.tmp_proj
*.log
*.vspscc
*.vssscc
.builds
*.pidb
*.svclog
*.scc

# Chutzpah Test files
_Chutzpah*

# Visual C++ cache files
ipch/
*.aps
*.ncb
*.opensdf
*.sdf
*.cachefile

# Visual Studio profiler
*.psess
*.vsp
*.vspx

# TFS 2012 Local Workspace
$tf/

# Guidance Automation Toolkit
*.gpState

# ReSharper is a .NET coding add-in
_ReSharper*/
*.[Rr]e[Ss]harper
*.DotSettings.user

# JustCode is a .NET coding addin-in
.JustCode

# TeamCity is a build add-in
_TeamCity*

# DotCover is a Code Coverage Tool
*.dotCover

# NCrunch
_NCrunch_*
.*crunch*.local.xml

# MightyMoose
*.mm.*
AutoTest.Net/

# Web workbench (sass)
.sass-cache/

# Installshield output folder
[Ee]xpress/

# DocProject is a documentation generator add-in
DocProject/buildhelp/
DocProject/Help/*.HxT
DocProject/Help/*.HxC
DocProject/Help/*.hhc
DocProject/Help/*.hhk
DocProject/Help/*.hhp
DocProject/Help/Html2
DocProject/Help/html

# Click-Once directory
publish/

# Publish Web Output
*.[Pp]ublish.xml
*.azurePubxml
# TODO: Comment the next line if you want to checkin your web deploy settings
# but database connection strings (with potential passwords) will be unencrypted
*.pubxml
*.publishproj

# NuGet Packages
*.nupkg
# The packages folder can be ignored because of Package Restore
**/packages/*
# except build/, which is used as an MSBuild target.
!**/packages/build/
# Uncomment if necessary however generally it will be regenerated when needed
#!**/packages/repositories.config

# Windows Azure Build Output
csx/
*.build.csdef

# Windows Store app package directory
AppPackages/

# Others
*.[Cc]ache
ClientBin/
~$*
*~
*.dbmdl
*.dbproj.schemaview
*.pfx
*.publishsettings
node_modules/
bower_components/

# RIA/Silverlight projects
Generated_Code/

# Backup & report files from converting an old project file
# to a newer Visual Studio version. Backup files are not needed,
# because we have git ;-)
_UpgradeReport_Files/
Backup*/
UpgradeLog*.XML
UpgradeLog*.htm

# SQL Server files
*.mdf
*.ldf

# Business Intelligence projects
*.rdl.data
*.bim.layout
*.bim_*.settings

# Microsoft Fakes
FakesAssemblies/

# Node.js Tools for Visual Studio
.ntvs_analysis.dat

# Visual Studio 6 build log
*.plg

# Visual Studio 6 workspace options file
*.opt

#Custom
project.lock.json
*.pyc
/.editorconfig

\.idea/

# Codealike UID
codealike.json`
}

func CreateBot(botName, platform, lang, pm string) {
	hostService := "zeabur.com"

	if err := os.Mkdir(botName, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir(filepath.Join(botName, "src"), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if lang == "Swift" || lang == "Java" {
		os.RemoveAll(filepath.Join(botName, "src"))
	}

	if err := os.Mkdir(filepath.Join(botName, "config"), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	botConfig := viper.New()

	dockerImage := "botway-local/" + botName

	botConfig.AddConfigPath(botName)
	botConfig.SetConfigName(".botway")
	botConfig.SetConfigType("yaml")

	botConfig.SetDefault("bot.lang", lang)
	botConfig.SetDefault("bot.name", botName)
	botConfig.SetDefault("bot.host_service", hostService)

	if pm != "continue" {
		botConfig.SetDefault("bot.package_manager", pm)
	}

	botConfig.SetDefault("bot.type", platform)
	botConfig.SetDefault("bot.start_cmd", startCmd(botName, lang, pm))

	botConfig.SetDefault("docker.image", dockerImage)
	botConfig.SetDefault("docker.enable_buildkit", true)
	botConfig.SetDefault("docker.cmds.build", "docker build -t "+dockerImage+" .")
	botConfig.SetDefault("docker.cmds.run", "docker run -it "+dockerImage)

	if platform == "discord" {
		guildsFile := os.WriteFile(filepath.Join(botName, "config", "guilds.json"), []byte("{}"), 0644)

		if guildsFile != nil {
			panic(guildsFile)
		}
	}

	if err := botConfig.SafeWriteConfig(); err != nil {
		if os.IsNotExist(err) {
			err = botConfig.WriteConfig()

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := botConfig.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal(err)
		}
	}

	dotGitIgnoreFileContent := ""

	l := ""

	if lang == "python" {
		l = "Python"
	} else if lang == "nodejs" || lang == "typescript" || lang == "deno" {
		l = "Node"
	} else if lang == "go" {
		l = "Go"
	} else if lang == "ruby" {
		l = "Ruby"
	} else if lang == "java" || lang == "kotlin" {
		l = "Java"
	} else if lang == "rust" {
		l = "Rust"
	} else if lang == "csharp" {
		l = "C#"
	} else if lang == "dart" {
		l = "Dart"
	} else if lang == "php" {
		l = "PHP"
	} else if lang == "cpp" {
		l = "C++"
	} else if lang == "nim" {
		l = "Nim"
	} else if lang == "swift" {
		l = "Swift"
	} else if lang == "c" {
		l = "C"
	} else if lang == "crytal" {
		l = "Crystal"
	}

	respone, status, _, err := api.BasicGet("https://raw.githubusercontent.com/github/gitignore/main/"+l+".gitignore", "GET", "", "", "", "", true, 0, nil)

	if err != nil {
		fmt.Println(err.Error())
	}

	if status == "404" || status == "401" {
		fmt.Println("404")
		os.Exit(0)
	}

	if lang == "deno" {
		respone += "fileloader.ts"
	}

	if lang == "csharp" {
		respone = csharpGitIgnore()
	}

	if lang == "java" || lang == "kotlin" {
		respone += "\n.gradle\nbuild"
	}

	if lang == "swift" {
		respone += "\n.build\nPackage.resolved"
	}

	if lang == "crystal" {
		respone = `/docs/
	/lib/
	/bin/
	/.shards/
	*.dwarf
	
	# Libraries don't need dependency lock
	# Dependencies will be locked in applications that use them
	/shard.lock`
	}

	dotGitIgnoreFileContent = respone + "\n*.lock\nbotway-tokens.env\nbotway.json"

	dotGitIgnoreFile := os.WriteFile(filepath.Join(botName, ".gitignore"), []byte(dotGitIgnoreFileContent), 0644)
	dotGitKeepFile := os.WriteFile(filepath.Join(botName, "config", ".gitkeep"), []byte(""), 0644)
	readmeFile := os.WriteFile(filepath.Join(botName, "README.md"), []byte(Content("bot-readme.md", "resources", "", "")), 0644)
	dockerComposeFile := os.WriteFile(filepath.Join(botName, "docker-compose.yaml"), []byte(Content("dockerfiles/compose/docker-compose.yaml", "botway", "", "")), 0644)

	if dotGitIgnoreFile != nil {
		log.Fatal(dotGitIgnoreFile)
	}

	if readmeFile != nil {
		log.Fatal(readmeFile)
	}

	if dotGitKeepFile != nil {
		log.Fatal(dotGitKeepFile)
	}

	if dockerComposeFile != nil {
		log.Fatal(dockerComposeFile)
	}
}

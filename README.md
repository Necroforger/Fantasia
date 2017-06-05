

# Fantasia
<!-- TOC -->

- [Fantasia](#fantasia)
- [Installing](#installing)
- [Running](#running)
- [Flags](#flags)
- [Config samples](#config-samples)
    - [Selfbot](#selfbot)
- [Modules](#modules)

<!-- /TOC -->

______________

![Example help menu](http://i.imgur.com/jC4lFxE.png)

# Installing
Go to the [releases](https://github.com/Necroforger/Fantasia/releases) and download a version or.

`go get -U github.com/Necroforger`

Navigate to GOPATH/github.com/Necroforger/Fantasia and use `go build` to create an executeable.

# Running
Execute the bot and it should generate a `config.toml` file. Fill this in with your bot information and execute the bot again. You can copy the sample config to get started quickly.

# Flags

Executing with flags is optional unless you want to use the same config
With multiple tokens, or use a config file stored in a path other than `./config.toml`

| Flag | Description           |
|------|-----------------------|
| t    | Specify the bot token |
| c    | Config file path      |
| s    | Enable selfbot mode   |
# Config samples

## Selfbot
```toml
# Bot tokens must be prefixed by 'bot'
Token = "REPLACE_WITH_ACCOUNT_TOKEN"

# Disabled commands is a list of commands you want to be disabled when the
# Bot starts. Yes, you can disable the help command.
DisabledCommands = []


# Setting selfbot to true prevents the bot from responding to users that are not you
[System]
  Prefix = "./"
  Selfbot = true 
  YoutubeDLPath = "youtube-dl"

# Any modules set to false are disabled.
# If inverted is true, all modules are enabled by default
# And any modules set to true are disabled.
[Modules]
  Inverted = false
  Audio = true
  Eval = true
  General = true
  Images = true
  Information = true
  Musicplayer = true
  Roles = true

  # Custom audio commands
  [Modules.AudioConfig]

    # Soundclip commands are in the format of ["Command Name", "Description", "url", "url"...]
    # If more than one URL is present, the command will choose a random one from the list.
    # If the path is not prefixed by http:// or https:// it will attempt to get the clip from the file system.
    # If the path is a folder, it will get a random clip from the folder.
    SoundclipCommands = [
      ["granddad", "GRAND DAD", "https://youtu.be/gIcmIjfdE78"],
    ]

  # Custom image commands
  [Modules.ImagesConfig]

    # Booru commands allow you to request an image or list of images from a booru.
    # usage: boorucommand "list of tags". Supports all boorus supported by
    # https://github.com/Necroforger/boorudl
    BooruCommands = [
      ["danbooru",  "http://danbooru.donmai.us"],
      ["safebooru", "https://safebooru.org/"],
      ["googleimg", "http://google.com"],
    ]

    # Image commands are in the format of ["Command Name", "Description", "url", "url"...]
    # If more than one URL is present, the command will choose a random one from the list.
    # If the path is not prefixed by http:// or https:// it will attempt to get the image from the file system.
    # If the path is a folder, it will get a random image from the folder.
    ImageCommands = [
        ["cirno",   "cir-no",
        "https://nokywln.files.wordpress.com/2011/09/807720-20920920cirno20touhou20e291a81.jpg?w=500"],

        ["cirnopi", "Cirno calculates the exact value of pi", 
        "https://lh6.googleusercontent.com/-7kCspKNcZvU/VQRWMo4fb7I/AAAAAAAABIg/fwBfrgrCcx0/w800-h800/cirno_PI.jpg"],
    ]

# Ffmpeg must be in your path and DCA-RS must be in the same folder
# As your bot if you intend to use audio
# You can change the paths here if you want to.
[Dream]
  FfmpegPath = "ffmpeg"
  DcaRsPath = "./dca-rs"
```


# Modules

| Module      | Description                                                                                 |
|-------------|---------------------------------------------------------------------------------------------|
| Audio       | Simple youtube audio player                                                                 |
| General     | General bot commands                                                                        |
| Information | Gets information about the bot. Dynamically generates a help menu based on enabled commands |
| Roles       | Role managing module                                                                        |
| eval        | Module for evaluating code                                                                  |
| Images      | Various image commands                                                                      |
| Musicplayer | Do not use, work in progress                                                                |
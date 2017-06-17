

# Fantasia
<!-- TOC -->

- [Fantasia](#fantasia)
- [Dependencies](#dependencies)
    - [Audio dependencies](#audio-dependencies)
- [Installing](#installing)
- [Running](#running)
- [Flags](#flags)
- [Config Example](#config-example)
- [Modules](#modules)

<!-- /TOC -->

______________
# Dependencies
## Audio dependencies
* [dca-rs](https://github.com/nstafie/dca-rs/releases)
* [ffmpeg](https://ffmpeg.org/)
* [youtubedl](https://rg3.github.io/youtube-dl/) - optional

You need to have [dca-rs](https://github.com/nstafie/dca-rs/releases) in the same directory as your bot executeable in order to use audio playing commands. This is required to convert the audio to opus format.

[ffmpeg](https://ffmpeg.org/) is also required for converting audio to opus format. It should be [installed to your path.](http://www.wikihow.com/Install-FFmpeg-on-Windows)

Installing [youtubedl](https://rg3.github.io/youtube-dl/) will allow you to queue videos in the media player from a variety of sources, such as soundcloud and facebook rather than specifically youtube. [It should be added to your path similarly to ffmpeg.](http://www.wikihow.com/Install-YouTube-DL.py-to-Download-YouTube-Videos-to-Your-PC). If you choose not to use youtube-dl, set UseYotubeDL in the MusicPlayer config to false and it will resort to using the golang downloader, [rylio/ytdl](https://github.com/rylio/ytdl)



# Installing
Go to the [releases](https://github.com/Necroforger/Fantasia/releases) and download a version or.

`go get -U github.com/Necroforger/Fantasia`

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
| p    | Bot prefix            |


# Config Example

```toml
# All bot tokens are prefixed by 'Bot '. Ex 'Bot MsasmdJAsndjANsjh23'
Token = "Bot MsasmdJAsndjANsjh23"
DisabledCommands = []

# If selfbot is enabled, your bot must be run on a user token and it will only respond
# To itself. Certain modules may not work on user accounts.
[System]
  Prefix = ">"
  Selfbot = false

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

    # Category name for sound clips
    SoundClipCommandsCategory = "Sound clips"

    # Soundclip commands are in the format of ["Command Name", "Description", "url", "url"...]
    # If more than one URL is present, the command will choose a random one from the list.>
    # If the path is not prefixed by http:// or https:// it will attempt to get the clip from the file system.
    # If the path is a folder, it will get a random clip from the folder.
    SoundclipCommands = [
      ["granddad",   "[source](https://youtu.be/gIcmIjfdE78)", "https://youtu.be/gIcmIjfdE78"]
    ]

  [Modules.ImagesConfig]

    # Include image filtering commands
    FilterCommands = true

    # Controls the category of image filters
    ImageFiltersCategory = "Image filters"

    # Leaving this setting blank will set custom image commands category to the default
    # Module category, which would be Images.
    ImageCommandsCategory = ""

    # Image commands are in the format of ["Command Name", "Description", "url", "url"...]
    # If more than one URL is present, the command will choose a random one from the list.
    # If the path is not prefixed by http:// or https:// it will attempt to get the image from the file system.
    # If the path is a folder, it will get a random image from the folder.
    ImageCommands = [
      ["cirnopi", "cirno calculates the exact value of pi", 
          "https://lh6.googleusercontent.com/-7kCspKNcZvU/VQRWMo4fb7I/AAAAAAAABIg/fwBfrgrCcx0/w800-h800/cirno_PI.jpg"],
      ["highground", "its over anakin, I have the high ground",
          'https://cdn.discordapp.com/attachments/221341345539686400/321496580362338304/Icantevendrawastickfiguresoidont_25dec9985c1399cf20e3bd064a7a8571.jpg'],
      ["nothing_to_do_with_this", "I don't think I want anything to do with this",
          "https://cdn.discordapp.com/attachments/319171684105977857/321117932002476033/FAqSXDJ.png"],
    ]


    # BooruCommandsCategory allows you to change the category the booru commands
    # Appear under. The help menu will display it under a different field.
    BooruCommandsCategory = "Boorus"

    # Booru commands allow you to request an image or list of images from a booru.
    # usage: boorucommand "list of tags". Supports all boorus supported by
    # https://github.com/Necroforger/boorudl
    BooruCommands = [
        ["danbooru",  "http://danbooru.donmai.us"],
        ["safebooru", "https://safebooru.org/"],
        ["img",       "http://google.com"],
    ]

[Modules.MusicplayerConfig]
    # Use the music player subrouter.
    # This can prevent conflicts between commands with
    # The same name. Requires you to type `musicplayer` or `m`
    # Before every musicplayer command.
    UseSubrouter = true

    # Debug adds a testing song queue
    Debug = false

    # Use youtube-dl to queue and stream media
    UseYoutubeDL = true


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
| Musicplayer | Queue and play songs from a variety of sources                                                           |
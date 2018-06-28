

# Fantasia
<!-- TOC -->

- [Fantasia](#fantasia)
- [Dependencies](#dependencies)
	- [Audio dependencies](#audio-dependencies)
- [Installing](#installing)
- [Running](#running)
- [Flags](#flags)

<!-- /TOC -->

______________
# Dependencies
## Audio dependencies
* [ffmpeg](https://ffmpeg.org/)
* [youtubedl](https://rg3.github.io/youtube-dl/) - optional

[ffmpeg](https://ffmpeg.org/) is required for converting audio to opus format. It should be [installed to your path.](http://www.wikihow.com/Install-FFmpeg-on-Windows)

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



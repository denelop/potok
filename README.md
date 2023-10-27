# Potok

The video streaming solution converting live webcam streams to HLS.

## Requirements

- [Go](https://go.dev/) (for development only)
- [ffmpeg](https://www.ffmpeg.org/)
  - Install latest release for macOS
    ```sh
    brew install ffmpeg
    ```
  - Install latest release for Linux kernels >= v3.2.0
    ```sh
    wget https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz
    tar xvf ffmpeg-release-amd64-static.tar.xz
    sudo mv ffmpeg-*-amd64-static/ff* /usr/local/bin/
    ```

## Web Player

Use [HLS.js](https://github.com/video-dev/hls.js).

## Meaning of "potok"

[The literal transalation of "stream" in Bosnian/Croatian/Serbian language is "potok".](https://translate.google.com/?sl=en&tl=bs&text=stream&op=translate) It is actually a river stream.

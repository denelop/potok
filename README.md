<div align="center">
  <br />

![potok](logo.png)

  <p>Video transcoding solution converting live webcam streams to HLS (<a href="https://developer.apple.com/streaming">HTTP Live Streaming</a>) that browsers and players can consume. It also supports often needed requirements like watermarking and scaling the video streams.</p>
</div>

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

## Defining Streams

Streams are defined in a single YAML file following the format:

```yaml
- name: <Stream Name> (required)
  url: <URL of the IP camera> (required)
  rtsp_transport: <"tcp" or "udp"> (optional, "tcp" is default)
  scale: <width:height> (optional)
  watermarks:
    - file: <path to the image> (required)
      opacity: <0.1 - 1> (required)
      position: <"top-left" or "center" or "bottom-right"> (required)
```

Definitions will be transcoded and the main streaming playlist (video) will each be served over HTTP at:

```
http://localhost:53030/streaming/<Stream Name>/playlist.m3u8
```

## Web Player

Since only Apple products natively support HLS, you're recommended to use [HLS.js](https://github.com/video-dev/hls.js) as an open-source browser player.

## Meaning of "potok"

[The literal transalation of "stream" in Bosnian/Croatian/Serbian language is "potok".](https://translate.google.com/?sl=en&tl=bs&text=stream&op=translate) It is actually a river stream.

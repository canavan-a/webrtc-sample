# webrtc-sample

sinple WebRTC implementation

## Linux command

```
ffmpeg \
  -f v4l2 -i /dev/video0 \  # Camera input
  -f alsa -i hw:1 \  # Microphone input (adjust based on your system setup)
  -c:v vp8 \  # VP8 codec for video
  -c:a libopus \  # Opus codec for audio
  -b:a 128k \  # Audio bitrate (adjust as needed)
  -b:v 1M \  # Video bitrate (adjust as needed)
  -f webm \  # Output format WebM for WebCodecs compatibility
  -g 30 \  # Keyframe interval for video (adjust as needed)
  -r 30 \  # Frame rate (adjust as needed)
  -async 1 \  # Ensure audio and video sync
  -vsync 2 \  # Video sync mode
  -f webm - | websocat ws://yourwebsocketserver  # Pipe to WebSocket
```

## windows command

```
ffmpeg ^
  -f dshow -i video="Aidan's S23 (Windows Virtual Camera)"^
  -f dshow -i audio="Your Audio Device Name"^  # Microphone input (adjust for your device)
  -c:v vp8^  # VP8 codec for video
  -c:a libopus^  # Opus codec for audio
  -b:a 128k^  # Audio bitrate (adjust as needed)
  -b:v 1M^  # Video bitrate (adjust as needed)
  -f webm^  # Output format WebM for WebCodecs compatibility
  -g 30^  # Keyframe interval for video (adjust as needed)
  -r 30^  # Frame rate (adjust as needed)
  -async 1^  # Ensure audio and video sync
  -vsync 2^  # Video sync mode
  -f webm - | websocketd --ws-url=ws://localhost:5000/relay
```

# webrtc-sample

sinple WebRTC implementation

`ffmpeg -list_devices true -f dshow -i dummy`

# Start RTC server

## Linux command

```

```

## windows command

```
ffmpeg -f dshow -i video="FHD Camera" -f dshow -i audio="Microphone Array (Intel® Smart Sound Technology for Digital Microphones)" -c:v vp8 -pix_fmt yuv420p -crf 28 -c:a libopus -s 640x360 -r 15 -rtbufsize 1G -f rtp_mpegts rtp://127.0.0.1:5004

high speed and performance:

ffmpeg -f dshow -i video="FHD Camera" -f dshow -i audio="Microphone Array (Intel® Smart Sound Technology for Digital Microphones)" ^
  -c:v vp8 -pix_fmt yuv420p -crf 23 -preset ultrafast -c:a libopus -s 640x360 -r 15 ^
  -rtbufsize 2G -f rtp_mpegts rtp://127.0.0.1:5005

```

## test reciving the data

```
ffmpeg -loglevel debug -i rtp://127.0.0.1:5005?localaddr=127.0.0.1 -c:v vp8 -c:a opus -f wmv output.wmv



```

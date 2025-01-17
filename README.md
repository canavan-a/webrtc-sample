# webrtc-sample

sinple WebRTC implementation

`ffmpeg -list_devices true -f dshow -i dummy`

# Start RTC server

## Linux command

```

```

## windows command

```


# mpeg ts data

ffmpeg -f dshow -framerate 30 -i video="FHD Camera":audio="Microphone Array (Intel® Smart Sound Technology for Digital Microphones)" -vcodec libx264 -pix_fmt yuv420p -crf 23 -preset ultrafast -s 640x360 -rtbufsize 1M -f rtp_mpegts -acodec libopus -colorspace bt709 -color_range tv rtp://127.0.0.1:5005

# opus and h264 data

video + sdp_file
ffmpeg -f dshow -framerate 30 -i video="FHD Camera" -vcodec libx264 -pix_fmt yuv420p -preset ultrafast -tune zerolatency -crf 23 -x264-params "keyint=30:min-keyint=30:scenecut=0" -s 640x360 -rtbufsize 64M -f rtp -sdp_file video.sdp rtp://127.0.0.1:5005


audio + sdp_file
ffmpeg -f dshow -i audio="Microphone Array (Intel® Smart Sound Technology for Digital Microphones)" -acodec libopus -b:a 64k -ar 48000 -ac 2 -application lowdelay -frame_duration 10 -rtbufsize 64M -preset ultrafast -tune zerolatency -f rtp -sdp_file audio.sdp rtp://127.0.0.1:5006




```

## test reciving the data

```
# mpeg_ts input

ffmpeg -i rtp://127.0.0.1:5005 -c:v wmv2 -c:a wmav2 -b:v 1024k -b:a 192k -flush_packets 0 output.wmv

audio only:
ffmpeg -i rtp://127.0.0.1:5005 -map 0:a -c:a wmav2 -b:a 192k -flush_packets 0 output.wmv

video only:
ffmpeg -i rtp://127.0.0.1:5005 -map 0:v -c:v wmv2 -b:v 1024k -flush_packets 0 output.wmv

# h264 input (need correct sdp file)
ffmpeg -protocol_whitelist file,crypto,data,rtp,udp -i video.sdp -c copy -f mpegts output_video.ts

# opus input (need correct sdp file)
ffmpeg -protocol_whitelist file,crypto,data,rtp,udp -i audio.sdp -c:a copy -f ogg output_audio.ogg

```

import { useEffect, useRef, useState } from "react";

const STUN_SERVERS = [
  "stun:stun.l.google.com:19302",
  "stun:stun1.l.google.com:19302",
  "stun:stun2.l.google.com:19302",
];

export const Home = () => {
  // webRTC logic here
  const audioRef = useRef(null);
  const videoRef = useRef(null);

  const [offer, setOffer] = useState(null);
  const signalingServer = new WebSocket("ws://192.168.1.153:6789/relay");

  signalingServer.onopen = () => {
    console.log("ws open");
  };
  const servers = {
    iceServers: [
      {
        urls: [...STUN_SERVERS],
      },
    ],
    iceTransportPolicy: "all", // Allow both peer-to-peer and relay transport
    bundlePolicy: "max-bundle", // Ensures only one ICE candidate for both audio/video
    rtcpMuxPolicy: "require",
  };
  const pc = new RTCPeerConnection(servers);

  const start = async () => {
    const dataChannel = pc.createDataChannel("dummyChannel");

    signalingServer.onmessage = (event) => {
      const response = JSON.parse(event.data);
      console.log(response);
      if (response.type == "answer") {
        console.log(response);
        pc.setRemoteDescription(response);
      } else if (response.type == "candidate") {
        pc.addIceCandidate(response);
      }
    };

    const offerOptions = {
      offerToReceiveAudio: true,
      offerToReceiveVideo: true,
      voiceActivityDetection: false,
      iceRestart: true,
    };

    const offer = await pc.createOffer(offerOptions);
    pc.setLocalDescription(offer);

    signalingServer.send(JSON.stringify(offer));

    pc.onicecandidate = (event) => {
      if (event.candidate) {
        // Send ICE candidate to signaling server
        const payload = {
          type: "candidate",
          candidate: event.candidate.candidate,
          sdpMid: event.candidate.sdpMid,
          sdpMLineIndex: event.candidate.sdpMLineIndex,
        };
        signalingServer.send(JSON.stringify(payload));
      }
    };

    pc.ontrack = (event) => {
      const stream = event.streams[0]; // Get the first MediaStream
      console.log("Received stream:", stream);

      // Attach the stream to an audio or video element
      if (stream.getVideoTracks().length > 0) {
        videoRef.current.srcObject = stream;
      }

      if (stream.getAudioTracks().length > 0) {
        audioRef.current.srcObject = stream;
      }
    };
  };

  return (
    <>
      <div className="w-full h-screen flex items-center justify-center">
        <p>
          <button className="btn btn-md" onClick={start}>
            test
          </button>
          <button
            className="btn btn-md"
            onClick={() => {
              console.log(pc);
            }}
          >
            pc test
          </button>
          <video
            ref={videoRef}
            autoPlay
            playsInline
            className="w-full h-auto"
          />
          <audio ref={audioRef} autoPlay></audio>
        </p>
      </div>
    </>
  );
};

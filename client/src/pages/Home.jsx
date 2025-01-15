import { useEffect, useState } from "react";

const STUN_SERVERS = [
  "stun:stun.l.google.com:19302",
  "stun:stun1.l.google.com:19302",
  "stun:stun2.l.google.com:19302",
];

export const Home = () => {
  // webRTC logic here

  const [peerConnection, setPeerConnection] = useState(null);
  const [offer, setOffer] = useState(null);
  const [candidate, setCandidate] = useState(null);
  const signalingServer = new WebSocket("ws://localhost:6789/relay");

  const start = async () => {
    console.log("hello world");
    const servers = {
      iceServers: [
        {
          urls: [...STUN_SERVERS],
        },
      ],
    };
    const pc = new RTCPeerConnection(servers);

    const dataChannel = pc.createDataChannel("dummyChannel");

    const offer = await pc.createOffer();
    pc.setLocalDescription(offer);

    pc.onicecandidate = (event) => {
      console.log("onicecandidate event triggered");

      if (event.candidate) {
        console.log("ICE Candidate gathered: ", event.candidate);
        // Send ICE candidate to signaling server
        setCandidate(event.candidate);

        console.log(event.candidate);
      }
    };
  };

  return (
    <>
      <div className="w-full h-screen flex items-center justify-center">
        <p>
          <strong>ICE Candidate:</strong>{" "}
          {candidate ? JSON.stringify(candidate) : "Waiting for candidates..."}
          <button className="btn btn-md" onClick={start}>
            test
          </button>
        </p>
      </div>
    </>
  );
};

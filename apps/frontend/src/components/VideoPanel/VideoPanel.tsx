import React, { useEffect, useRef, useState } from "react";
import Peer, { MediaConnection } from "peerjs";
import "./styles.scss";
import { Button } from "antd";
import {
  ApiOutlined,
  AudioMutedOutlined,
  AudioOutlined,
  PhoneOutlined,
  VideoCameraOutlined,
} from "@ant-design/icons";

const VideoPanel = () => {
  const [currentId, setCurrentId] = useState<string | undefined>();
  const [partnerId, setPartnerId] = useState<string | undefined>();

  const remoteVideoRef = useRef<HTMLVideoElement>(null);
  const currentUserVideoRef = useRef<HTMLVideoElement>(null);
  const [peerInstance, setPeerInstance] = useState<Peer | null>(null);
  const [callInstance, setCallInstance] = useState<MediaConnection | null>(
    null
  );
  const [remoteCallInstance, setRemoteCallInstance] =
    useState<MediaConnection | null>(null);
  const [userStream, setUserStream] = useState<MediaStream | null>(null);
  const [videoOn, setVideoOn] = useState<boolean>(true);
  const [muteOn, setMuteOn] = useState<boolean>(false);
  const [isCalling, setIsCalling] = useState<boolean>(false);

  useEffect(() => {
    const matchId = localStorage.getItem("collabId")?.toString() ?? "";
    const currentUsername = localStorage.getItem("user")?.toString();
    const matchedUsername = localStorage.getItem("matchedUser")?.toString();

    setCurrentId(currentUsername + "-" + (matchId ?? ""));
    setPartnerId(matchedUsername + "-" + (matchId ?? ""));
  }, []);

  const handleCall = () => {
    navigator.mediaDevices
      .getUserMedia({
        video: true,
        audio: true,
      })
      .then((stream) => {
        if (peerInstance && partnerId) {
          const call = peerInstance?.call(partnerId, stream);
          setCallInstance(call);
          setIsCalling(true); // Set isCalling as true since it is the initiator
          if (call) {
            call.on("stream", (userVideoStream: MediaStream) => {
              if (remoteVideoRef.current) {
                remoteVideoRef.current.srcObject = userVideoStream;
              }
            });

            if (userStream) {
              const videoTrack = userStream.getVideoTracks()[0];

              const sender = call.peerConnection.getSenders().find((s) => {
                if (s.track) {
                  return s.track.kind === "video";
                }
              });
              if (sender) {
                sender.replaceTrack(videoTrack); // Replace the video track in the call
              }

              const audioTrack = userStream.getAudioTracks()[0];

              const sender2 = call.peerConnection.getSenders().find((s) => {
                if (s.track) {
                  return s.track.kind === "audio";
                }
              });
              if (sender2) {
                sender2.replaceTrack(audioTrack); // Replace the video track in the call
              }
            }
          }
        }
      });
  };

  useEffect(() => {
    if (currentId) {
      let peer: Peer;
      if (typeof window !== "undefined") {
        peer = new Peer(
          currentId
          //     {
          //   host: "localhost",
          //   port: 4444,
          //   path: "/",
          // }
        );

        setPeerInstance(peer);

        navigator.mediaDevices
          .getUserMedia({
            video: true,
            audio: true,
          })
          .then((stream) => {
            setUserStream(stream);
            if (currentUserVideoRef.current) {
              currentUserVideoRef.current.srcObject = stream;
            }

            peer.on("call", (call) => {
              call.answer(stream);
              setIsCalling(true); // When a peer initiates call, on answer, set the calling state to true
              setRemoteCallInstance(call); // Set remote call instance media connection
              call.on("stream", (userVideoStream) => {
                if (remoteVideoRef.current) {
                  remoteVideoRef.current.srcObject = userVideoStream;
                }
              });
            });
          });
      }
      return () => {
        if (peer) {
          peer.destroy();
        }
      };
    }
  }, [currentId]);

  // When remote peer initiates end call, we set isCalling to false
  useEffect(() => {
    remoteCallInstance?.on("close", () => {
      setIsCalling(false);
    });

    callInstance?.on("close", () => {
      setIsCalling(false);
    });
  }, [remoteCallInstance?.open, callInstance?.open]);

  const toggleVideo = () => {
    if (userStream) {
      const videoTrack = userStream.getVideoTracks()[0];

      if (videoTrack) {
        if (videoOn) {
          // Stop the video track
          videoTrack.enabled = false;

          setVideoOn(false);
          if (callInstance && remoteCallInstance?.open) {
            const sender = (
              callInstance.peerConnection || remoteCallInstance.peerConnection
            )
              .getSenders()
              .find((s) => {
                if (s.track) {
                  return s.track.kind === "video";
                }
              });
            if (sender) {
              sender.replaceTrack(videoTrack); // Replace the video track in the call
            }
          }
        } else {
          // Resume the video track
          videoTrack.enabled = true;
          setVideoOn(true);
        }
      }
    }
  };

  const handleCloseConnection = () => {
    if (callInstance && callInstance.open) {
      callInstance.close();
      setIsCalling(false);
    }

    if (remoteCallInstance && remoteCallInstance.open) {
      remoteCallInstance.close();
      setIsCalling(false);
    }
  };

  const toggleMute = () => {
    if (userStream) {
      const audioTrack = userStream.getAudioTracks()[0];

      if (audioTrack) {
        if (muteOn) {
          // Unmute the audio track
          audioTrack.enabled = true;
          setMuteOn(false);
        } else {
          // Mute the audio track
          audioTrack.enabled = false;
          setMuteOn(true);

          if (callInstance && remoteCallInstance?.open) {
            const sender = (
              callInstance.peerConnection || remoteCallInstance.peerConnection
            )
              .getSenders()
              .find((s) => {
                if (s.track) {
                  return s.track.kind === "audio";
                }
              });
            if (sender) {
              sender.replaceTrack(audioTrack); // Replace the video track in the call
            }
          }
        }
      }
    }
  };

  return (
    <div className="video-panel">
      {/* <p className="header-tag">Video Feed for: {currentUsername}</p> */}
      <video
        className="user-video"
        playsInline
        ref={currentUserVideoRef}
        autoPlay
        muted
      />
      <div className="buttons-container">
        <Button
          onClick={toggleVideo}
          icon={
            videoOn ? (
              <div className="icon-padding">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="1.2em"
                  height="1.2em"
                  viewBox="0 0 32 32"
                >
                  <path
                    fill="currentColor"
                    d="M21 26H4a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h17a2 2 0 0 1 2 2v4.06l5.42-3.87A1 1 0 0 1 30 9v14a1 1 0 0 1-1.58.81L23 19.94V24a2 2 0 0 1-2 2M4 8v16h17v-6a1 1 0 0 1 1.58-.81L28 21.06V10.94l-5.42 3.87A1 1 0 0 1 21 14V8Z"
                  ></path>
                </svg>
              </div>
            ) : (
              <div className="icon-padding">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="1.2em"
                  height="1.2em"
                  viewBox="0 0 32 32"
                >
                  <path
                    fill="currentColor"
                    d="M29.46 8.11a1 1 0 0 0-1 .08L23 12.06v-1.62l7-7L28.56 2L2 28.56L3.44 30l4-4H21a2 2 0 0 0 2-2v-4.06l5.42 3.87A1 1 0 0 0 30 23V9a1 1 0 0 0-.54-.89M28 21.06l-5.42-3.87A1 1 0 0 0 21 18v6H9.44L21 12.44V14a1 1 0 0 0 1.58.81L28 10.94zM4 24V8h16V6H4a2 2 0 0 0-2 2v16z"
                  ></path>
                </svg>
              </div>
            )
          }
          type={videoOn ? "default" : "primary"}
        >
          {videoOn ? "Off Cam" : "On Cam"}
        </Button>
        <Button
          onClick={toggleMute}
          type={!muteOn ? "default" : "primary"}
          icon={muteOn ? <AudioMutedOutlined /> : <AudioOutlined />}
        >
          {muteOn ? "On Mic" : "Off Mic"}
        </Button>
        <Button
          onClick={isCalling ? handleCloseConnection : handleCall}
          icon={
            isCalling ? (
              <div className="icon-padding">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="1.2em"
                  height="1.2em"
                  viewBox="0 0 24 24"
                >
                  <g fill="none">
                    <path
                      stroke="currentColor"
                      strokeLinecap="round"
                      strokeWidth={2}
                      d="m20 4l-4 4m0-4l4 4"
                    ></path>
                    <path
                      fill="currentColor"
                      d="m15.1 15.027l-.543-.516zm.456-.48l.544.517zm2.417-.335l-.374.65zm1.91 1.1l-.374.65zm.539 3.446l.543.517zm-1.42 1.496l-.545-.517zm-1.326.71l.074.745zm-9.86-4.489l.543-.516zm-4.813-9.51l-.749.041zm6.475 1.538l.543.517zm.156-2.81l.613-.433zM8.374 3.91l-.613.433zM5.26 3.609l.544.516zM3.691 5.26l-.543-.516zm7.372 7.795l.544-.517zm4.582 2.488l.455-.48l-1.088-1.033l-.455.48zm1.954-.682l1.91 1.1l.749-1.3l-1.911-1.1zm2.279 3.38l-1.42 1.495l1.087 1.034l1.42-1.496zm-2.275 1.975c-1.435.141-5.18.02-9.244-4.258l-1.087 1.033c4.429 4.663 8.654 4.898 10.478 4.717zm-9.244-4.258c-3.876-4.081-4.526-7.523-4.607-9.033l-1.498.08c.1 1.85.884 5.634 5.018 9.986zm1.376-6.637l.286-.302l-1.087-1.033l-.287.302zm.512-4.062L8.986 3.477l-1.225.866l1.26 1.783zm-5.53-2.168L3.149 4.745l1.088 1.033l1.57-1.653zm4.474 5.713a38 38 0 0 0-.545-.515l-.002.002l-.003.003l-.05.058a1.6 1.6 0 0 0-.23.427c-.098.275-.15.639-.084 1.093c.13.892.715 2.091 2.242 3.7l1.088-1.034c-1.428-1.503-1.78-2.428-1.846-2.884c-.032-.22 0-.335.013-.372l.008-.019l-.028.037l-.018.02s-.002 0-.545-.516m1.328 4.767c1.523 1.604 2.673 2.234 3.55 2.377c.451.073.816.014 1.092-.095a1.5 1.5 0 0 0 .421-.25l.036-.034l.014-.014l.007-.006l.003-.003l.001-.002s.002-.001-.542-.518c-.544-.516-.543-.517-.543-.518l.002-.001l.002-.003l.005-.005l.01-.01l.037-.032q.015-.008-.004.001c-.02.008-.11.04-.3.009c-.402-.066-1.27-.42-2.703-1.929zM8.986 3.477C7.972 2.043 5.944 1.8 4.718 3.092l1.087 1.033c.523-.55 1.444-.507 1.956.218zM3.752 6.926c-.022-.4.152-.8.484-1.148L3.148 4.745c-.536.564-.943 1.347-.894 2.261zm14.705 12.811c-.279.294-.57.452-.854.48l.147 1.492c.747-.073 1.352-.472 1.795-.939zM10.021 9.02c.968-1.019 1.036-2.613.226-3.76l-1.225.866c.422.597.357 1.392-.088 1.86zm9.488 6.942c.821.473.982 1.635.369 2.28l1.087 1.033c1.305-1.374.925-3.673-.707-4.613zm-3.409-.898c.385-.406.986-.497 1.499-.202l.748-1.3c-1.099-.632-2.46-.45-3.335.47z"
                    ></path>
                  </g>
                </svg>
              </div>
            ) : (
              <div className="icon-padding">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="1.2em"
                  height="1.2em"
                  viewBox="0 0 24 24"
                >
                  <g fill="none">
                    <path
                      stroke="currentColor"
                      strokeLinecap="round"
                      strokeWidth={2}
                      d="M13.5 2s2.334.212 5.303 3.182c2.97 2.97 3.182 5.303 3.182 5.303m-7.778-4.949s.99.282 2.475 1.767s1.768 2.475 1.768 2.475"
                    ></path>
                    <path
                      fill="currentColor"
                      d="m15.1 15.027l-.543-.516zm.456-.48l.544.517zm2.417-.335l-.374.65zm1.91 1.1l-.374.65zm.539 3.446l.543.517zm-1.42 1.496l-.545-.517zm-1.326.71l.074.745zm-9.86-4.489l.543-.516zm-4.813-9.51l-.749.041zm6.475 1.538l.543.517zm.156-2.81l.613-.433zM8.374 3.91l-.613.433zM5.26 3.609l.544.516zM3.691 5.26l-.543-.516zm7.372 7.795l.544-.517zm4.582 2.488l.455-.48l-1.088-1.033l-.455.48zm1.954-.682l1.91 1.1l.749-1.3l-1.911-1.1zm2.279 3.38l-1.42 1.495l1.087 1.034l1.42-1.496zm-2.275 1.975c-1.435.141-5.18.02-9.244-4.258l-1.087 1.033c4.429 4.663 8.654 4.898 10.478 4.717zm-9.244-4.258c-3.876-4.081-4.526-7.523-4.607-9.033l-1.498.08c.1 1.85.884 5.634 5.018 9.986zm1.376-6.637l.286-.302l-1.087-1.033l-.287.302zm.512-4.062L8.986 3.477l-1.225.866l1.26 1.783zm-5.53-2.168L3.149 4.745l1.088 1.033l1.57-1.653zm4.474 5.713a38 38 0 0 0-.545-.515l-.002.002l-.003.003l-.05.058a1.6 1.6 0 0 0-.23.427c-.098.275-.15.639-.084 1.093c.13.892.715 2.091 2.242 3.7l1.088-1.034c-1.428-1.503-1.78-2.428-1.846-2.884c-.032-.22 0-.335.013-.372l.008-.019l-.028.037l-.018.02s-.002 0-.545-.516m1.328 4.767c1.523 1.604 2.673 2.234 3.55 2.377c.451.073.816.014 1.092-.095a1.5 1.5 0 0 0 .421-.25l.036-.034l.014-.014l.007-.006l.003-.003l.001-.002s.002-.001-.542-.518c-.544-.516-.543-.517-.543-.518l.002-.001l.002-.003l.005-.005l.01-.01l.037-.032q.015-.008-.004.001c-.02.008-.11.04-.3.009c-.402-.066-1.27-.42-2.703-1.929zM8.986 3.477C7.972 2.043 5.944 1.8 4.718 3.092l1.087 1.033c.523-.55 1.444-.507 1.956.218zM3.752 6.926c-.022-.4.152-.8.484-1.148L3.148 4.745c-.536.564-.943 1.347-.894 2.261zm14.705 12.811c-.279.294-.57.452-.854.48l.147 1.492c.747-.073 1.352-.472 1.795-.939zM10.021 9.02c.968-1.019 1.036-2.613.226-3.76l-1.225.866c.422.597.357 1.392-.088 1.86zm9.488 6.942c.821.473.982 1.635.369 2.28l1.087 1.033c1.305-1.374.925-3.673-.707-4.613zm-3.409-.898c.385-.406.986-.497 1.499-.202l.748-1.3c-1.099-.632-2.46-.45-3.335.47z"
                    ></path>
                  </g>
                </svg>
              </div>
            )
          }
          type="primary"
          danger={isCalling}
        >
          {isCalling ? "End" : "Call"}
        </Button>
      </div>
      {/* <p className="header-tag">Video Feed for: {matchedUsername}</p> */}
      {isCalling && (
        <video
          className="matched-user-video"
          playsInline
          ref={remoteVideoRef}
          autoPlay
        />
      )}
    </div>
  );
};

export default VideoPanel;

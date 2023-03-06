import React, { createContext, useEffect, useRef, useState } from 'react';
import * as JsSIP from 'jssip';
import { UA } from 'jssip';

import {
  EndEvent,
  IncomingAckEvent,
  OutgoingAckEvent,
  RTCSession,
  RTCSessionEventMap,
} from 'jssip/lib/RTCSession';

import {
  IncomingRTCSessionEvent,
  OutgoingRTCSessionEvent,
  UAConfiguration,
} from 'jssip/lib/UA';

import axios from 'axios';

export const WebRTCContext = createContext(null);

export const WebRTCContextProvider = (props: any) => {
  const webSocketUrl = `${process.env.NEXT_PUBLIC_WEBRTC_WEBSOCKET_URL}`;
  const [inCall, setInCall] = useState(false);
  const [callerId, setCallerId] = useState('');
  const [ringing, setRinging] = useState(false);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [userAgent, setUserAgent] = useState<UA>();
  const [session, setSession] = useState<RTCSession>();
  const remoteVideo = useRef<HTMLVideoElement>();
  const [isCallOnHold, setIsCallOnHold] = useState(false);
  const [isCallMuted, setIsCallMuted] = useState(false);

  useEffect(() => {
    const refreshCredentials = () => {
      axios
        .get(
          `/oasis-api/call_credentials?service=sip&username=` +
            'gabi@openline.ai',
        )
        .then((res) => {
          console.error('Got a key: ' + JSON.stringify(res.data));

          setCredentials(res.data.username, res.data.password);
          if (!userAgent) {
            startUA(res.data.username, res.data.password);
          }
          setTimeout(() => {
            refreshCredentials();
          }, (res.data.ttl * 3000) / 4);
        });
    };
    refreshCredentials();
  }, []);
  const setCredentials = (user: string, pass: string) => {
    setUsername(user);
    setPassword(pass);
    userAgent?.set('authorization_user', user);
    userAgent?.set('password', pass);
  };
  const startUA = (username: string, password: string) => {
    if (userAgent) {
      console.log('UA already started! ignoring request');
      return;
    }
    const socket: JsSIP.Socket = new JsSIP.WebSocketInterface(webSocketUrl);
    const configuration: UAConfiguration = {
      sockets: [socket],
      uri: from,
    };

    configuration.authorization_user = username;
    configuration.password = password;
    console.error('Got a configuration: ' + JSON.stringify(configuration));
    JsSIP.debug.enable('JsSIP:*');
    const ua: UA = new JsSIP.UA(configuration);
    ua.on(
      'newRTCSession',
      ({
        originator,
        session: rtcSession,
        request,
      }: IncomingRTCSessionEvent | OutgoingRTCSessionEvent) => {
        if (originator === 'local') return;

        if (inCall) {
          rtcSession.terminate({ status_code: 486 });
          return;
        }
        setSession(rtcSession);
        setRinging(true);
        setInCall(true);
        setCallerId(rtcSession.remote_identity.uri.toString());

        console.error(
          'Got a call for ' + rtcSession.remote_identity.uri.toString(),
        );
        rtcSession.on('accepted', () => {
          console.log('call accepted');
          if (remoteVideo.current) {
            remoteVideo.current.srcObject =
              session?.connection.getRemoteStreams()[0]
                ? session?.connection.getRemoteStreams()[0]
                : null;
            remoteVideo.current.play();
          }
        });
        rtcSession.on('ended', (e: EndEvent) => {
          console.log('call ended with cause: ' + JSON.stringify(e.cause));
          setInCall(false);
        });

        rtcSession.on('failed', (e: EndEvent) => {
          console.log('call failed with cause: ' + JSON.stringify(e.cause));
          setRinging(false);
        });
      },
    );
    ua.start();
    setUserAgent(ua);
  };

  const makeCall = (destination: string) => {
    const eventHandlers: Partial<RTCSessionEventMap> = {
      progress: function () {
        console.log('call is in progress');
        setInCall(true);
      },
      failed: function (e: EndEvent) {
        console.log('call failed with cause: ' + JSON.stringify(e.cause));
        setInCall(false);
      },
      ended: function (e: EndEvent) {
        console.log('call ended with cause: ' + JSON.stringify(e.cause));
        setInCall(false);
      },
      confirmed: function (e: IncomingAckEvent | OutgoingAckEvent) {
        console.log('call confirmed');
        setInCall(true);
      },
    };

    const options: any = {
      eventHandlers: eventHandlers,
      mediaConstraints: { audio: true, video: false },
    };
    if (process.env.NEXT_PUBLIC_TURN_SERVER) {
      options['pcConfig'] = {
        iceServers: [
          {
            urls: [process.env.NEXT_PUBLIC_TURN_SERVER],
            username: process.env.NEXT_PUBLIC_TURN_USER,
            credential: process.env.NEXT_PUBLIC_TURN_USER,
          },
        ],
      };
    }

    setInCall(true);
    const newSession: RTCSession | undefined = userAgent?.call(
      'sip:' + destination + '@oasis.openline.ai',
      options,
    );
    setSession(newSession);
    const peerConnection = newSession?.connection;
    peerConnection?.addEventListener('addstream', (event: any) => {
      if (remoteVideo.current) {
        remoteVideo.current.srcObject = event.stream;
      }
      remoteVideo.current?.play();
    });
  };

  const answerCall = () => {
    setInCall(true);
    setRinging(false);
    session?.answer();
  };

  const hangupCall = () => {
    setInCall(false);
    setRinging(false);
    if (session) {
      session.terminate();
    }
  };

  const holdCall = () => {
    session?.hold();
    setIsCallOnHold(true);
  };

  const unHoldCall = () => {
    session?.unhold();
    setIsCallOnHold(false);
  };

  const muteCall = () => {
    session?.mute();
    setIsCallMuted(true);
  };

  const unMuteCall = () => {
    session?.unmute();
    setIsCallMuted(false);
  };
  const sendDtmf = (digit: string) => {
    session?.sendDTMF(digit);
  };

  const value = {
    inCall,
    setInCall,
    callerId,
    setCallerId,
    ringing,
    setRinging,
    remoteVideo,
    makeCall,
    answerCall,
    hangupCall,
    holdCall,
    unHoldCall,
    muteCall,
    unMuteCall,
    sendDtmf,
    isCallOnHold,
    isCallMuted,
  };

  return (
    <WebRTCContext.Provider value={value}>
      {props.children}
    </WebRTCContext.Provider>
  );
};

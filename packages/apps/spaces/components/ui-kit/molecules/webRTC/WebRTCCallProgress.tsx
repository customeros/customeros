import * as React from 'react';
import { useContext, useRef } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faMicrophone,
  faMicrophoneSlash,
  faPause,
  faPhone,
  faPhoneSlash,
  faPlay,
} from '@fortawesome/free-solid-svg-icons';

import { OverlayPanel } from 'primereact/overlaypanel';
import { Button } from 'primereact/button';
import { WebRTCContext } from '../../../../context/web-rtc';
import {useRecoilValue} from "recoil";
import {userData} from "../../../../state";

export const WebRTCCallProgress: React.FC<any> = () => {
  const {
    inCall,
    isCallMuted,
    muteCall,
    unMuteCall,
    isCallOnHold,
    holdCall,
    unHoldCall,
    sendDtmf,
    hangupCall,
  } = useContext(WebRTCContext);

  const from = useRecoilValue(userData);
  const phoneContainerRef = useRef<OverlayPanel>(null);

  const toggleMute = () => {
    if (isCallMuted) {
      unMuteCall();
    } else {
      muteCall();
    }
  };

  const toggleHold = () => {
    if (isCallOnHold) {
      unHoldCall();
    } else {
      holdCall();
    }
  };

  const makeButton = (number: string) => {
    return (
      <button
        className='btn btn-primary btn-lg m-1'
        key={'dtmf-' + number}
        onClick={() => {
          sendDtmf(number);
        }}
      >
        {number}
      </button>
    );
  };

  const dialpad_matrix = new Array(4);
  for (let i = 0, digit = 1; i < 3; i++) {
    dialpad_matrix[i] = new Array(3);
    for (let j = 0; j < 3; j++, digit++) {
      dialpad_matrix[i][j] = makeButton(digit.toString());
    }
  }
  dialpad_matrix[3] = new Array(3);
  dialpad_matrix[3][0] = makeButton('*');
  dialpad_matrix[3][1] = makeButton('0');
  dialpad_matrix[3][2] = makeButton('#');

  const dialpad_rows = [];
  for (let i = 0; i < 4; i++) {
    dialpad_rows.push(
      <div
        key={'dtmf-row-' + i}
        className='d-flex flex-row justify-content-center'
      >
        {dialpad_matrix[i]}
      </div>,
    );
  }

  return (
    <>
      {inCall && (
        <>
          <Button
            className='p-button-rounded p-button-success p-2'
            onClick={(e: any) => phoneContainerRef?.current?.toggle(e)}
          >
            <FontAwesomeIcon icon={faPhone} fontSize={'16px'} />
          </Button>

          <OverlayPanel ref={phoneContainerRef} dismissable>
            <div
              style={{ position: 'relative', width: '100%', height: '100%' }}
            >
              <div className='font-bold text-center'>In call with</div>
              <div className='font-bold text-center mb-3'>{dialpad_rows}</div>

              <div className='font-bold text-center mb-3'>{from.identity}</div>
              <div className='mb-3'>
                <Button onClick={() => toggleMute()} className='mr-2'>
                  <FontAwesomeIcon
                    icon={isCallMuted ? faMicrophone : faMicrophoneSlash}
                    className='mr-2'
                  />{' '}
                  {isCallMuted ? 'Unmute' : 'Mute'}
                </Button>
                <Button onClick={() => toggleHold()} className='mr-2'>
                  <FontAwesomeIcon
                    icon={isCallOnHold ? faPlay : faPause}
                    className='mr-2'
                  />{' '}
                  {isCallOnHold ? 'Release hold' : 'Hold'}
                </Button>
                <Button
                  onClick={() => hangupCall()}
                  className='p-button-danger mr-2'
                >
                  <FontAwesomeIcon icon={faPhoneSlash} className='mr-2' />{' '}
                  Hangup
                </Button>
              </div>
            </div>
          </OverlayPanel>
        </>
      )}
    </>
  );
};

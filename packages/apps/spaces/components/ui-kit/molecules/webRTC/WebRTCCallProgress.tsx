import * as React from 'react';
import { useContext } from 'react';
import { WebRTCContext } from '../../../../context/web-rtc';
import { Dialog } from 'primereact/dialog';
import { useRecoilValue } from 'recoil';
import { default as Play } from '@spaces/atoms/icons/Play';
import PhoneSlashed from '@spaces/atoms/icons/PhoneSlashed';
import Pause from '@spaces/atoms/icons/Pause';
import MicrophoneSlashed from '@spaces/atoms/icons/MicrophoneSlashed';
import Microphone from '@spaces/atoms/icons/Microphone';
import { Button } from '@spaces/atoms/button';
import { IconButton } from '@spaces/atoms/icon-button/IconButton';
import { callParticipant } from '../../../../state';
import styles from './web-rtc.module.scss';
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
    ringing,
    callerId,
  } = useContext(WebRTCContext) as any;
  const { identity } = useRecoilValue(callParticipant);

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

  const getRows = () => {
    const makeButton = (number: string) => {
      return (
        <Button
          mode='secondary'
          key={'dtmf-' + number}
          onClick={() => {
            sendDtmf(number);
          }}
        >
          {number}
        </Button>
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
        <div key={'dtmf-row-' + i} className={styles.dialNumbersRow}>
          {dialpad_matrix[i]}
        </div>,
      );
    }

    return dialpad_rows;
  };

  if (!inCall) {
    return null;
  }

  return (
    <Dialog
      visible={inCall && !ringing}
      modal={false}
      className={styles.overlayContentWrapper}
      closable={false}
      closeOnEscape={false}
      draggable={false}
      onHide={() => console.log()}
    >
      <article>
        <h1 className={styles.sectionTitle}>In call with {identity}</h1>

        <div className={styles.dialNumbers}>{getRows()}</div>

        <div className={styles.actionButtonsRow}>
          <IconButton
            size='xxs'
            mode='primary'
            label={isCallMuted ? 'Unmute' : 'Mute'}
            onClick={() => toggleMute()}
            icon={
              isCallMuted ? (
                <Microphone height={16} width={16} color='#fff' />
              ) : (
                <MicrophoneSlashed height={16} width={16} color='#fff' />
              )
            }
          />

          <IconButton
            size='xxs'
            label={isCallOnHold ? 'Resume' : 'Put on hold'}
            mode='primary'
            onClick={() => toggleHold()}
            icon={
              isCallOnHold ? (
                <Play height={16} color='#fff' />
              ) : (
                <Pause height={16} color='#fff' />
              )
            }
          />

          <IconButton
            size='xxs'
            label='Hang up'
            onClick={() => hangupCall()}
            mode='danger'
            icon={<PhoneSlashed height={16} color='#fff' />}
          />
        </div>
      </article>
    </Dialog>
  );
};

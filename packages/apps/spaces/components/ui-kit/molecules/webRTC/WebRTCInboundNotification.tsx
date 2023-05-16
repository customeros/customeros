import React, { useContext } from 'react';
import { Dialog } from 'primereact/dialog';
import { WebRTCContext } from '../../../../context/web-rtc';
import styles from './web-rtc.module.scss';
import { IconButton } from '@spaces/atoms/icon-button/IconButton';
import Phone from '@spaces/atoms/icons/Phone';
export const WebRTCInboundNotification: React.FC<any> = () => {
  const { inCall, ringing, remoteVideo, answerCall, hangupCall, callerId } =
    useContext(WebRTCContext) as any;
  return (
    <>
      <video
        controls={false}
        hidden={false} //!isInCall
        ref={remoteVideo}
        autoPlay
        style={{ width: '0px', height: '0px', position: 'absolute' }}
      />

      <Dialog
        visible={ringing && inCall}
        modal={false}
        className={styles.incomingCallContainer}
        style={{ position: 'absolute', top: '25px' }}
        closable={false}
        closeOnEscape={false}
        draggable={false}
        onHide={() => console.log()}
        footer={
          <div className={styles.actionButtonsRow}>
            <IconButton
              label='Answer the phone'
              mode='primary'
              onClick={() => answerCall()}
              icon={<Phone />}
            />
            <IconButton
              label='Hang up'
              mode='danger'
              onClick={() => hangupCall()}
              icon={<Phone style={{ transform: 'rotate(133deg)' }} />}
            />
          </div>
        }
      >
        <div className={styles.incomingCall}>
          Incoming call from
          <span>{callerId}</span>
        </div>
      </Dialog>
    </>
  );
};

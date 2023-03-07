import React, { useContext } from 'react';
import { Dialog } from 'primereact/dialog';
import { WebRTCContext } from '../../../../context/web-rtc';
import styles from './web-rtc.module.scss';
import { Button, IconButton, Phone, Times } from '../../atoms';
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
              mode='primary'
              onClick={() => answerCall()}
              icon={<Phone />}
            />
            <IconButton
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

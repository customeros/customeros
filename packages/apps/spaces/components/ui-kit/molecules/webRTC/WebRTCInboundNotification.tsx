import React, { useContext } from 'react';
import { Dialog } from 'primereact/dialog';
import { Button } from 'primereact/button';
import { WebRTCContext } from '../../../../context/web-rtc';

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
        style={{ background: 'red', position: 'absolute', top: '25px' }}
        closable={false}
        closeOnEscape={false}
        draggable={false}
        onHide={() => console.log()}
        footer={
          <div>
            <Button
              label='Accept the call'
              icon='pi pi-check'
              onClick={() => answerCall()}
              className='p-button-success'
            />
            <Button
              label='Reject the call'
              icon='pi pi-times'
              onClick={() => hangupCall()}
              className='p-button-danger'
            />
          </div>
        }
      >
        <div
          className='w-full text-center font-bold'
          style={{ fontSize: '25px' }}
        >
          Incoming call from {callerId}
        </div>
      </Dialog>
    </>
  );
};

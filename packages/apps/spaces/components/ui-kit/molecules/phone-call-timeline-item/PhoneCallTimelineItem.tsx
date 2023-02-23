import React from 'react';
import styles from './phone-call-timeline-item.module.scss';
interface Props {
  phoneCallParties: {
    callingParty: string;
    calledParty: string;
  };
  duration: string;
}

export const PhoneCallTimelineItem: React.FC<Props> = ({
  phoneCallParties,
  duration,
}) => {
  return (
    <div className={styles.phoneCallContainer}>
      <div> {phoneCallParties?.callingParty || 'unknown user'} </div>
      <div className='flex flex-column'>
        <div className={styles.phoneCallVoiceWave}>
          <span />
          <span />
          <span />
          <span />
          <span />
        </div>
        {/*<Moment duration={duration}/>*/}
      </div>

      <div> {phoneCallParties?.calledParty || 'unknown user'} </div>
    </div>
  );
};

import React, { useRef, useState } from 'react';
import classNames from 'classnames';
import Timekeeper from 'react-timekeeper';
import styles from './time-picker.module.scss';
import { DateTimeUtils } from '../../../../../../utils';
import { useDetectClickOutside } from '../../../../../../hooks';
import { toast } from 'react-toastify';

interface TimePickerProps {
  alignment: 'left' | 'right';
  dateTime: Date;
  label: string;
  onUpdateTime: (newDate: Date) => void;
}

export const TimePicker: React.FC<TimePickerProps> = ({
  alignment,
  dateTime,
  label,
  onUpdateTime,
}) => {
  const [time, setTime] = useState(dateTime);
  const [timePickerOpen, setTimePickerOpen] = useState(false);
  const timePickerWrapperRef = useRef(null);

  useDetectClickOutside(timePickerWrapperRef, () => {
    setTimePickerOpen(false);
  });
  return (
    <div
      ref={timePickerWrapperRef}
      className={classNames(
        styles.date,
        styles[alignment],
        `time-picker-${alignment}`,
      )}
    >
      <button
        className={classNames(styles.timeWrapper, styles[alignment], alignment)}
        onClick={() => setTimePickerOpen(!timePickerOpen)}
      >
        <span className={styles.tinyTitle}>{label}</span>
        <span>
          {time ? (
            DateTimeUtils.formatTime(time?.toString())
          ) : (
            <span>&#8211;&#8211; : &#8211;&#8211; </span>
          )}
        </span>
      </button>
      {timePickerOpen && (
        <Timekeeper
          hour24Mode
          // switchToMinuteOnHourSelect
          // closeOnMinuteSelect
          time={time ? DateTimeUtils.formatTime(time.toString()) : undefined}
          onChange={(e) => {
            const date = new Date(time).setHours(e.hour, e.minute);
            try {
              const newDateTime = new Date(date);
              console.log('🏷️ ----- date: ', newDateTime);
              onUpdateTime(newDateTime);
              setTime(newDateTime);
            } catch (e) {
              toast.error('Invalid date selected');
            }
          }}
        />
      )}
    </div>
  );
};

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
}

export const TimePicker: React.FC<TimePickerProps> = ({
  alignment,
  dateTime,
  label,
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
          time={time ? DateTimeUtils.formatTime(time.toString()) : undefined}
          onChange={(e) => {
            const date = (time || new Date()).setHours(e.hour, e.minute);
            try {
              const newDateTime = new Date(date);
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

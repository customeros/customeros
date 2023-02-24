import React from 'react';
import { useRouter } from 'next/router';
import styles from './fullscreen-mode.module.scss';
import { Button, ChevronLeft } from '../../atoms';

interface Props {
  fullScreenMode: boolean;
  children: React.ReactNode;
  classNames?: string;
}

export const FullScreenModeLayout: React.FC<Props> = ({
  fullScreenMode,
  children,
  classNames,
}) => {
  const router = useRouter();

  return (
    <div
      className={
        fullScreenMode
          ? `${classNames} ${styles.fullScreenModeContainer}`
          : classNames
      }
    >
      {fullScreenMode && (
        <div style={{ width: '40px' }}>
          <Button icon={<ChevronLeft />} onClick={() => router.push('/')}>
            Back
          </Button>
        </div>
      )}
      {children}
    </div>
  );
};

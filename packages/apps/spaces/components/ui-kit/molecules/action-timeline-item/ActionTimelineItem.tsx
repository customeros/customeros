import React from 'react';
import styles from './action.module.scss';
import { Action } from '@spaces/graphql';
import { Company } from '@spaces/atoms/icons';

interface Props {
  action: Action;
}

export const ActionTimelineItem: React.FC<Props> = ({ action }) => {
  return (
    <>
      <div className={styles.actionWrapper}>
        {action.actionType === 'CREATED' && (
          <div className={styles.action}>
            <div className={styles.actionIcon}>
              <Company width='24' height='24' viewBox='0 0 24 24' />
            </div>
            <div className={styles.actionLabel}>Created</div>
          </div>
        )}
      </div>
    </>
  );
};

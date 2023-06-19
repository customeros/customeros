import { forwardRef, useRef } from 'react';
import type { MenuProps } from 'primereact/menu';

import { IconButton } from '@spaces/atoms/icon-button';
import { EllipsesV } from '@spaces/atoms/icons';
import { OverlayPanel } from '@spaces/atoms/overlay-panel';

import styles from './TActions.module.scss';

interface TActionsProps {
  model: MenuProps['model'];
}

export const TActions = forwardRef<HTMLDivElement, TActionsProps>(
  ({ model }, ref) => {
    const op = useRef(null);

    return (
      <div className={styles.actionHeader} ref={ref}>
        <IconButton
          label='Actions'
          className={styles.actionsMenuButton}
          id={'finder-actions-dropdown-button'}
          mode='secondary'
          size='xxxxs'
          //@ts-expect-error revisit
          onClick={(e) => op?.current?.toggle(e)}
          icon={
            <EllipsesV
              height={24}
              width={24}
              style={{ transform: 'rotate(90deg)' }}
            />
          }
        />

        <OverlayPanel
          ref={op}
          model={model}
          style={{
            maxHeight: '400px',
            height: 'fit-content',
            overflowX: 'hidden',
            overflowY: 'auto',
            bottom: 0,
          }}
        />
      </div>
    );
  },
);

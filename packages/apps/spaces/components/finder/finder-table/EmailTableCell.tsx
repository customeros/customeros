import { useRecoilValue } from 'recoil';
import React, { useRef } from 'react';
import { Button } from '@spaces/atoms/button';
import { OverlayPanel } from '@spaces/atoms/overlay-panel';
import styles from './finder-table.module.scss';
import { Contact } from '../../../graphQL/__generated__/generated';
import { FinderCell } from './FinderTableCell';
import { uuidv4 } from '../../../utils';

export const EmailTableCell = ({ emails }: { emails: Contact['emails'] }) => {
  const op = useRef(null);

  if (!emails?.length) {
    return <span>-</span>;
  }

  if (emails.length === 1) {
    return <FinderCell label={emails[0]?.email || '-'} />;
  }
  const primary = (emails || []).find((data: any) => data.primary);

  return (
    <div>
      <Button
        role='button'
        mode='text'
        style={{ padding: 0 }}
        // @ts-expect-error revisit
        onClick={(e) => op?.current?.toggle(e)}
      >
        <FinderCell label={primary?.email || emails[0]?.email || '-'} />
        <span style={{ marginLeft: '8px' }}>(...)</span>
      </Button>
      <OverlayPanel
        ref={op}
        style={{
          maxHeight: '400px',
          height: 'fit-content',
          overflowX: 'hidden',
          overflowY: 'auto',
          bottom: 0,
        }}
      >
        <ul className={styles.adressesList}>
          {emails
            .filter((d: any) => !!d?.email)
            .map((data: any) => (
              <li
                key={data.id}
                style={{ display: 'flex' }}
                className={styles.emailList}
              >
                <FinderCell label={data.email} />
              </li>
            ))}
        </ul>
      </OverlayPanel>
    </div>
  );
};

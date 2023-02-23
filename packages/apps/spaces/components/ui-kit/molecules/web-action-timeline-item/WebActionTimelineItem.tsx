import React from 'react';
import format from 'date-fns/format';
import Link from 'next/link';

import styles from './web-action-timeline-item.module.scss';
import { Globe } from '../../atoms';
import { capitalizeFirstLetter } from '../../../../utils';

// interface Props extends ContactWebAction {
//     contactName?: string
// }

export const WebActionTimelineItem = ({
  startedAt,
  pageTitle,
  pageUrl,
  application,
  engagedTime,
  contactName,
  ...rest
}: any): JSX.Element => {
  return (
    <div className={styles.x}>
      <Globe />
      <article className={`${styles.actionContainer}`}>
        <div className='flex align-items-center'>
          <div>
            {contactName && <div className='text-gray-700'>{contactName}</div>}

            <div>
              <span className='mr-1 text-gray-700'>Visited: </span>
              <Link
                href={pageUrl}
                style={{ fontWeight: 'bolder' }}
                className='overflow-hidden text-overflow-ellipsis'
              >
                {pageTitle}
              </Link>
            </div>

            <div className='flex text-gray-700'>
              <span className='mr-2 font-bolder text-gray-700'>Duration:</span>

              <span>{engagedTime ? engagedTime : '-'} minutes</span>
            </div>
          </div>
        </div>
        <div>
          <div className='flex flex-column'>
            <div>
              <span className='mr-2 font-bolder text-gray-700'>
                Started at:
              </span>

              {format(startedAt, 'dd/mm/yyyy h:mm a')}
            </div>
          </div>
          <div>
            <span className='mr-1 text-gray-700 font-bolder'>
              Accessed from:{' '}
            </span>
            <span className={styles.actionDevice}>
              {application ? capitalizeFirstLetter(application) : '-'}
            </span>
          </div>
        </div>
      </article>
    </div>
  );
};

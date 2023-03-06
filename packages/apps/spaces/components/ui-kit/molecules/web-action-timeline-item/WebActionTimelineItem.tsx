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
        <div>
          <div>
            {contactName && <div className='text-gray-700'>{contactName}</div>}

            <div>
              <span>Visited: </span>
              <Link
                href={pageUrl}
                target={'_blank'}
                style={{ fontWeight: 'bolder' }}
              >
                {pageTitle && pageTitle !== '' && <span>{pageTitle}</span>}

                {(!pageTitle || pageTitle === '') && <span>{pageUrl}</span>}
              </Link>
            </div>

            <div className='flex text-gray-700'>
              <span className='mr-2 font-bolder text-gray-700'>
                Duration:&nbsp;
              </span>

              <span>{engagedTime} minutes</span>
            </div>
          </div>
        </div>
        <div>
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

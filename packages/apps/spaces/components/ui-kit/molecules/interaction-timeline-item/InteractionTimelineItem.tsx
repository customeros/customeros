import React from 'react';

import styles from './interaction-timeline-item.module.scss';
import { Globe } from '../../atoms';
import Sitemap from '../../atoms/icons/Sitemap';
import ShoppingBag from '../../atoms/icons/ShoppingBag';
import { uuidv4 } from '../../../../utils';

// interface Props extends ContactWebAction {
//     contactName?: string
// }

export const InteractionTimelineItem = ({
  startedAt,
  name,
  status,
  type,
  channel,
  events,
  ...rest
}: any): JSX.Element => {
  return (
    <div className={styles.x}>
      <ShoppingBag />
      <article className={`${styles.actionContainer}`}>
        <div>
          <div>
            {name && <div className='text-gray-700 mb-3'>{name}</div>}

            {events.map((event: any) => {
              return <div key={uuidv4()}>{event.content}</div>;
            })}
          </div>
        </div>
      </article>
    </div>
  );
};

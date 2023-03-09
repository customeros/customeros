import React from 'react';

import styles from './interaction-timeline-item.module.scss';
import ShoppingBag from '../../atoms/icons/ShoppingBag';
import { uuidv4 } from '../../../../utils';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faImage } from '@fortawesome/free-solid-svg-icons';

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
      <ShoppingBag className={styles.icon} />
      <article className={`${styles.actionContainer}`}>
        <div>
          <div>
            {name && (
              <div
                className='text-gray-700 mb-3'
                dangerouslySetInnerHTML={{ __html: name }}
              ></div>
            )}

            {events.map((event: any) => {
              if (event.content.toLowerCase().indexOf('photos') > -1) {
                const linksExtracted = event.content.substring(
                  event.content.indexOf(':') + 1,
                  event.content.length,
                );
                const links = linksExtracted
                  .split(',')
                  .map((photo: any, index: any) => {
                    if (photo.indexOf('http') == -1) return '';
                    return (
                      <a
                        target={'_blank'}
                        rel='noreferrer'
                        key={uuidv4()}
                        href={photo}
                        className={styles.photoLink}
                      >
                        <div
                          className={'flex flex-column align-items-center mr-5'}
                        >
                          <FontAwesomeIcon icon={faImage} size='2x' />
                          Photo {index + 1}
                        </div>
                      </a>
                    );
                  });

                if (links.filter((link: any) => link != '').length == 0) {
                  return '';
                }

                return (
                  <div key={uuidv4()} className={'flex flex-column'}>
                    <div className={'text-gray-700 mb-2'}>Photos:</div>
                    <div className={'flex flex-row'}>{links}</div>
                  </div>
                );
              } else {
                const keyExtracted = event.content.substring(
                  0,
                  event.content.indexOf(':') + 1,
                );
                const valueExtracted = event.content.substring(
                  event.content.indexOf(':') + 1,
                  event.content.length,
                );
                return (
                  <div key={uuidv4()} className={'flex flex-row mb-1'}>
                    <div className={'text-gray-700 mr-1'}>{keyExtracted}</div>
                    <div>{valueExtracted}</div>
                  </div>
                );
              }
            })}
          </div>
        </div>
      </article>
    </div>
  );
};

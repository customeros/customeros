import React from 'react';

import styles from './interaction-timeline-item.module.scss';
import { uuidv4 } from '../../../../utils';
import { Avatar } from '../../atoms';
import { useRecoilValue } from 'recoil';
import { userData } from '../../../../state';

export const InteractionTimelineItem = ({ name, events }: any): JSX.Element => {
  const { identity } = useRecoilValue(userData);
  return (
    <div className={styles.folder}>
      <article>
        <div className={styles.content}>
          <div>
            <div className={styles.title}>
              <div className='flex align-items-center'>
                {/*{contactId && <ContactAvatar contactId={contactId} size={35} />}*/}
                {<Avatar name={'todo'} surname='' isSquare size={35} />}

                {name && (
                  <div
                    className='text-gray-700 ml-2'
                    dangerouslySetInnerHTML={{ __html: name }}
                  ></div>
                )}
              </div>
              <Avatar name={identity} surname='' size={35} isSquare />
            </div>

            <div className={styles.events}>
              {events?.map((event: any) => {
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
                          IMG {index + 1}
                        </a>
                      );
                    });

                  if (links.filter((link: any) => link != '').length == 0) {
                    return '';
                  }

                  return (
                    <div key={uuidv4()} className={styles.dataWrapper}>
                      <b className={styles.label}>Photos:</b>
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
                    <div key={uuidv4()} className={styles.dataWrapper}>
                      <b className={styles.label}>{keyExtracted}</b>
                      <div>{valueExtracted}</div>
                    </div>
                  );
                }
              })}
            </div>
          </div>
        </div>
      </article>
    </div>
  );
};

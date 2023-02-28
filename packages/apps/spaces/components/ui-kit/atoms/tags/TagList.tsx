import React from 'react';
import styles from './tags.module.scss';
import { capitalizeFirstLetter } from '../../../../utils';
import classNames from 'classnames';

export const TagsList = ({
  tags,
  onTagDelete,
  readOnly,
}: {
  tags: Array<{ name: string; id: string }>;
  readOnly?: boolean;
  onTagDelete?: (id: string) => void;
}) => {
  return (
    <ul
      className={classNames(styles.tagsList, {
        [styles.tagListPresentation]: readOnly,
      })}
    >
      {tags?.map((tag: { name: string; id: string }) => (
        <li key={tag.id} className={styles.tag}>
          {capitalizeFirstLetter(tag.name)?.split('_')?.join(' ')}
          {!readOnly && (
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-expect-error
            <span className='ml-2' onClick={(e) => onTagDelete(tag.id)}>
              x
            </span>
          )}
        </li>
      ))}
    </ul>
  );
};

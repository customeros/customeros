import React, { ReactNode } from 'react';
import styles from './tags.module.scss';
import { capitalizeFirstLetter } from '../../../../utils';
import classNames from 'classnames';

export const TagsList = ({
  tags,
  onTagDelete,
  readOnly,
  children,
}: {
  tags: Array<{ name: string; id: string }>;
  readOnly?: boolean;
  onTagDelete?: (id: string) => void;
  children?: ReactNode;
}) => {
  return (
    <ul
      className={classNames(styles.tagsList, {
        [styles.tagListPresentation]: readOnly,
      })}
    >
      <>
        {(tags ?? []).map((tag: { name: string; id: string }) => {
          const label = capitalizeFirstLetter(tag?.name ?? '')
            .split('_')
            .join(' ');

          return (
            <li key={tag.id} className={styles.tag}>
              {label}
              {!readOnly && onTagDelete && (
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                <span
                  className={styles.deleteButton}
                  onClick={() => onTagDelete(tag.id)}
                >
                  x
                </span>
              )}
            </li>
          );
        })}
      </>

      {children && <li key='add-tag-input'>{children}</li>}
    </ul>
  );
};

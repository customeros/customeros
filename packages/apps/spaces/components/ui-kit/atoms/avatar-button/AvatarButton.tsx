import React, { FC } from 'react';
import { Avatar } from '../avatar';
import { StaticImageData } from 'next/image';
import styles from './avatar-button.module.scss';
import { User } from '../icons';

interface Props extends Partial<HTMLButtonElement> {
  image?: StaticImageData;
  onClick?: () => void;
  ariaLabel: string;
}

export const AvatarButton: FC<Props> = ({ image, onClick, ariaLabel }) => {
  return (
    <div
      onClick={onClick}
      aria-label={ariaLabel}
      role='button'
      tabIndex={0}
      className={styles.button}
    >
      {image ? (
        <Avatar
          name={ariaLabel.split(' ')?.[0] || ''}
          surname={ariaLabel.split(' ')?.[1] || ''}
          size={30}
        />
      ) : (
        <User />
      )}
    </div>
  );
};

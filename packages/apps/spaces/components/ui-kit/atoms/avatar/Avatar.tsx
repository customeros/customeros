import React, { ReactNode, useMemo } from 'react';
import styles from './avatar.module.scss';
import { getInitialsColor } from './utils';
import classNames from 'classnames';
import { Tooltip } from '../tooltip';
import { uuidv4 } from '../../../../utils';

interface AvatarProps {
  name: string;
  surname: string;
  size: number;
  image?: ReactNode;
  imageHeight?: number;
  imageWidth?: number;
  isSquare?: boolean;
  id?: string;
}

export const Avatar: React.FC<AvatarProps> = ({
  name,
  surname,
  size,
  image,
  imageWidth,
  imageHeight,
  isSquare = false,
  id,
  ...rest
}) => {
  const initials = `${name?.charAt(0)}${surname?.charAt(0)}`;

  const color = useMemo(() => getInitialsColor(initials || 'A'), [initials]);
  const avatarStyle = {
    width: `${size}px`,
    height: `${size}px`,
    backgroundColor: image ? 'var(--gradient-base)' : color,
    fontSize: size > 40 ? 'var(--font-size-lg)' : 'var(--font-size-xxs)',
  };
  const tooltipId =
    (name || surname) && `avatar${uuidv4().split('-').join('')}`;
  return (
    <>
      {tooltipId && (
        <Tooltip
          content={`${name || ''} ${surname || ''}`}
          target={`#${tooltipId}`}
          position='top'
          showDelay={0}
          autoHide={false}
        />
      )}

      <div
        id={tooltipId || ''}
        className={classNames(styles.avatar, {
          [styles.square]: isSquare,
        })}
        style={avatarStyle}
      >
        {image || initials}
      </div>
    </>
  );
};

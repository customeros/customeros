import React from 'react';
import Image, { StaticImageData } from 'next/image';
import styles from './avatar.module.scss';
import { getInitialsColor } from './utils';

interface AvatarProps {
  name: string;
  surname: string;
  size: number;
  image?: StaticImageData;
  imageHeight?: number;
  imageWidth?: number;
}

function hashString(str: string): number {
  let hash = 0;
  if (str.length === 0) {
    return hash;
  }
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i);
    hash = (hash << 5) - hash + char;
    hash |= 0; // Convert to 32bit integer
  }
  return hash;
}
export const Avatar: React.FC<AvatarProps> = ({
  name,
  surname,
  size,
  image,
  imageWidth,
  imageHeight,
  ...rest
}) => {
  if (image) {
    return (
      <>
        <Image
          {...rest}
          src={image}
          alt={`${name} ${surname}`}
          height={imageHeight || 40}
          width={imageWidth}
        />
      </>
    );
  }

  const initials = `${name.charAt(0)}${surname.charAt(0)}`;
  const color = getInitialsColor(initials);

  const avatarStyle = {
    width: `${size}px`,
    height: `${size}px`,
    backgroundColor: color,
  };

  return (
    <div className={styles.avatar} style={avatarStyle}>
      {initials}
    </div>
  );
};

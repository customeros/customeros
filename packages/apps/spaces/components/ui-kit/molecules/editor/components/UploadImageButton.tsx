import React, { useRef } from 'react';
import { IconButton } from '@spaces/atoms/icon-button/IconButton';
import Image from '@spaces/atoms/icons/Image';

export const UploadImageButton = ({ onFileChange }: { onFileChange: any }) => {
  const inputRef = useRef<HTMLInputElement | null>(null);

  const handleUploadClick = () => {
    inputRef.current?.click();
  };

  return (
    <>
      <IconButton
        id='custom-button'
        type={'button'}
        label='Insert picture'
        onClick={handleUploadClick}
        isSquare
        mode='text'
        size='xxxxs'
        style={{ padding: '2px', background: 'transparent' }}
        icon={
          <Image
            color='#757473'
            height={18}
          />
        }
      ></IconButton>
      <input
        type='file'
        ref={inputRef}
        onChange={onFileChange}
        style={{ display: 'none' }}
      />
    </>
  );
};

import React, { ChangeEvent, useRef } from 'react';
import Image from 'next/image';
import { Button } from '../../atoms';

interface Props {
  onFileChange: (e: ChangeEvent<HTMLInputElement>) => void;
  onSubmit: any;
  onCancel?: () => void;
  label: string;
  saving: boolean;
}
export const RichTextHeader = ({
  onFileChange,
  onSubmit,
  label,
  onCancel,
  saving,
}: Props) => {
  const inputRef = useRef<HTMLInputElement | null>(null);
  const handleUploadClick = () => {
    inputRef.current?.click();
  };

  return (
    <span className='flex justify-content-end'>
      <span className='flex flex-grow-1'>
        <button className='ql-bold' aria-label='Bold'></button>
        <button className='ql-italic' aria-label='Italic'></button>
        <button className='ql-underline' aria-label='Underline'></button>
        <button className='ql-strike' aria-label='Strike'></button>

        <button className='ql-link' aria-label='Link'></button>
        <button className='ql-code-block' aria-label='Code block'></button>
        <button className='ql-blockquote' aria-label='Blockquote'></button>

        <button
          id='custom-button'
          type={'button'}
          aria-label='Insert picture'
          onClick={() => handleUploadClick()}
        >
          <Image
            src='/icons/image.svg'
            alt='Insert image'
            color={'#6c757d'}
            width={20}
            height={20}
          />
        </button>
      </span>
      <div className='editor_save'>
        {onCancel && (
          <Button onClick={onCancel} mode='secondary' className='secondary'>
            Cancel
          </Button>
        )}
        <Button
          onClick={onSubmit}
          disabled={saving}
          mode='primary'
          className='primary'
        >
          {saving ? 'Saving...' : label}
        </Button>
      </div>
      <input
        type='file'
        ref={inputRef}
        onChange={onFileChange}
        style={{ display: 'none' }}
      />
    </span>
  );
};

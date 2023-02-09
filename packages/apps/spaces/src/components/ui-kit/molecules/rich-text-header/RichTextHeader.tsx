import { Image } from '../../atoms';
import { ChangeEvent, useRef } from 'react';

interface Props {
  onFileChange: (e: ChangeEvent<HTMLInputElement>) => void;
}
export const RichTextHeader = ({ onFileChange }: Props) => {
  const inputRef = useRef<HTMLInputElement | null>(null);
  const handleUploadClick = () => {
    inputRef.current?.click();
  };
  return (
    <span className='ql-formats'>
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
        <Image color={'#6c757d'} />
      </button>

      <input
        type='file'
        ref={inputRef}
        onChange={onFileChange}
        style={{ display: 'none' }}
      />
    </span>
  );
};

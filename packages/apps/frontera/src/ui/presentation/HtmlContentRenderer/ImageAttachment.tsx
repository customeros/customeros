import { useState, ImgHTMLAttributes } from 'react';

import { FileX03 } from '@ui/media/icons/FileX03';

export const ImageAttachment = (props: ImgHTMLAttributes<HTMLImageElement>) => {
  const [hasError, setHasError] = useState(false);

  if (hasError) {
    return (
      <div className='flex items-center gap-1'>
        <FileX03 color='gray.500' />
        <span className='text-gray-500'>Attachment missing</span>
      </div>
    );
  }

  //TODO:refactor to use Image component
  return (
    <img
      {...props}
      src={props.src}
      width={props.width}
      height={props.height}
      className='mt-2 rounded-[4px]'
      alt={props.alt || 'Attachment'}
      onError={() => setHasError(true)}
    />
  );
};

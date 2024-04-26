import { useState } from 'react';
import Image, { ImageProps } from 'next/image';

import { FileX03 } from '@ui/media/icons/FileX03';

export const ImageAttachment = (props: ImageProps) => {
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
    <Image
      {...props}
      alt={props.alt || 'Attachment'}
      className='mt-2 rounded-[4px]'
      onError={() => setHasError(true)}
      src={props.src}
      width={props.width}
      height={props.height}
    />
  );
};

import { useState } from 'react';

import { FileX03 } from '@ui/media/icons/FileX03';
import { ClientImage, ClientImageProps } from '@ui/media/Image';

export const ImageAttachment = (props: ClientImageProps) => {
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
    <ClientImage
      mt='2'
      borderRadius='4px'
      onError={() => setHasError(true)}
      {...props}
    />
  );
};

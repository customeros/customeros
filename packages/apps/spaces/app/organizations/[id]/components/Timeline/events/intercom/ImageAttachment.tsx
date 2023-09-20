import { useState } from 'react';

import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { FileX03 } from '@ui/media/icons/FileX03';
import { ClientImage, ClientImageProps } from '@ui/media/Image';

export const ImageAttachment = (props: ClientImageProps) => {
  const [hasError, setHasError] = useState(false);

  if (hasError) {
    return (
      <Flex align='center' gap='1'>
        <FileX03 color='gray.500' />
        <Text color='gray.500'>Attachment missing</Text>
      </Flex>
    );
  }

  return (
    <ClientImage
      mt='2'
      borderRadius='4px'
      onError={() => setHasError(true)}
      {...props}
    />
  );
};

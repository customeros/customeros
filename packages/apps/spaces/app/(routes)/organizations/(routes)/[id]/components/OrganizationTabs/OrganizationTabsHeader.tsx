'use client';

import { Image } from '@ui/media/Image';
import { Flex } from '@ui/layout/Flex';

export const OrganizationTabsHeader = ({
  children,
}: {
  children?: React.ReactNode;
}) => {
  return (
    <Flex
      h='32'
      w='full'
      flexDir='column'
      align='center'
      justify='flex-start'
      position='relative'
      borderTopRadius='2xl'
      bg="url('/backgrounds/organization/org-banner-1.jpeg')"
      bgSize='cover'
    >
      {children}
    </Flex>
  );
};

export const OrganizationLogo = ({ src }: { src: string }) => {
  return (
    <Flex
      mt='-0.5'
      w='25%'
      h='9'
      px='3'
      pt='2'
      pb='6px'
      bg='white'
      position='relative'
      borderBottomRadius='4'
    >
      <Image
        fill
        src={src}
        boxSize='120px'
        alt='Organization Logo'
        sx={{
          objectFit: 'contain',
        }}
      />
    </Flex>
  );
};

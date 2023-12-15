'use client';

import { Flex } from '@ui/layout/Flex';
import { useDisclosure } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { Flag04 } from '@ui/media/icons/Flag04';
import { OnboardingDetails } from '@graphql/types';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';

import { OnboardingStatusModal } from './OnboardingStatusModal';

interface OnboardingStatusProps {
  data?: OnboardingDetails | null;
}

export const OnboardingStatus = ({ data }: OnboardingStatusProps) => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <>
      <Flex
        mt='1'
        gap='4'
        w='full'
        align='center'
        onClick={onOpen}
        cursor='pointer'
        overflow='visible'
        justify='flex-start'
      >
        <FeaturedIcon colorScheme='primary'>
          <Flag04 />
        </FeaturedIcon>

        <Flex>
          <Text mr='1' fontWeight='semibold'>
            Oboarding
          </Text>
          <Text color='gray.500'> on track for 3 days</Text>
        </Flex>
      </Flex>

      <OnboardingStatusModal data={data} isOpen={isOpen} onClose={onClose} />
    </>
  );
};

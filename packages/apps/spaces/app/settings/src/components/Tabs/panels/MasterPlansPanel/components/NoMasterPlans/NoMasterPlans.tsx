import { Flex } from '@ui/layout/Flex';
import { Text } from '@ui/typography/Text';
import { Center } from '@ui/layout/Center';
import { Flag04 } from '@ui/media/icons/Flag04';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';

import { NoMasterPlansMenu } from './NoMasterPlansMenu';
import { useMasterPlansMethods } from '../../hooks/useMasterPlansMethods';

export const NoMasterPlans = () => {
  const { isPending, handleCreateDefault, handleCreateFromScratch } =
    useMasterPlansMethods();

  return (
    <Flex w='full' h='full' justify='center' pt='20'>
      <Center h='fit-content' flexDir='column'>
        <FeaturedIcon colorScheme='primary' size='lg'>
          <Flag04 />
        </FeaturedIcon>
        <Text mt='2' fontWeight='medium'>
          Create a new master plan
        </Text>
        <Text mt='1' textAlign='center' color='gray.600'>
          Help your customers be more successful by creating <br /> a new
          onboarding master plan
        </Text>
        <Flex mt='6'>
          <NoMasterPlansMenu
            isLoading={isPending}
            onCreateDefault={handleCreateDefault}
            onCreateFromScratch={handleCreateFromScratch}
          />
        </Flex>
      </Center>
    </Flex>
  );
};

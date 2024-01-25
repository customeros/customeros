import { useMemo } from 'react';
import { useRouter } from 'next/navigation';

import { Button } from '@ui/form/Button';
import { useDisclosure } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useMasterPlansQuery } from '@shared/graphql/masterPlans.generated';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
  MenuDivider,
} from '@ui/overlay/Menu';

export const OnboardingMenu = () => {
  const router = useRouter();
  const client = getGraphQLClient();
  const { isOpen, onClose, onOpen } = useDisclosure();
  const { data, isPending } = useMasterPlansQuery(client);

  const activeMasterPlans = useMemo(
    () => data?.masterPlans?.filter((m) => !m.retired),
    [data?.masterPlans],
  );

  const handleEditMasterPlans = () => {
    router.push('/settings?tab=master-plans&show=active');
  };

  return (
    <Menu isOpen={isOpen} onClose={onClose} placement='bottom-end'>
      <MenuButton
        size='sm'
        as={Button}
        variant='ghost'
        color='gray.500'
        onClick={onOpen}
        fontWeight='normal'
        isDisabled={isPending}
        leftIcon={<Plus color='gray.400' />}
      >
        Add plan
      </MenuButton>
      <MenuList maxW='200px' maxH='300px' overflowY='auto'>
        {activeMasterPlans?.map((m) => (
          <MenuItem key={m.id}>
            <Text noOfLines={1}>{m.name}</Text>
          </MenuItem>
        ))}
        <MenuItem>Plan 2</MenuItem>
        <MenuDivider />
        <MenuItem onClick={handleEditMasterPlans}>Edit master plans</MenuItem>
      </MenuList>
    </Menu>
  );
};

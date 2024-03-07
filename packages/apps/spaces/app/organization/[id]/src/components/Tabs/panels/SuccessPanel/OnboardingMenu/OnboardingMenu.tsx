import { useMemo } from 'react';
import { useRouter, useParams } from 'next/navigation';

import { Button } from '@ui/form/Button';
import { useDisclosure } from '@ui/utils';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { Map01 } from '@ui/media/icons/Map01';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useMasterPlansQuery } from '@shared/graphql/masterPlans.generated';
import {
  Menu,
  MenuItem,
  MenuList,
  MenuButton,
  MenuDivider,
} from '@ui/overlay/Menu';

import { usePlanMutations } from '../hooks/usePlanMutations';

export const OnboardingMenu = () => {
  const router = useRouter();
  const client = getGraphQLClient();
  const organizationId = (useParams()?.id ?? '') as string;
  const { isOpen, onClose, onOpen } = useDisclosure();
  const { data, isPending } = useMasterPlansQuery(client);

  const { createOnboardingPlan } = usePlanMutations({ organizationId });

  const activeMasterPlans = useMemo(
    () => data?.masterPlans?.filter((m) => !m.retired),
    [data?.masterPlans],
  );

  const handleEditMasterPlans = () => {
    router.push('/settings?tab=master-plans&show=active');
  };

  const handleCreateOnboardingPlan =
    (masterPlanId: string, name: string) => () => {
      createOnboardingPlan.mutate({
        input: {
          name,
          masterPlanId,
          organizationId,
        },
      });
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
        isLoading={createOnboardingPlan.isPending}
      >
        Add plan
      </MenuButton>
      <MenuList maxW='200px' maxH='300px' overflowY='auto'>
        {activeMasterPlans?.map((m) => (
          <MenuItem
            key={m.id}
            onClick={handleCreateOnboardingPlan(m.id, m.name)}
          >
            <Text noOfLines={1}>{m.name}</Text>
          </MenuItem>
        ))}
        <MenuDivider
          mx='2'
          borderBottom='unset'
          borderTop='1px dashed'
          borderColor='gray.300'
        />
        <MenuItem
          icon={<Map01 color='gray.500' />}
          onClick={handleEditMasterPlans}
        >
          Edit master plans
        </MenuItem>
      </MenuList>
    </Menu>
  );
};

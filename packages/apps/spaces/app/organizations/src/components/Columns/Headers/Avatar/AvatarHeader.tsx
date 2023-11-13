import { Flex } from '@ui/layout/Flex';
import { Plus } from '@ui/media/icons/Plus';
import { Tooltip } from '@ui/overlay/Tooltip';
import { IconButton } from '@ui/form/IconButton';
import { useOrganizationsPageMethods } from '@organizations/hooks/useOrganizationsPageMethods';

export const AvatarHeader = () => {
  const { createOrganization } = useOrganizationsPageMethods();

  const handleCreateOrganization = () => {
    createOrganization.mutate({ input: { name: '' } });
  };

  return (
    <Flex w='42px' align='center' justify='center'>
      <Tooltip label='Create an organization'>
        <IconButton
          size='sm'
          variant='ghost'
          aria-label='create organization'
          onClick={handleCreateOrganization}
          isLoading={createOrganization.isLoading}
          icon={<Plus color='gray.400' boxSize='5' />}
        />
      </Tooltip>
    </Flex>
  );
};

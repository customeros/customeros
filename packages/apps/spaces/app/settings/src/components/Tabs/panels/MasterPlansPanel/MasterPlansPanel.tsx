import { useRouter, useSearchParams } from 'next/navigation';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Plus } from '@ui/media/icons/Plus';
import { IconButton } from '@ui/form/IconButton';
import { Grid, GridItem } from '@ui/layout/Grid';

export const MasterPlansPanel = () => {
  return (
    <Grid templateColumns='1fr 2fr' h='full'>
      <GridItem borderRight='1px solid' borderRightColor='gray.200' p='4'>
        <Flex align='center' justify='space-between' mb='2'>
          <Text fontWeight='semibold'>Master Plans</Text>
          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Add Master Plan'
            icon={<Plus color='gray.400' />}
          />
        </Flex>

        <VStack align='flex-start'>
          {Array.from({ length: 5 }).map((_, i) => (
            <MasterPlanNavItem key={i} id={String(i)}>
              Master Plan {i + 1}
            </MasterPlanNavItem>
          ))}
        </VStack>
      </GridItem>
    </Grid>
  );
};

interface MasterPlanNavItemProps {
  id: string;
  children?: React.ReactNode;
}

const MasterPlanNavItem = ({ id, children }: MasterPlanNavItemProps) => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const isActive = searchParams?.get('planId') === id;

  const handleClick = () => {
    const newParams = new URLSearchParams(searchParams ?? '');
    newParams.set('planId', id);

    router.push(`?${newParams.toString()}`);
  };

  return (
    <Button
      px='3'
      w='full'
      fontSize='sm'
      fontWeight='normal'
      onClick={handleClick}
      justifyContent='flex-start'
      bg={isActive ? 'gray.100' : 'transparent'}
      _hover={{
        bg: 'gray.100',
      }}
      _active={{
        bg: 'gray.200',
      }}
    >
      {children}
    </Button>
  );
};

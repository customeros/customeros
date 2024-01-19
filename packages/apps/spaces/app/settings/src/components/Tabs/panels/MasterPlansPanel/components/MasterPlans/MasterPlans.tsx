import { useRouter, useSearchParams } from 'next/navigation';

import { MasterPlansQuery } from '@settings/graphql/masterPlans.generated';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { VStack } from '@ui/layout/Stack';
import { Skeleton } from '@ui/presentation/Skeleton';

interface MasterPlanNavItemProps {
  id: string;
  children?: React.ReactNode;
}

const MasterPlanItem = ({ id, children }: MasterPlanNavItemProps) => {
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

const LoadingMasterplans = () => {
  return (
    <VStack align='flex-start'>
      {Array.from({ length: 5 }).map((_, i) => (
        <Skeleton key={i} h='6' w='full' />
      ))}
    </VStack>
  );
};

const NoMasterplans = () => {
  return <Flex>No master plans created yet</Flex>;
};

interface MasterPlansProps {
  isLoading?: boolean;
  masterPlans?: MasterPlansQuery['masterPlans'];
}

export const MasterPlans = ({ masterPlans, isLoading }: MasterPlansProps) => {
  if (isLoading) return <LoadingMasterplans />;
  if (!masterPlans) return <NoMasterplans />;

  return (
    <VStack align='flex-start'>
      {masterPlans.map(({ id, name }) => (
        <MasterPlanItem key={id} id={id}>
          {name}
        </MasterPlanItem>
      ))}
    </VStack>
  );
};

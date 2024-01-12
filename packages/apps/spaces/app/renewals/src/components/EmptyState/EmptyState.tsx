import { useRouter } from 'next/navigation';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Center } from '@ui/layout/Center';
import { Text } from '@ui/typography/Text';
import { EmptyTable } from '@ui/media/logos/EmptyTable';

import HalfCirclePattern from '../../../../src/assets/HalfCirclePattern';

export const EmptyState = () => {
  const router = useRouter();

  const options = {
    title: 'No contracts created yet',
    description:
      'Currently, you have not been assigned to any organizations.\n' +
      '\n' +
      'Head to your list of organizations and assign yourself as an owner to one of them.',
    buttonLabel: 'Go to Organizations',
    onClick: () => {
      router.push(`/organizations`);
    },
  };

  return (
    <Center h='100%' bg='white'>
      <Flex direction='column' height={500} width={500}>
        <Flex position='relative'>
          <EmptyTable
            width='152px'
            height='120'
            position='absolute'
            top='25%'
            right='35%'
          />
          <HalfCirclePattern height={500} width={500} />
        </Flex>
        <Flex
          flexDir='column'
          textAlign='center'
          align='center'
          top='5vh'
          transform='translateY(-230px)'
        >
          <Text color='gray.900' fontSize='md' fontWeight='semibold'>
            {options.title}
          </Text>
          <Text maxW='400px' fontSize='sm' color='gray.600' my={1}>
            {options.description}
          </Text>

          <Button
            onClick={options.onClick}
            mt='2'
            w='min-content'
            variant='outline'
            fontSize='sm'
          >
            {options.buttonLabel}
          </Button>
        </Flex>
      </Flex>
    </Center>
  );
};

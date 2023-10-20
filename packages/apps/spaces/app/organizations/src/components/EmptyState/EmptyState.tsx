import { Center } from '@ui/layout/Center';
import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { Text } from '@ui/typography/Text';
import { EmptyTable } from '@ui/media/logos/EmptyTable';
import HalfCirclePattern from '../../../../src/assets/HalfCirclePattern';

interface EmptyStateProps {
  onClick: () => void;
  buttonLabel: string;
  title: string;
  description: string;
}

const EmptyState = ({
  onClick,
  title,
  buttonLabel,
  description,
}: EmptyStateProps) => {
  return (
    <Center
      h='100%'
      bg='white'
      borderRadius='2xl'
      border='1px solid'
      borderColor='gray.200'
    >
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
            {title}
          </Text>
          <Text maxW='400px' fontSize='sm' color='gray.600' my={1}>
            {description}
          </Text>

          <Button
            onClick={onClick}
            mt='2'
            w='min-content'
            variant='outline'
            fontSize='sm'
          >
            {buttonLabel}
          </Button>
        </Flex>
      </Flex>
    </Center>
  );
};

export default EmptyState;

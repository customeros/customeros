import { Flex } from '@ui/layout/Flex';
import { Skeleton } from '@ui/presentation/Skeleton';
import { Card, CardBody } from '@ui/presentation/Card';

export const MessageCardSkeleton = () => {
  return (
    <Card
      variant='outline'
      size='md'
      fontSize='14px'
      background='white'
      flexDirection='row'
      position='unset'
      boxShadow='xs'
      borderColor='gray.200'
      w='full'
    >
      <CardBody p='3'>
        <Flex gap='4'>
          <Skeleton
            w='10'
            h='10'
            borderRadius='6px'
            startColor='gray.300'
            endColor='gray.100'
          />

          <Flex flex='1' flexDir='column' gap='2'>
            <Skeleton
              h='16px'
              w='25%'
              startColor='gray.300'
              endColor='gray.100'
            />
            <Skeleton
              h='12px'
              w='50%'
              startColor='gray.300'
              endColor='gray.100'
            />
          </Flex>
        </Flex>
      </CardBody>
    </Card>
  );
};

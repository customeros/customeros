import React from 'react';
import { Card } from '@ui/layout/Card';
import { CardHeader, Flex, Heading, IconButton } from '@chakra-ui/react';
import { Plus } from '@ui/media/icons/Plus';
import { CardBody } from '@chakra-ui/card';
import Link from 'next/link';

interface BranchesProps {
  branches: any;
}

export const Branches: React.FC<BranchesProps> = ({ branches }) => {
  return (
    <Card size='sm'>
      <CardHeader
        display='flex'
        alignItems='center'
        justifyContent='space-between'
      >
        <Heading fontSize={'md'}>Branches</Heading>
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Add'
          onClick={() => null}
          icon={<Plus boxSize='4' />}
        />
      </CardHeader>
      <CardBody as={Flex}>
        {branches.map((e) => (
          <Link key={e.href} href={e.href} /> // todo 719
        ))}
      </CardBody>
    </Card>
  );
};

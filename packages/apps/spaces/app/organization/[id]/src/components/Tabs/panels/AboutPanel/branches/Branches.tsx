import React from 'react';
import { Card } from '@ui/layout/Card';
import { CardHeader, Heading, IconButton, VStack } from '@chakra-ui/react';
import { Plus } from '@ui/media/icons/Plus';
import { CardBody } from '@chakra-ui/card';
import { Link } from '@ui/navigation/Link';

import { Organization } from '@graphql/types';

interface BranchesProps {
  branches?: Organization['subsidiaries'];
}

export const Branches: React.FC<BranchesProps> = ({ branches = [] }) => {
  return (
    <Card size='sm' width='full' mt={2}>
      <CardHeader
        display='flex'
        alignItems='center'
        justifyContent='space-between'
        pb={4}
      >
        <Heading fontSize={'md'}>Branches</Heading>
        {/*<IconButton TODO*/}
        {/*  size='xs'*/}
        {/*  variant='ghost'*/}
        {/*  aria-label='Add'*/}
        {/*  onClick={() => null}*/}
        {/*  icon={<Plus boxSize='4' />}*/}
        {/*/>*/}
      </CardHeader>
      <CardBody as={VStack} pt={0} gap={2} alignItems='baseline'>
        {branches.map(({ organization }) =>
          organization?.website ? (
            <Link
              noOfLines={1}
              wordBreak='keep-all'
              href={`/organization/${organization.id}?tab=about`}
              key={`subsidiaries-${organization.id}`}
              color='gray.700'
              _hover={{ color: 'primary.600' }}
            >
              {organization.name}
            </Link>
          ) : null,
        )}
      </CardBody>
    </Card>
  );
};

'use client';
import React from 'react';
import { useParams, useRouter } from 'next/navigation';

import { motion } from 'framer-motion';

import { Flex } from '@ui/layout/Flex';
import { Invoice } from '@graphql/types';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetInvoicesQuery } from '@shared/graphql/getInvoices.generated';
import { InvoicesTable } from '@organization/src/components/Tabs/panels/InvoicesPanel/InvoicesTable';
import { OrganizationPanel } from '@organization/src/components/Tabs/panels/OrganizationPanel/OrganizationPanel';

const slideUpVariants = {
  initial: { y: '100%', opacity: 0 },
  animate: {
    y: 0,
    opacity: 1,
    transition: { type: 'interia', stiffness: 100 },
  },
  exit: { y: '100%', opacity: 0, transition: { duration: 3 } },
};
export const InvoicesPanel = () => {
  const id = useParams()?.id as string;
  const client = getGraphQLClient();
  const router = useRouter();
  const { data } = useGetInvoicesQuery(client, {
    pagination: {
      page: 0,
      limit: 50,
    },
    organizationId: id,
  });

  return (
    <OrganizationPanel title='Account'>
      <motion.div
        variants={slideUpVariants}
        initial='initial'
        animate='animate'
        exit='exit'
        style={{ width: '100%' }}
      >
        <Flex justifyContent='space-between' mb={2}>
          <Text fontSize='sm' fontWeight='semibold'>
            Invoices
          </Text>
          <IconButton
            aria-label='Go back'
            variant='ghost'
            size='xs'
            icon={<ChevronDown color='gray.400' />}
            onClick={() => router.push(`?tab=account`)}
          />
        </Flex>
        <Flex mx={-5}>
          <InvoicesTable
            invoices={(data?.invoices?.content as Array<Invoice>) ?? []}
            totalElements={data?.invoices.totalElements}
          />
        </Flex>
      </motion.div>
    </OrganizationPanel>
  );
};

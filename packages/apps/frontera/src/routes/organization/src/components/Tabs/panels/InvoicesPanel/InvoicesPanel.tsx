import { useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

import { motion } from 'framer-motion';
import { useIsRestoring } from '@tanstack/react-query';

import { Invoice } from '@graphql/types';
import { Table } from '@ui/presentation/Table';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { EmptyState } from '@shared/components/Invoice/EmptyState/EmptyState';
import { useInfiniteInvoices } from '@shared/components/Invoice/hooks/useInfiniteInvoices';
import { columns } from '@organization/components/Tabs/panels/InvoicesPanel/Columns/Columns';
import { OrganizationPanel } from '@organization/components/Tabs/panels/OrganizationPanel/OrganizationPanel';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

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
  const isRestoring = useIsRestoring();
  const navigate = useNavigate();
  const tableRef = useRef(null);
  const { handleOpenInvoice } = useTimelineEventPreviewMethodsContext();

  const {
    invoiceFlattenData,
    totalInvoicesCount,
    isFetching,
    fetchNextPage,
    hasNextPage,
  } = useInfiniteInvoices(id);
  if (totalInvoicesCount === 0) {
    return (
      <div className='flex justify-center'>
        <EmptyState />
      </div>
    );
  }

  return (
    <OrganizationPanel title='Account'>
      <motion.div
        key='invoices'
        variants={slideUpVariants}
        initial='initial'
        animate='animate'
        exit={{ x: -500, opacity: 0 }}
        style={{ width: '100%' }}
      >
        <div className='flex justify-between mb-2'>
          <p className='text-sm font-semibold'>Invoices</p>
          <IconButton
            aria-label='Go back'
            variant='ghost'
            size='xs'
            icon={<ChevronDown className='text-gray-400' />}
            onClick={() => navigate(`?tab=account`)}
          />
        </div>
        <div className='mx-[-5px]'>
          <Table<Invoice>
            data={invoiceFlattenData ?? []}
            columns={columns}
            enableRowSelection={false}
            fullRowSelection={true}
            onFullRowSelection={(id) => id && handleOpenInvoice(id)}
            canFetchMore={hasNextPage}
            onFetchMore={fetchNextPage}
            tableRef={tableRef}
            isLoading={isRestoring ? false : isFetching}
            totalItems={isRestoring ? 10 : totalInvoicesCount}
            rowHeight={4}
            borderColor='gray.100'
            contentHeight={'80vh'}
          />
        </div>
      </motion.div>
    </OrganizationPanel>
  );
};

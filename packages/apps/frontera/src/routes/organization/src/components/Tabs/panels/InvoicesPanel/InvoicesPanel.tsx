import { useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

import { motion } from 'framer-motion';
import { Store } from '@store/store.ts';
import { observer } from 'mobx-react-lite';

import { Invoice } from '@graphql/types';
import { Table } from '@ui/presentation/Table';
import { useStore } from '@shared/hooks/useStore';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { EmptyState } from '@shared/components/Invoice/EmptyState/EmptyState';
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
export const InvoicesPanel = observer(() => {
  const id = useParams()?.id as string;
  const navigate = useNavigate();
  const tableRef = useRef(null);
  const { handleOpenInvoice } = useTimelineEventPreviewMethodsContext();

  const store = useStore();
  const invoices = store.invoices
    .toComputedArray((a) => a)
    .filter((e) => e?.value?.organization?.metadata?.id === id);
  if (!store.invoices.isLoading && invoices.length === 0) {
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
          <Table<Store<Invoice>>
            data={invoices ?? []}
            totalItems={invoices.length}
            columns={columns}
            enableRowSelection={false}
            fullRowSelection={true}
            onFullRowSelection={(id) => id && handleOpenInvoice(id)}
            tableRef={tableRef}
            isLoading={store.invoices.isLoading}
            rowHeight={4}
            borderColor='gray.100'
            contentHeight={'80vh'}
          />
        </div>
      </motion.div>
    </OrganizationPanel>
  );
});

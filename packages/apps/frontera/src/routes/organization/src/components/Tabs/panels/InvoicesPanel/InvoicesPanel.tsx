import { useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

import { motion } from 'framer-motion';
import { observer } from 'mobx-react-lite';
import { InvoiceStore } from '@store/Invoices/Invoice.store.ts';

import { Table } from '@ui/presentation/Table';
import { useStore } from '@shared/hooks/useStore';
import { ChevronDown } from '@ui/media/icons/ChevronDown';
import { IconButton } from '@ui/form/IconButton/IconButton';
import { EmptyState } from '@shared/components/Invoice/EmptyState/EmptyState';
import { columns } from '@organization/components/Tabs/panels/InvoicesPanel/Columns/Columns';
import { OrganizationPanel } from '@organization/components/Tabs/shared/OrganizationPanel/OrganizationPanel';
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
    .filter(
      (e) => e?.value?.organization?.metadata?.id === id && !e.value.dryRun,
    );

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
        initial='initial'
        animate='animate'
        style={{ width: '100%' }}
        variants={slideUpVariants}
        exit={{ x: -500, opacity: 0 }}
      >
        <div className='flex justify-between mb-2 mr-4'>
          <p className='text-sm font-semibold'>Invoices</p>
          <IconButton
            size='xs'
            variant='ghost'
            aria-label='Go back'
            onClick={() => navigate(`?tab=account`)}
            icon={<ChevronDown className='text-gray-400' />}
          />
        </div>
        <div className='-ml-6 max-w-[447px]'>
          <Table<InvoiceStore>
            rowHeight={4}
            columns={columns}
            tableRef={tableRef}
            data={invoices ?? []}
            borderColor='gray.100'
            contentHeight={'80vh'}
            fullRowSelection={true}
            enableRowSelection={false}
            totalItems={invoices.length}
            isLoading={store.invoices.isLoading}
            onFullRowSelection={(id) => id && handleOpenInvoice(id)}
          />
        </div>
      </motion.div>
    </OrganizationPanel>
  );
});

import { useNavigate } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { Plus } from '@ui/media/icons/Plus';
import { Table } from '@ui/presentation/Table';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { LinkExternal01 } from '@ui/media/icons/LinkExternal01';

import { columns } from './columns';

export const MailboxTable = observer(() => {
  const store = useStore();
  const navigate = useNavigate();
  const flowsPreset = store.tableViewDefs.flowsPreset;

  const goToBuy = () => navigate('/settings?tab=mailboxes&view=buy');

  const data = store.mailboxes.toArray().map((v) => v.value);

  return (
    <div className='w-full'>
      <div className='pl-6 pr-3 pt-[5px] pb-[5px] flex items-center justify-between'>
        <div className='flex items-center gap-2'>
          <h2 data-test='mailboxes-header' className='font-semibold text-md'>
            Mailboxes
          </h2>
          <Button
            size='xxs'
            rightIcon={<LinkExternal01 className='size-3' />}
            onClick={() => navigate(`/finder?preset=${flowsPreset}`)}
          >
            Jump to flows
          </Button>
        </div>
        <Button
          size='xs'
          onClick={goToBuy}
          leftIcon={<Plus />}
          colorScheme='primary'
          dataTest='new-mailboxes'
        >
          New mailboxes
        </Button>
      </div>

      <Table
        data={data}
        rowHeight={28}
        columns={columns}
        contentHeight={'calc(100vh - 40px)'}
      />
    </div>
  );
});

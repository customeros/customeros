import type { Mailbox } from '@store/Settings/Mailbox.store';

import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { createColumnHelper } from '@ui/presentation/Table';

import { UserCell } from './cells';

type ColumnDatum = Mailbox;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

export const columns: Column[] = [
  columnHelper.accessor('mailbox', {
    id: 'mailbox',
    minSize: 320,
    cell: (props) => <p>{props.getValue()}</p>,
    header: () => <p className='text-sm'>Mailboxes</p>,
    skeleton: () => null,
  }),
  columnHelper.accessor('mailbox', {
    id: 'user',
    minSize: 200,
    cell: (props) => <UserCell id={props.getValue()} />,
    header: () => <p className='text-sm'>User</p>,
    skeleton: () => null,
  }),
  columnHelper.accessor('rampUpCurrent', {
    id: 'rampUpCurrent',
    minSize: 128,
    cell: (props) => <p>{props.getValue()}</p>,
    header: () => <p className='text-sm'>Daily Email limit</p>,
    skeleton: () => null,
  }),
  columnHelper.accessor('scheduledEmails', {
    id: 'scheduledEmail',
    minSize: 128,
    cell: (props) => <p>{props.getValue()}</p>,
    header: () => <p className='text-sm'>Scheduled emails</p>,
    skeleton: () => null,
  }),
  // columnHelper.accessor('currentFlowIds', {
  //   id: 'currentFlowIds',
  //   minSize: 128,
  //   cell: (props) => <p>{JSON.stringify(props.getValue())}</p>,
  //   header: () => <p className='text-sm'>Current Flows</p>,
  //   skeleton: () => null,
  // }),
];

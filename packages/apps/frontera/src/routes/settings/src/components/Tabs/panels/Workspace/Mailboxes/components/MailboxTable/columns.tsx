import type { Mailbox } from '@store/Settings/Mailbox.store';

import { ColumnDef as ColumnDefinition } from '@tanstack/react-table';

import { createColumnHelper } from '@ui/presentation/Table';

import { UserCell, MailboxCell } from './cells';
import { SenderHeader } from './header/SenderHeader';
import { RampUpCurrentHeader } from './header/RampUpCurrentHeader';

type ColumnDatum = Mailbox;

// REASON: we do not care about exhaustively typing this TValue type
// eslint-disable-next-line @typescript-eslint/no-explicit-any
type Column = ColumnDefinition<ColumnDatum, any>;

const columnHelper = createColumnHelper<ColumnDatum>();

export const columns: Column[] = [
  columnHelper.accessor('mailbox', {
    id: 'mailbox',
    minSize: 320,
    cell: (props) => <MailboxCell mailbox={props.getValue()} />,
    header: () => <p className='text-sm ml-8'>Mailboxes</p>,
    skeleton: () => null,
  }),
  columnHelper.accessor('mailbox', {
    id: 'user',
    minSize: 200,
    cell: (props) => <UserCell id={props.getValue()} />,
    header: SenderHeader,
    skeleton: () => null,
  }),
  columnHelper.accessor('rampUpCurrent', {
    id: 'rampUpCurrent',
    minSize: 200,
    cell: (props) => <p>{props.getValue()} emails / day</p>,
    header: RampUpCurrentHeader,
    skeleton: () => null,
  }),
  // columnHelper.accessor('status', {
  //   id: 'Status',
  //   minSize: 128,
  //   cell: (props) => <p>{props.getValue()}</p>,
  //   header: () => <p className='text-sm'>Status</p>,
  //   skeleton: () => null,
  // }),
  columnHelper.accessor('scheduledEmails', {
    id: 'scheduledEmail',
    minSize: 128,
    cell: (props) => <p>{props.getValue()}</p>,
    header: () => <p className='text-sm'>Scheduled Emails</p>,
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

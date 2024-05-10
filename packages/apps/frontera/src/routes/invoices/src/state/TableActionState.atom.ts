import { atom, useRecoilState } from 'recoil';

import { InvoiceStatus } from '@graphql/types';

type TableActionState = {
  targetId: string;
  isConfirming: boolean;
  targetStatus: InvoiceStatus | null;
};

export const TableActionsAtom = atom<TableActionState>({
  key: 'table-action-state',
  default: {
    targetId: '',
    targetStatus: null,
    isConfirming: false,
  },
});

export const useTableActionState = () => {
  return useRecoilState(TableActionsAtom);
};

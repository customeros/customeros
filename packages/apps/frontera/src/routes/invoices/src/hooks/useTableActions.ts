import { useEffect, useCallback } from 'react';

import { useLocalObservable } from 'mobx-react-lite';

import { useStore } from '@shared/hooks/useStore';
import { InvoiceStatus } from '@shared/types/__generated__/graphql.types';

interface TableActionsState {
  targetId: string;
  isConfirming: boolean;
  targetStatus: InvoiceStatus | null;
  setTableActionState: (
    state: Omit<TableActionsState, 'setTableActionState'>,
  ) => void;
}

export const useTableActions = () => {
  const tableState = useLocalObservable<TableActionsState>(() => ({
    targetId: '',
    targetStatus: null,
    isConfirming: false,
    setTableActionState({
      targetId,
      targetStatus,
      isConfirming,
    }: {
      targetId: string;
      isConfirming: boolean;
      targetStatus: InvoiceStatus | null;
    }) {
      this.targetId = targetId;
      this.targetStatus = targetStatus;
      this.isConfirming = isConfirming;
    },
  }));

  const store = useStore();
  const invoice = store.invoices?.value?.get(tableState.targetId);

  const reset = () => {
    tableState.setTableActionState({
      targetId: '',
      targetStatus: null,
      isConfirming: false,
    });
  };

  const onConfirm = useCallback(() => {
    if (tableState.targetId && tableState.targetStatus) {
      invoice?.update((prev) => ({
        ...prev,
        status: tableState.targetStatus,
      }));
    }
  }, [tableState.targetId, tableState.targetStatus]);

  useEffect(() => {
    if (!tableState.isConfirming) {
      onConfirm();
    }
  }, [tableState.isConfirming, onConfirm]);

  return {
    reset,
    targetId: tableState.targetId,
    onConfirm,
    isConfirming: tableState.isConfirming,
  };
};

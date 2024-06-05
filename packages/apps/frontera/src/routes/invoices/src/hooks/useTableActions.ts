import { useEffect, useCallback } from 'react';

import { useStore } from '@shared/hooks/useStore';

import { useTableActionState } from '../state/TableActionState.atom';

export const useTableActions = () => {
  const [tableActionState, setTableActionState] = useTableActionState();
  const store = useStore();
  const { targetId, targetStatus, isConfirming } = tableActionState;
  const invoice = store.invoices?.value?.get(targetId);

  const reset = () => {
    setTableActionState({
      targetId: '',
      targetStatus: null,
      isConfirming: false,
    });
  };

  const onConfirm = useCallback(() => {
    if (targetId && targetStatus) {
      invoice?.update((prev) => ({
        ...prev,
        status: targetStatus,
      }));
    }
  }, [targetId, targetStatus]);

  useEffect(() => {
    if (!isConfirming) {
      onConfirm();
    }
  }, [isConfirming, onConfirm]);

  return {
    reset,
    targetId,
    onConfirm,
    isConfirming,
  };
};

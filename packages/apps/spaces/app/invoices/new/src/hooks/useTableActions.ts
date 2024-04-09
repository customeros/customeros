import { useEffect, useCallback } from 'react';

import { produce } from 'immer';
import { useQueryClient } from '@tanstack/react-query';

import { toastError } from '@ui/presentation/Toast';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useInvoicesMeta } from '@shared/state/InvoicesMeta.atom';
import { useInfiniteGetInvoicesQuery } from '@shared/graphql/getInvoices.generated';
import { useUpdateInvoiceStatusMutation } from '@shared/graphql/updateInvoiceStatus.generated';

import { useTableActionState } from '../state/TableActionState.atom';

export const useTableActions = () => {
  const [tableActionState, setTableActionState] = useTableActionState();
  const queryClient = useQueryClient();
  const client = getGraphQLClient();
  const [invoicesMeta] = useInvoicesMeta();

  const { targetId, targetStatus, isConfirming } = tableActionState;

  const reset = () => {
    setTableActionState({
      targetId: '',
      targetStatus: null,
      isConfirming: false,
    });
  };

  const queryKey = useInfiniteGetInvoicesQuery.getKey(invoicesMeta.getInvoices);
  const updateInvoiceStatus = useUpdateInvoiceStatusMutation(client, {
    onMutate: ({ input }) => {
      queryClient.cancelQueries({ queryKey });
      const page = invoicesMeta.getInvoices.pagination.page;

      const { previousEntries } = useInfiniteGetInvoicesQuery.mutateCacheEntry(
        queryClient,
        invoicesMeta.getInvoices,
      )((cache) =>
        produce(cache, (draft) => {
          const foundInvoice = draft.pages[page].invoices.content.find(
            (i) => i.metadata.id === tableActionState.targetId,
          );

          if (foundInvoice) {
            foundInvoice.status = input.status;
          }
        }),
      );

      return { previousEntries };
    },
    onError: (_, __, context) => {
      reset();
      if (context?.previousEntries) {
        queryClient.setQueryData(queryKey, context.previousEntries);
      }
      toastError(
        `We couldn't update the invoice status`,
        'update-invoice-status-finder',
      );
    },
    onSettled: () => {
      reset();
      setTimeout(() => {
        queryClient.invalidateQueries({ queryKey });
      }, 500);
    },
  });

  const onConfirm = useCallback(() => {
    if (targetId && targetStatus) {
      updateInvoiceStatus.mutate({
        input: {
          id: targetId,
          status: targetStatus,
          patch: true,
        },
      });
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
    isPending: updateInvoiceStatus.isPending,
  };
};

import { useState, useContext, createContext, PropsWithChildren } from 'react';

import { Invoice } from '@graphql/types';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';

export enum ContractStatusModalMode {
  Start = 'Start',
  End = 'End',
  Renew = 'Renew',
}
interface ContractModalStatusContextState {
  isModalOpen: boolean;
  nextInvoice?: Invoice | null;
  onStatusModalClose: () => void;
  mode: ContractStatusModalMode | null;
  committedPeriodInMonths: number | string;
  onStatusModalOpen: (mode: ContractStatusModalMode | null) => void;
}

const ContractPanelStateContext =
  createContext<ContractModalStatusContextState>({
    isModalOpen: false,
    onStatusModalOpen: () => null,
    onStatusModalClose: () => null,
    mode: null,
    committedPeriodInMonths: 0,
  });

export const useContractModalStatusContext = () => {
  return useContext(ContractPanelStateContext);
};

export const ContractModalStatusContextProvider = ({
  children,
  id,
  nextInvoice,
  upcomingInvoices,
  committedPeriodInMonths,
}: PropsWithChildren & {
  id: string;
  nextInvoice?: string;
  upcomingInvoices?: Array<Invoice> | null;
  committedPeriodInMonths?: number | string | null;
}) => {
  const [mode, setMode] = useState<ContractStatusModalMode | null>(null);
  const { onOpen, onClose, open } = useDisclosure({
    id: `status-contract-modal-${id}`,
  });

  const onStatusModalOpen = (mode: ContractStatusModalMode | null) => {
    onOpen();
    setMode(mode);
  };
  const onStatusModalClose = () => {
    onClose();
    setMode(null);
  };

  const upcomingInvoice = upcomingInvoices?.find(
    (invoice) => invoice.issued === nextInvoice,
  );

  return (
    <ContractPanelStateContext.Provider
      value={{
        mode,
        isModalOpen: open,
        onStatusModalOpen,
        onStatusModalClose,
        nextInvoice: upcomingInvoice,
        committedPeriodInMonths: committedPeriodInMonths || 0,
      }}
    >
      {children}
    </ContractPanelStateContext.Provider>
  );
};

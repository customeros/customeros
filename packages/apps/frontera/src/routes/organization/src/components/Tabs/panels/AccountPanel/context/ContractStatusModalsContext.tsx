import { useState, useContext, createContext, PropsWithChildren } from 'react';

import { Invoice } from '@graphql/types';
import { useDisclosure } from '@ui/utils/hooks/useDisclosure';

export enum ContractStatusModalMode {
  Start = 'Start',
  End = 'End',
  Renew = 'Renew',
  Delete = 'Delete',
}
interface ContractModalStatusContextState {
  isModalOpen: boolean;
  onStatusModalClose: () => void;
  mode: ContractStatusModalMode | null;
  onStatusModalOpen: (mode: ContractStatusModalMode | null) => void;
}

const ContractPanelStateContext =
  createContext<ContractModalStatusContextState>({
    isModalOpen: false,
    onStatusModalOpen: () => null,
    onStatusModalClose: () => null,
    mode: null,
  });

export const useContractModalStatusContext = () => {
  return useContext(ContractPanelStateContext);
};

export const ContractModalStatusContextProvider = ({
  children,
  id,
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

  return (
    <ContractPanelStateContext.Provider
      value={{
        mode,
        isModalOpen: open,
        onStatusModalOpen,
        onStatusModalClose,
      }}
    >
      {children}
    </ContractPanelStateContext.Provider>
  );
};

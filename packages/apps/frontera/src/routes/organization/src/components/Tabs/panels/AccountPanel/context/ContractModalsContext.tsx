import {
  useState,
  Dispatch,
  useContext,
  createContext,
  SetStateAction,
  PropsWithChildren,
} from 'react';

import { useDisclosure } from '@ui/utils/hooks/useDisclosure';

export enum EditModalMode {
  ContractDetails,
  BillingDetails,
}

interface ContractPanelState {
  isEditModalOpen: boolean;
  onEditModalOpen: () => void;
  onEditModalClose: () => void;
  editModalMode: EditModalMode;
  onChangeModalMode: Dispatch<SetStateAction<EditModalMode>>;
}

const ContractPanelStateContext = createContext<ContractPanelState>({
  isEditModalOpen: false,
  onEditModalOpen: () => null,
  onEditModalClose: () => null,
  onChangeModalMode: () => null,
  editModalMode: EditModalMode.ContractDetails,
});

export const useContractModalStateContext = () => {
  return useContext(ContractPanelStateContext);
};

export const ContractModalsContextProvider = ({
  children,
  id,
}: PropsWithChildren & { id: string }) => {
  const [editModalMode, setEditModalMode] = useState<EditModalMode>(
    EditModalMode.ContractDetails,
  );

  const {
    onOpen: onEditModalOpen,
    onClose: onEditModalClose,
    open: isEditModalOpen,
  } = useDisclosure({
    id: `edit-contract-modal-${id}`,
  });

  return (
    <ContractPanelStateContext.Provider
      value={{
        isEditModalOpen,
        onEditModalOpen,
        onEditModalClose,
        editModalMode,
        onChangeModalMode: setEditModalMode,
      }}
    >
      {children}
    </ContractPanelStateContext.Provider>
  );
};

import {
  useState,
  Dispatch,
  useContext,
  createContext,
  SetStateAction,
  PropsWithChildren,
} from 'react';

import { useDisclosure } from '@ui/utils/hooks/useDisclosure';

interface ServiceLineItemsModalsState {
  isEditModalOpen: boolean;
  onEditModalOpen: () => void;
  onEditModalClose: () => void;
  focusedItemId: string | null;
  onSelectFocusedItem: Dispatch<SetStateAction<string | null>>;
}

const ServiceLineItemsModalsStateContext =
  createContext<ServiceLineItemsModalsState>({
    isEditModalOpen: false,
    onEditModalOpen: () => null,
    onEditModalClose: () => null,
    onSelectFocusedItem: () => null,
    focusedItemId: null,
  });

export const useServiceLineItemsModalsContext = () => {
  return useContext(ServiceLineItemsModalsStateContext);
};

export const ServiceLineItemsModalsContextProvider = ({
  children,
  id,
}: PropsWithChildren & { id: string }) => {
  const [focusedItemId, setFocusedItemId] = useState<null | string>(null);

  const {
    onOpen: onEditModalOpen,
    onClose: onEditModalClose,
    open: isEditModalOpen,
  } = useDisclosure({
    id: `service-line-items-modal-${id}`,
  });

  return (
    <ServiceLineItemsModalsStateContext.Provider
      value={{
        isEditModalOpen,
        onEditModalOpen,
        onEditModalClose,
        focusedItemId,
        onSelectFocusedItem: setFocusedItemId,
      }}
    >
      {children}
    </ServiceLineItemsModalsStateContext.Provider>
  );
};

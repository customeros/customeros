import { useState, useContext, createContext, PropsWithChildren } from 'react';

import {
  useDisclosure,
  UseDisclosureReturn,
} from '@ui/utils/hooks/useDisclosure';

// Moved to upperscope due to error in safari https://linear.app/customer-os/issue/COS-619/scrollbar-overlaps-the-renewal-modals-in-safari

interface ModalContextMethods {
  modal: UseDisclosureReturn;
}
interface AccountPanelState {
  isModalOpen: boolean;
}

const modalDefaultState: UseDisclosureReturn = {
  onClose: () => null,
  onOpen: () => null,
  open: false,
  onToggle: () => null,
  isControlled: false,
  getButtonProps: () => ({}),
  getDisclosureProps: () => ({}),
};

const UpdatePanelModalStateContext = createContext<{
  setIsPanelModalOpen: (newState: boolean) => void;
}>({
  setIsPanelModalOpen: () => null,
});
const ARRInfoModalContext = createContext<ModalContextMethods>({
  modal: modalDefaultState,
});

const UpdateRenewalDetailsContext = createContext<ModalContextMethods>({
  modal: modalDefaultState,
});

const AccountPanelStateContext = createContext<AccountPanelState>({
  isModalOpen: false,
});

export const useUpdatePanelModalStateContext = () => {
  return useContext(UpdatePanelModalStateContext);
};

export const useARRInfoModalContext = () => {
  return useContext(ARRInfoModalContext);
};

export const useAccountPanelStateContext = () => {
  return useContext(AccountPanelStateContext);
};

export const useUpdateRenewalDetailsContext = () => {
  return useContext(UpdateRenewalDetailsContext);
};

export const AccountModalsContextProvider = ({
  children,
}: PropsWithChildren) => {
  const [isPanelModalOpen, setIsPanelModalOpen] = useState<boolean>(false);
  const arrForecastInfoModal = useDisclosure({
    id: 'arr-forecast-info-modal',
  });

  const updateRenewalDetailsModal = useDisclosure({
    id: 'update-renewal-details-modal',
  });

  return (
    <ARRInfoModalContext.Provider
      value={{
        modal: arrForecastInfoModal,
      }}
    >
      <UpdatePanelModalStateContext.Provider
        value={{
          setIsPanelModalOpen,
        }}
      >
        <UpdateRenewalDetailsContext.Provider
          value={{
            modal: updateRenewalDetailsModal,
          }}
        >
          <AccountPanelStateContext.Provider
            value={{
              isModalOpen:
                arrForecastInfoModal.open ||
                updateRenewalDetailsModal.open ||
                isPanelModalOpen,
            }}
          >
            {children}
          </AccountPanelStateContext.Provider>
        </UpdateRenewalDetailsContext.Provider>
      </UpdatePanelModalStateContext.Provider>
    </ARRInfoModalContext.Provider>
  );
};

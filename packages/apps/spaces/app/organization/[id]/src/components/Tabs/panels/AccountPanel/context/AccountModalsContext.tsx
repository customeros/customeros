import { useMemo, useContext, createContext, PropsWithChildren } from 'react';

import { useDisclosure, UseDisclosureReturn } from '@ui/utils';

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
  isOpen: false,
  onToggle: () => null,
  isControlled: false,
  getButtonProps: () => null,
  getDisclosureProps: () => null,
};

const AddServiceModalContext = createContext<ModalContextMethods>({
  modal: modalDefaultState,
});
const UpdateServiceModalContext = createContext<ModalContextMethods>({
  modal: modalDefaultState,
});
const ARRInfoModalContext = createContext<ModalContextMethods>({
  modal: modalDefaultState,
});

const AccountPanelStateContext = createContext<AccountPanelState>({
  isModalOpen: false,
});

export const useAddServiceModalContext = () => {
  return useContext(AddServiceModalContext);
};
export const useUpdateServiceModalContext = () => {
  return useContext(UpdateServiceModalContext);
};
export const useARRInfoModalContext = () => {
  return useContext(ARRInfoModalContext);
};
export const useAccountPanelStateContext = () => {
  return useContext(AccountPanelStateContext);
};

export const AccountModalsContextProvider = ({
  children,
}: PropsWithChildren) => {
  const arrForecastInfoModal = useDisclosure({
    id: 'arr-forecast-info-modal',
  });
  const addServiceModal = useDisclosure({
    id: 'add-service-modal',
  });
  // const addRenewalDetailsModal = useDisclosure({
  //   id: 'add-renewal-details-modal',
  // });

  const updateServiceModal = useDisclosure({
    id: 'update-service-modal',
  });

  const isModalOpen = useMemo(() => {
    return (
      arrForecastInfoModal.isOpen || addServiceModal.isOpen
      // addRenewalDetailsModal.isOpen
    );
  }, [
    arrForecastInfoModal.isOpen,
    addServiceModal.isOpen,
    // addRenewalDetailsModal.isOpen,
  ]);

  return (
    <ARRInfoModalContext.Provider
      value={{
        modal: arrForecastInfoModal,
      }}
    >
      <UpdateServiceModalContext.Provider
        value={{
          modal: updateServiceModal,
        }}
      >
        <AddServiceModalContext.Provider
          value={{
            modal: addServiceModal,
          }}
        >
          <AccountPanelStateContext.Provider
            value={{
              isModalOpen,
            }}
          >
            {children}
          </AccountPanelStateContext.Provider>
        </AddServiceModalContext.Provider>
      </UpdateServiceModalContext.Provider>
    </ARRInfoModalContext.Provider>
  );
};

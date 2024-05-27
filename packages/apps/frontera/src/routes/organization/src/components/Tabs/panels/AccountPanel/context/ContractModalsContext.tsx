import { useForm } from 'react-inverted-form';
import {
  useState,
  Dispatch,
  useContext,
  createContext,
  SetStateAction,
  PropsWithChildren,
} from 'react';

import { useDeepCompareEffect } from 'rooks';

import { useDisclosure } from '@ui/utils/hooks/useDisclosure';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGetContractQuery } from '@organization/graphql/getContract.generated';
import { ContractDetailsDto } from '@organization/components/Tabs/panels/AccountPanel/Contract/ContractBillingDetailsModal/ContractDetails.dto';
import {
  BillingDetailsDto,
  BillingAddressDetailsFormDto,
} from '@organization/components/Tabs/panels/AccountPanel/Contract/BillingAddressDetails/BillingAddressDetailsForm.dto';

export enum EditModalMode {
  ContractDetails,
  BillingDetails,
}

interface ContractPanelState {
  isEditModalOpen: boolean;
  onEditModalOpen: () => void;
  onEditModalClose: () => void;
  editModalMode: EditModalMode;
  detailsState: ContractDetailsDto | null;
  addressState: BillingAddressDetailsFormDto | null;
  onChangeModalMode: Dispatch<SetStateAction<EditModalMode>>;
}

const ContractPanelStateContext = createContext<ContractPanelState>({
  isEditModalOpen: false,
  onEditModalOpen: () => null,
  onEditModalClose: () => null,
  onChangeModalMode: () => null,
  editModalMode: EditModalMode.ContractDetails,
  addressState: null,
  detailsState: null,
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
  const client = getGraphQLClient();

  const {
    onOpen: onEditModalOpen,
    onClose: onEditModalClose,
    open: isEditModalOpen,
  } = useDisclosure({
    id: `edit-contract-modal-${id}`,
  });
  const { data } = useGetContractQuery(
    client,
    {
      id,
    },
    {
      enabled: isEditModalOpen && !!id,
      refetchOnMount: true,
    },
  );
  const defaultValues = new ContractDetailsDto(data?.contract);
  const formId = `billing-details-form-${id}`;

  const { state, setDefaultValues } = useForm<ContractDetailsDto>({
    formId,
    defaultValues,
    stateReducer: (_, _action, next) => {
      return next;
    },
  });

  const addressDetailsDefaultValues = new BillingDetailsDto(data?.contract);

  const { state: addressState, setDefaultValues: setDefaultAddressValues } =
    useForm<BillingAddressDetailsFormDto>({
      formId: 'billing-details-address-form',
      defaultValues: addressDetailsDefaultValues,
      stateReducer: (_, _action, next) => {
        return next;
      },
    });

  useDeepCompareEffect(() => {
    setDefaultValues(defaultValues);
  }, [defaultValues]);
  useDeepCompareEffect(() => {
    setDefaultAddressValues(addressDetailsDefaultValues);
  }, [addressDetailsDefaultValues]);

  return (
    <ContractPanelStateContext.Provider
      value={{
        isEditModalOpen,
        onEditModalOpen,
        onEditModalClose,
        editModalMode,
        onChangeModalMode: setEditModalMode,
        addressState: addressState.values,
        detailsState: state.values,
      }}
    >
      {children}
    </ContractPanelStateContext.Provider>
  );
};

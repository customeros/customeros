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
import { useGetContractQuery } from '@organization/src/graphql/getContract.generated';
import { ContractDetailsDto } from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/ContractBillingDetailsModal/ContractDetails.dto';
import {
  BillingDetailsDto,
  BillingAddressDetailsFormDto,
} from '@organization/src/components/Tabs/panels/AccountPanel/ContractNew/BillingAddressDetails/BillingAddressDetailsForm.dto';

export enum EditModalMode {
  ContractDetails,
  BillingDetails,
  MakeLive,
  RenewContract,
  PauseInvoicing,
  ResumeInvoicing,
}

interface ContractPanelState {
  addressState: any;
  detailsState: any;
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
  addressState: {},
  detailsState: {},
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

  const { state, setDefaultValues } = useForm({
    formId,
    defaultValues,
    stateReducer: (_, action, next) => {
      return next;
    },
  });

  const addressDetailsDefailtValues = new BillingDetailsDto(data?.contract);

  const { state: addressState, setDefaultValues: setDefaultAddressValues } =
    useForm<BillingAddressDetailsFormDto>({
      formId: 'billing-details-address-form',
      defaultValues: addressDetailsDefailtValues,
      stateReducer: (_, action, next) => {
        return next;
      },
    });

  useDeepCompareEffect(() => {
    setDefaultValues(defaultValues);
  }, [defaultValues]);
  useDeepCompareEffect(() => {
    setDefaultAddressValues(addressDetailsDefailtValues);
  }, [addressDetailsDefailtValues]);

  return (
    <ContractPanelStateContext.Provider
      value={{
        isEditModalOpen,
        onEditModalOpen,
        onEditModalClose,
        editModalMode,
        onChangeModalMode: setEditModalMode,
        addressState,
        detailsState: state,
      }}
    >
      {children}
    </ContractPanelStateContext.Provider>
  );
};

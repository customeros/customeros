import { toast } from 'react-toastify';
import {
  RemoveHealthIndicatorMutationVariables,
  SetHealthIndicatorMutationVariables,
  useRemoveHealthIndicatorMutation,
  useSetHealthIndicatorMutation,
} from '@spaces/graphql';

interface Result {
  saving: boolean;
  onSetHealthIndicator: ({
    variables,
  }: {
    variables: SetHealthIndicatorMutationVariables;
  }) => void;
  onRemoveHealthIndicator: ({
    variables,
  }: {
    variables: RemoveHealthIndicatorMutationVariables;
  }) => void;
}
export const useSetHealthIndicator = (): Result => {
  const [removeHealthIndicator, { loading: savingRemove }] =
    useRemoveHealthIndicatorMutation({
      onError: () => {
        toast.error('Something went wrong while setting health indicator');
      },
    });

  const [setHealthIndicator, { loading: savingSetNew }] =
    useSetHealthIndicatorMutation({
      onError: () => {
        toast.error('Something went wrong while setting health indicator');
      },
    });

  return {
    onSetHealthIndicator: setHealthIndicator,
    onRemoveHealthIndicator: removeHealthIndicator,
    saving: savingRemove || savingSetNew,
  };
};

import { QueryKey } from '@tanstack/react-query';

import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import {
  User,
  Action,
  ActionType,
  BilledType,
  DataSource,
  InputMaybe,
  ServiceLineItem,
  ServiceLineItemUpdateInput,
  ServiceLineItemBulkUpdateItem,
} from '@graphql/types';

type UpdateInput = {
  user: string;
  content: string;
  metadata: string;
  actionType: string;
};
const update = ({ content, metadata, actionType, user }: UpdateInput) => ({
  __typename: 'Action',
  id: Math.random().toString(),
  createdAt: new Date(),
  updatedAt: '',
  actionType,
  appSource: 'customeros-optimistic-update',
  source: 'customeros-optimistic-update',
  metadata,
  actionCreatedBy: {
    firstName: user,
    lastName: '',
  },
  content,
});

export const getUpdateServiceEvents = (
  prev?: ServiceLineItem,
  next?: ServiceLineItemBulkUpdateItem | null,
  user?: string,
  // @ts-expect-error fixme later
  updateTimelineCache: (event: Action) => void,
) => {
  const metadata = JSON.stringify({
    price: next?.price,
    previousPrice: prev?.price,
    billedType: next?.billed,
  });

  if (prev?.price !== next?.price) {
    const decreased =
      parseFloat(`${prev?.price}`) > parseFloat(`${next?.price}`);

    const event = update({
      metadata,
      user: user ?? '',
      actionType: ActionType.ServiceLineItemPriceUpdated,
      content: `${user} ${
        next?.isRetroactiveCorrection ? 'retroactively ' : ''
      }${decreased ? 'decreased' : 'increased'} the price for ${next?.name}`,
    });
    updateTimelineCache(event as Action);
  }

  if (prev?.billingCycle !== next?.billed) {
    const event = update({
      metadata,
      user: user ?? '',
      actionType: ActionType.ServiceLineItemBilledTypeUpdated,
      content: `${user} ${
        next?.isRetroactiveCorrection ? 'retroactively ' : ''
      } changed the billing cycle for ${next?.name}`,
    });
    updateTimelineCache(event as Action);
  }

  if (prev?.quantity !== next?.quantity) {
    const decreased = parseFloat(prev?.quantity) > parseFloat(next?.quantity);
    const event = update({
      metadata,
      user: user ?? '',
      actionType: ActionType.ServiceLineItemQuantityUpdated,
      content: `${user} ${
        next?.isRetroactiveCorrection ? 'retroactively ' : ''
      }${decreased ? 'decreased' : 'increased'} the quantity of ${next?.name}`,
    });
    updateTimelineCache(event as Action);
  }
};

export const updateTimelineCacheWithNewServiceLineItem = ({
  input,
  user,
  updateTimelineCache,
  contractName,
  timelineQueryKey,
}: {
  user: string;
  contractName: string;
  timelineQueryKey: QueryKey;
  input: InputMaybe<ServiceLineItemBulkUpdateItem>;
  updateTimelineCache: (
    event: Action & { actionCreatedBy: Pick<User, 'firstName' | 'lastName'> },
    queryKey: QueryKey,
  ) => void;
}) => {
  const isRecurring = [
    BilledType.Annually,
    BilledType.Monthly,
    BilledType.Quarterly,
  ].includes(input?.billed as BilledType);
  const metadata = JSON.stringify({
    price: input?.price,
    billedType: input?.billed,
  });
  const actionType = isRecurring
    ? ActionType.ServiceLineItemBilledTypeRecurringCreated
    : input?.billed === BilledType.Usage
    ? ActionType.ServiceLineItemBilledTypeUsageCreated
    : ActionType.ServiceLineItemBilledTypeOnceCreated;
  updateTimelineCache(
    {
      __typename: 'Action',
      id: Math.random().toString(),
      createdAt: new Date(),
      actionType,
      appSource: 'customeros-optimistic-update',
      source: DataSource.Openline,
      metadata,
      actionCreatedBy: {
        firstName: user,
        lastName: '',
      },
      content: `${user} added a ${
        isRecurring
          ? 'recurring'
          : input?.billed === BilledType.Usage
          ? 'use based'
          : 'one-time'
      } service to ${contractName}: ${input?.name} , at ${formatCurrency(
        input?.price ?? 0,
      )}`,
    },
    timelineQueryKey,
  );
};

export const updateTimelineCacheAfterServiceLineItemChange = ({
  prevServiceLineItems,
  newServiceLineItems,
  user,
  updateTimelineCache,
  timelineQueryKey,
  contractName,
}: {
  user: string;
  contractName: string;
  timelineQueryKey: QueryKey;
  prevServiceLineItems: Array<ServiceLineItem>;
  updateTimelineCache: (event: Action, queryKey: QueryKey) => void;
  newServiceLineItems: InputMaybe<ServiceLineItemBulkUpdateItem>[];
}) => {
  // const deletedServiceLineItems = prevServiceLineItems.filter((item) => {
  //   const x =
  //     newServiceLineItems.findIndex(
  //       (e) =>
  //         e.serviceLineItemId === item.id ||
  //         e.serviceLineItemId === item.parentId,
  //     ) === -1;
  //
  //   return !item.endedAt && x;
  // });

  const newItems = newServiceLineItems.filter(
    (element) => !element?.serviceLineItemId,
  );

  if (newItems.length) {
    newItems.forEach((element) =>
      updateTimelineCacheWithNewServiceLineItem({
        contractName,
        input: element,
        timelineQueryKey,
        user,
        updateTimelineCache,
      }),
    );
  }

  const changedItems = newServiceLineItems.filter((obj1) =>
    prevServiceLineItems.some(
      (obj2) =>
        (obj1?.serviceLineItemId === obj2.metadata.id ||
          obj2.parentId === obj1?.serviceLineItemId) &&
        (obj2.description !== obj1?.name ||
          obj1?.quantity !== obj2.quantity ||
          obj1?.price !== obj2.price ||
          obj1?.billed !== obj2.billingCycle),
    ),
  );

  if (changedItems.length) {
    changedItems.forEach((element) => {
      getUpdateServiceEvents(
        prevServiceLineItems.find(
          (e) => e.metadata.id === element?.serviceLineItemId,
        ),
        element as ServiceLineItemUpdateInput,
        user,
        (event: Action) => updateTimelineCache(event, timelineQueryKey),
      );
    });
  }
};

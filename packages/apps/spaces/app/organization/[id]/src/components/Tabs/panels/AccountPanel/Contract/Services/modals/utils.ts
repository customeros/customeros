import {
  Action,
  ActionType,
  ServiceLineItem,
  ServiceLineItemUpdateInput,
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
  prev: ServiceLineItem,
  next: ServiceLineItemUpdateInput,
  user: string,
  updateTimelineCache: (event: Action) => void,
) => {
  const metadata = JSON.stringify({
    price: next.price,
    previousPrice: prev.price,
    billedType: next.billed,
  });

  if (prev?.price !== next?.price) {
    const decreased = parseFloat(`${prev.price}`) > parseFloat(`${next.price}`);

    const event = update({
      metadata,
      user,
      actionType: ActionType.ServiceLineItemPriceUpdated,
      content: `${user} ${
        next.isRetroactiveCorrection ? 'retroactively ' : ''
      }${decreased ? 'decreased' : 'increased'} the price for ${next.name}`,
    });
    updateTimelineCache(event as Action);
  }

  if (prev?.billed !== next?.billed) {
    const event = update({
      metadata,
      user,
      actionType: ActionType.ServiceLineItemBilledTypeUpdated,
      content: `${user} ${
        next.isRetroactiveCorrection ? 'retroactively ' : ''
      } changed the billing cycle for ${next.name}`,
    });
    updateTimelineCache(event as Action);
  }

  if (prev?.quantity !== next?.quantity) {
    const decreased = parseFloat(prev.quantity) > parseFloat(next.quantity);
    const event = update({
      metadata,
      user,
      actionType: ActionType.ServiceLineItemQuantityUpdated,
      content: `${user} ${
        next.isRetroactiveCorrection ? 'retroactively ' : ''
      }${decreased ? 'decreased' : 'increased'} the quantity of ${next.name}`,
    });
    updateTimelineCache(event as Action);
  }
};

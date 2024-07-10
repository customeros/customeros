import {
  BilledType,
  ServiceLineItem,
} from '@shared/types/__generated__/graphql.types';

export const groupServicesByParentId = (services: ServiceLineItem[]) => {
  const { subscription, once } = services.reduce<{
    once: ServiceLineItem[];
    subscription: ServiceLineItem[];
  }>(
    (acc, item) => {
      const isSubscription = [
        BilledType.Monthly,
        BilledType.Quarterly,
        BilledType.Annually,
      ].includes(item.billingCycle);

      const key: 'subscription' | 'once' = isSubscription
        ? 'subscription'
        : 'once';
      acc[key].push(item);

      return acc;
    },
    { subscription: [], once: [] },
  );

  const groupByParentId = (
    services: ServiceLineItem[],
  ): Record<string, ServiceLineItem[]> => {
    return services.reduce<Record<string, ServiceLineItem[]>>(
      (acc, service) => {
        const parentId = service?.parentId || service?.metadata?.id;
        if (parentId) {
          if (!acc[parentId]) {
            acc[parentId] = [];
          }
          acc[parentId].push(service);
        }

        return acc;
      },
      {},
    );
  };

  const sortByServiceStarted = (
    group: ServiceLineItem[],
  ): ServiceLineItem[] => {
    return group.sort(
      (a, b) =>
        new Date(a?.serviceStarted).getTime() -
        new Date(b?.serviceStarted).getTime(),
    );
  };

  const filterGroups = (groups: ServiceLineItem[][]): ServiceLineItem[][] => {
    return groups.filter((group) =>
      group.some((service) => service?.serviceEnded === null),
    );
  };

  const getCurrentLineItem = (group: ServiceLineItem[]): ServiceLineItem => {
    const today = new Date();

    return (
      group.reduce<ServiceLineItem | null>((currentService, service) => {
        const serviceStarted = new Date(service.serviceStarted);
        const serviceEnded = service.serviceEnded
          ? new Date(service.serviceEnded)
          : null;

        const isActiveToday =
          (serviceStarted <= today &&
            (!serviceEnded || serviceEnded > today)) ||
          (serviceStarted.toDateString() === today.toDateString() &&
            (!serviceEnded || serviceEnded > today));

        if (currentService) {
          const currentStarted = new Date(currentService.serviceStarted);

          if (isActiveToday) {
            return serviceStarted > currentStarted ? service : currentService;
          }

          if (!isActiveToday && serviceStarted <= today) {
            return serviceStarted > currentStarted ? service : currentService;
          }
        } else if (isActiveToday) {
          return service;
        }

        return currentService;
      }, null as ServiceLineItem | null) || group[0]
    );
  };

  const processGroups = (
    services: ServiceLineItem[],
  ): { group: ServiceLineItem[]; currentLineItem: ServiceLineItem }[] => {
    const groupedServices = groupByParentId(services);
    const sortedGroups =
      Object.values(groupedServices).map(sortByServiceStarted);
    const filteredGroups = filterGroups(sortedGroups);

    return filteredGroups.map((group) => ({
      group,
      currentLineItem: getCurrentLineItem(group),
    }));
  };

  return {
    subscription: processGroups(subscription),
    once: processGroups(once),
  };
};

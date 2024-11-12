export const organizationsTestState = {
  createdOrganizationIds: new Set<string>(),
};

export const trackOrganization = (organizationId: string) => {
  organizationsTestState.createdOrganizationIds.add(organizationId);
};

export const getCreatedOrganizationIds = () => {
  return Array.from(organizationsTestState.createdOrganizationIds);
};

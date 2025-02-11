export const contactsTestState = {
  createdContactsIds: new Set<string>(),
};

export const trackContact = (organizationId: string) => {
  contactsTestState.createdContactsIds.add(organizationId);
};

export const getCreatedContact = () => {
  return Array.from(contactsTestState.createdContactsIds);
};

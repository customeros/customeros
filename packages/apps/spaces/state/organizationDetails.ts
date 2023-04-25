import { atom } from 'recoil';

interface OrganizationDetailsEdit {
  isEditMode: boolean;
}

interface OrganizationNewItemsToEdit {
  timelineEvents: Array<{ id: string }>;
}
export const organizationDetailsEdit = atom<OrganizationDetailsEdit>({
  key: 'organizationDetailsEdit', // unique ID (with respect to other atoms/selectors)
  default: {
    isEditMode: false,
  },
});

export const organizationNewItemsToEdit = atom<OrganizationNewItemsToEdit>({
  key: 'organizationNewItemsToEdit', // unique ID (with respect to other atoms/selectors)
  default: {
    timelineEvents: [],
  },
});

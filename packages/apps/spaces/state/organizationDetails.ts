import { atom } from 'recoil';

interface ContactDetailsEdit {
  isEditMode: boolean;
}

export const organizationDetailsEdit = atom<ContactDetailsEdit>({
  key: 'organizationDetailsEdit', // unique ID (with respect to other atoms/selectors)
  default: {
    isEditMode: false,
  },
});

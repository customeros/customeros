import { atom } from 'recoil';

interface ContactDetailsEdit {
  isEditMode: boolean;
}

export const contactDetailsEdit = atom<ContactDetailsEdit>({
  key: 'contactDetailsEdit', // unique ID (with respect to other atoms/selectors)
  default: {
    isEditMode: false,
  },
});

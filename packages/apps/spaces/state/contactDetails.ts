import { atom } from 'recoil';

interface ContactDetailsEdit {
  isEditMode: boolean;
}
interface ContactNewItemsToEdit {
  timelineEvents: Array<{ id: string }>;
}

export const contactDetailsEdit = atom<ContactDetailsEdit>({
  key: 'contactDetailsEdit', // unique ID (with respect to other atoms/selectors)
  default: {
    isEditMode: false,
  },
});

export const contactNewItemsToEdit = atom<ContactNewItemsToEdit>({
  key: 'contactNewItemsToEdit', // unique ID (with respect to other atoms/selectors)
  default: {
    timelineEvents: [],
  },
});

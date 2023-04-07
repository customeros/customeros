import { DeleteContactMutation, useDeleteContactMutation } from './types';
import { useRouter } from 'next/router';
import { toast } from 'react-toastify';

interface Props {
  id: string;
}

interface Result {
  onArchiveContact: () => Promise<
    DeleteContactMutation['contact_SoftDelete'] | null
  >;
}
export const useArchiveContact = ({ id }: Props): Result => {
  const { push } = useRouter();
  const [archiveContactMutation, { loading, error, data }] =
    useDeleteContactMutation({
      variables: {
        id,
      },
    });
  const handleArchiveContact: Result['onArchiveContact'] = async () => {
    try {
      const response = await archiveContactMutation();
      if (response) {
        push('/').then(() =>
          toast.success('Contact successfully archived!', {}),
        );
      }

      return response.data?.contact_SoftDelete ?? null;
    } catch (err) {
      toast.error('Something went wrong while archiving contact!', {
        toastId: `archive-contact-error-${id}`,
      });
      return null;
    }
  };

  return {
    onArchiveContact: handleArchiveContact,
  };
};

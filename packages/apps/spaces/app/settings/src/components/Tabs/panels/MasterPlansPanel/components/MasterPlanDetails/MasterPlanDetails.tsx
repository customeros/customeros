import { useMemo, useEffect } from 'react';
import { useForm } from 'react-inverted-form';

import { Flex } from '@ui/layout/Flex';
import { FormInput } from '@ui/form/Input';
import { IconButton } from '@ui/form/IconButton';
import { DotsVertical } from '@ui/media/icons/DotsVertical';

interface MasterPlanDetailsProps {
  name: string;
}

type MasterPlanForm = {
  name: string;
};

const formId = 'master-plan-details-form';

export const MasterPlanDetails = ({ name }: MasterPlanDetailsProps) => {
  const defaultValues = useMemo<MasterPlanForm>(() => ({ name }), [name]);

  const { setDefaultValues } = useForm<MasterPlanForm>({
    formId,
    defaultValues,
  });

  useEffect(() => {
    setDefaultValues(defaultValues);
  }, [name]);

  return (
    <>
      <Flex align='center' justify='space-between' mb='2'>
        <FormInput
          name='name'
          formId={formId}
          variant='unstyled'
          borderRadius='unset'
          fontWeight='semibold'
        />
        <IconButton
          size='xs'
          variant='ghost'
          aria-label='Master Plan Options'
          icon={<DotsVertical color='gray.400' />}
        />
      </Flex>
    </>
  );
};

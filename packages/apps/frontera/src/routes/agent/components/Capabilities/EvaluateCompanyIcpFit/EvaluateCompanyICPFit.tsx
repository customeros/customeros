import { useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { EditIcpQualificationCriteriaUsecase } from '@domain/usecases/agents/capabilities/edit-icp-qualification-criteria.usecase';
import { EditIcpDisqualificationCriteriaUsecase } from '@domain/usecases/agents/capabilities/edit-icp-disqualification-criteria.usecase';

import { Icon } from '@ui/media/Icon';
import { useStore } from '@shared/hooks/useStore';
import { Textarea } from '@ui/form/Textarea/Textarea';
import {
  ScrollAreaRoot,
  ScrollAreaThumb,
  ScrollAreaViewport,
  ScrollAreaScrollbar,
} from '@ui/utils/ScrollArea';

const disqualificationCriteriaUsecase =
  new EditIcpDisqualificationCriteriaUsecase();
const qualificationCriteriaUsecase = new EditIcpQualificationCriteriaUsecase();
export const EvaluateCompanyIcpFit = observer(() => {
  const { id } = useParams<{ id: string }>();
  const store = useStore();

  useEffect(() => {
    if (!id || store.agents.value.size === 0) return;
    disqualificationCriteriaUsecase.setAgentId(id);
    qualificationCriteriaUsecase.setAgentId(id);
    disqualificationCriteriaUsecase.init();
    qualificationCriteriaUsecase.init();
  }, [id, store.agents.value.size]);

  return (
    <article className='flex flex-col gap-4 -mr-4'>
      <h1 className='text-sm font-medium pr-4'>Evaluate company ICP fit</h1>
      <ScrollAreaRoot>
        <ScrollAreaViewport className=''>
          <div className=' h-[90vh] pr-4'>
            <div>
              <h2 className='text-sm font-medium mb-1'>
                Qualification criteria
              </h2>
              <p className='text-sm'>
                What criteria define an ideal customer? You can list critical
                qualifiers like industry types or describe key qualifications,
                such as specific regions, company size, or founding date.
              </p>

              <Textarea
                variant='outline'
                value={qualificationCriteriaUsecase.inputValue}
                onBlur={() => qualificationCriteriaUsecase.execute()}
                placeholder='Enter keywords or key qualifications....'
                className='mt-3 px-2 py-1 text-sm bg-transparent resize-none min-h-[72px]  overflow-y-hidden'
                onChange={(e) => {
                  qualificationCriteriaUsecase.setInputValue(e.target.value);

                  if (qualificationCriteriaUsecase.validationError) {
                    qualificationCriteriaUsecase.validate();
                  }
                }}
              />
              {(qualificationCriteriaUsecase.capabilityErrors ||
                qualificationCriteriaUsecase.validationError) && (
                <div className='bg-error-50 text-error-700 px-2 py-1 rounded-[4px] '>
                  <Icon stroke='none' className='mr-2' name='dot-single' />
                  <span className='text-sm'>
                    {qualificationCriteriaUsecase.capabilityErrors}
                    {qualificationCriteriaUsecase.validationError}
                  </span>
                </div>
              )}
            </div>

            {/*<div className='flex flex-col'>*/}
            {/*  <p className='text-sm font-medium'>Add 5 ideal customers</p>*/}
            {/*  <p className='text-sm'>*/}
            {/*    Add at least 5 company websites that match your ideal customer*/}
            {/*    profile. Feel free to add more.*/}
            {/*  </p>*/}

            {/*  <Button*/}
            {/*    size='xs'*/}
            {/*    variant='ghost'*/}
            {/*    colorScheme='primary'*/}
            {/*    // onClick={usecase.open}*/}
            {/*    className='w-fit mt-1'*/}
            {/*    leftIcon={<Icon name='plus-circle' />}*/}
            {/*  >*/}
            {/*    Add website*/}
            {/*  </Button>*/}
            {/*</div>*/}

            <div>
              <h2 className='text-sm font-medium mb-1'>
                Disqualification criteria
              </h2>
              <p className='text-sm'>
                What criteria should disqualify a company from being an ideal
                customer? You can list keywords like industry types or describe
                detailed deal-breakers, such as specific regions, company size,
                or founding date.
              </p>

              <Textarea
                variant='outline'
                value={disqualificationCriteriaUsecase.inputValue}
                onBlur={() => disqualificationCriteriaUsecase.execute()}
                placeholder={'Enter keywords or detailed dealbreakers...'}
                className='mt-3 px-2 py-1 text-sm bg-transparent resize-none min-h-[72px]  overflow-y-hidden'
                onChange={(e) => {
                  disqualificationCriteriaUsecase.setInputValue(e.target.value);

                  if (disqualificationCriteriaUsecase.validationError) {
                    disqualificationCriteriaUsecase.validate();
                  }
                }}
              />
              {(disqualificationCriteriaUsecase.capabilityErrors ||
                disqualificationCriteriaUsecase.validationError) && (
                <div className='bg-error-50 text-error-700 px-2 py-1 rounded-[4px] '>
                  <Icon stroke='none' className='mr-2' name='dot-single' />
                  <span className='text-sm'>
                    {disqualificationCriteriaUsecase.capabilityErrors}
                    {disqualificationCriteriaUsecase.validationError}
                  </span>
                </div>
              )}
            </div>
          </div>
        </ScrollAreaViewport>
        <ScrollAreaScrollbar orientation='vertical'>
          <ScrollAreaThumb />
        </ScrollAreaScrollbar>
      </ScrollAreaRoot>
    </article>
  );
});

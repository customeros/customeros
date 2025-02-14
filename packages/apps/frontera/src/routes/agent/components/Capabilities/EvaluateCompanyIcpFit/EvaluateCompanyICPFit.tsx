import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { EditIcpDomainsUsecase } from '@domain/usecases/agents/capabilities/edit-icp-domains.usecase';
import { EditIcpQualificationCriteriaUsecase } from '@domain/usecases/agents/capabilities/edit-icp-qualification-criteria.usecase';
import { EditIcpDisqualificationCriteriaUsecase } from '@domain/usecases/agents/capabilities/edit-icp-disqualification-criteria.usecase';

import { Icon } from '@ui/media/Icon';
import { Button } from '@ui/form/Button/Button';
import { IconButton } from '@ui/form/IconButton';
import { Textarea } from '@ui/form/Textarea/Textarea';
import { getFormattedLink } from '@utils/getExternalLink.ts';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import {
  ScrollAreaRoot,
  ScrollAreaThumb,
  ScrollAreaViewport,
  ScrollAreaScrollbar,
} from '@ui/utils/ScrollArea';

import { IdealCustomersModal } from './IdealCustomersModal';

const disqualificationCriteriaUsecase =
  new EditIcpDisqualificationCriteriaUsecase();
const qualificationCriteriaUsecase = new EditIcpQualificationCriteriaUsecase();
const editIcpDomainsUsecase = new EditIcpDomainsUsecase();

export const EvaluateCompanyIcpFit = observer(() => {
  const { id } = useParams<{ id: string }>();
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    if (!id) return;
    disqualificationCriteriaUsecase.init(id);
    qualificationCriteriaUsecase.init(id);
    editIcpDomainsUsecase.init(id);
  }, [id]);

  return (
    <article className='flex flex-col gap-4 -mr-4'>
      <h1 className='text-sm font-medium pr-4'>Evaluate company ICP fit</h1>
      <ScrollAreaRoot>
        <ScrollAreaViewport className=''>
          <div className='h-[90vh] pr-4 flex flex-col gap-4'>
            <div className='flex flex-col'>
              <h2 className='text-sm font-medium mb-1'>
                Add 5 or more of your ideal customers
              </h2>
              <p className='text-sm'>
                Add 5 or more websites of companies that match your ideal
                customer profile
              </p>

              {editIcpDomainsUsecase.icpCompanyExamples.size > 0 && (
                <ul className='mt-1 flex flex-col gap-y-1'>
                  {Array.from(editIcpDomainsUsecase.icpCompanyExamples).map(
                    (website) => (
                      <li
                        key={`ideal-icp${website}`}
                        className='flex items-center gap-x-2 text-sm mx-2 group'
                      >
                        <Icon
                          stroke='none'
                          name='dot-single'
                          className={'text-grayModern-500'}
                        />
                        {getFormattedLink(website)}

                        <Menu modal={false}>
                          <MenuButton asChild>
                            <IconButton
                              size='xxs'
                              variant='ghost'
                              aria-label='more'
                              icon={<Icon name='dots-vertical' />}
                              className='ml-2 invisible group-hover:visible'
                            />
                          </MenuButton>
                          <MenuList>
                            <MenuItem
                              onClick={() => {
                                editIcpDomainsUsecase.select(website);
                              }}
                            >
                              <Icon
                                name='x-circle'
                                className='text-grayModern-500'
                              />
                              Remove
                            </MenuItem>
                          </MenuList>
                        </Menu>
                      </li>
                    ),
                  )}
                </ul>
              )}

              <Button
                size='xs'
                variant='ghost'
                colorScheme='primary'
                className='w-fit mt-1'
                onClick={() => setIsOpen(true)}
                leftIcon={<Icon name='plus-circle' />}
              >
                Add website
              </Button>

              {editIcpDomainsUsecase.icpCompanyExamplesError && (
                <p className='text-sm text-error-500 mt-1'>
                  {editIcpDomainsUsecase.icpCompanyExamplesError}
                </p>
              )}
            </div>

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
                }}
              />
              {qualificationCriteriaUsecase.capabilityErrors && (
                <div className='bg-error-50 text-error-700 px-2 py-1 rounded-[4px] '>
                  <Icon stroke='none' className='mr-2' name='dot-single' />
                  <span className='text-sm'>
                    {qualificationCriteriaUsecase.capabilityErrors}
                  </span>
                </div>
              )}
            </div>

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
                }}
              />
              {disqualificationCriteriaUsecase.capabilityErrors && (
                <div className='bg-error-50 text-error-700 px-2 py-1 rounded-[4px] '>
                  <Icon stroke='none' className='mr-2' name='dot-single' />
                  <span className='text-sm'>
                    {disqualificationCriteriaUsecase.capabilityErrors}
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
      <IdealCustomersModal
        isOpen={isOpen}
        usecase={editIcpDomainsUsecase}
        onClose={() => setIsOpen(false)}
      />
    </article>
  );
});

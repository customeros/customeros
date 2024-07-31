import { FC } from 'react';

import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';
import { IssueBgPattern } from '@ui/media/logos/IssueBgPattern';
import { IssueWithAliases } from '@organization/components/Timeline/types';
import { MarkdownContentRenderer } from '@ui/presentation/MarkdownContentRenderer/MarkdownContentRenderer';
import { useTimelineEventPreviewMethodsContext } from '@organization/components/Timeline/shared/TimelineEventPreview/context/TimelineEventPreviewContext';

function getStatusColor(status: string) {
  if (status === 'solved' || status === 'closed') {
    return 'gray';
  }

  return 'blue';
}

export const IssueStub: FC<{ data: IssueWithAliases }> = ({ data }) => {
  const { openModal } = useTimelineEventPreviewMethodsContext();
  const statusColorScheme = getStatusColor(data.issueStatus);

  return (
    <div className='flex cursor-pointer !imortant hover:transition-all hover:filter hover:drop-shadow-[0_2px_2px_rgba(16,24,40,0.09)] w-[510px] h-[110px] relative hover:duration-200 hover:ease-oute'>
      <div
        className='left-[-1px] top-[-1px] w-[502px] h-[110px] bg-gray-200 absolute'
        style={{
          clipPath:
            'polygon(0% 5.554%,0% 5.554%,0.016% 4.653%,0.061% 3.798%,0.134% 3.001%,0.232% 2.274%,0.351% 1.626%,0.491% 1.071%,0.648% 0.62%,0.821% 0.283%,1.005% 0.073%,1.2% 0%,84.4% 0%,84.4% 0%,84.417% 0.901%,84.467% 1.756%,84.545% 2.553%,84.651% 3.28%,84.781% 3.928%,84.932% 4.483%,85.103% 4.934%,85.289% 5.271%,85.489% 5.481%,85.7% 5.554%,85.7% 5.554%,85.911% 5.481%,86.111% 5.271%,86.298% 4.934%,86.468% 4.483%,86.619% 3.928%,86.749% 3.28%,86.855% 2.553%,86.934% 1.756%,86.983% 0.901%,87% 0%,98.8% 0%,98.8% 0%,98.995% 0.073%,99.179% 0.283%,99.352% 0.62%,99.509% 1.071%,99.649% 1.626%,99.768% 2.274%,99.866% 3.001%,99.939% 3.798%,99.984% 4.653%,100% 5.554%,100% 94.446%,100% 94.446%,99.984% 95.347%,99.939% 96.202%,99.866% 96.999%,99.768% 97.726%,99.649% 98.374%,99.509% 98.929%,99.352% 99.38%,99.179% 99.717%,98.995% 99.927%,98.8% 100%,87% 100%,87% 100%,86.983% 99.099%,86.934% 98.244%,86.855% 97.447%,86.749% 96.72%,86.619% 96.072%,86.468% 95.517%,86.298% 95.066%,86.111% 94.729%,85.911% 94.519%,85.7% 94.446%,85.7% 94.446%,85.489% 94.519%,85.289% 94.729%,85.103% 95.066%,84.932% 95.517%,84.781% 96.072%,84.651% 96.72%,84.545% 97.447%,84.467% 98.244%,84.417% 99.099%,84.4% 100%,1.2% 100%,1.2% 100%,1.005% 99.927%,0.821% 99.717%,0.648% 99.38%,0.491% 98.929%,0.351% 98.374%,0.232% 97.726%,0.134% 96.999%,0.061% 96.202%,0.016% 95.347%,0% 94.446%,0% 5.554%)',
        }}
      />
      <div
        onClick={() => openModal(data.id)}
        className='flex w-[500px] h-[108px] bg-white relative transition-all duration-200 ease out text-sm'
        style={{
          clipPath:
            'polygon(0px 6px, 0px 6px, 0.07852983px 5.02676847px, 0.30588384px 4.10353536px, 0.66970881px 3.24265389px, 1.15765152px 2.45647728px, 1.75735875px 1.75735875px, 2.45647728px 1.15765152px, 3.24265389px 0.66970881px, 4.10353536px 0.30588384px, 5.02676847px 0.07852983px, 6px 9.9333241925913E-32px, 422px 0px, 422px 0px, 422.08507px 0.97323153px, 422.33136px 1.89646464px, 422.72549px 2.75734611px, 423.25408px 3.54352272px, 423.90375px 4.24264125px, 424.66112px 4.84234848px, 425.51281px 5.33029119px, 426.44544px 5.69411616px, 427.44563px 5.92147017px, 428.5px 6px, 428.5px 6px, 429.55437px 5.92147017px, 430.55456px 5.69411616px, 431.48719px 5.33029119px, 432.33888px 4.84234848px, 433.09625px 4.24264125px, 433.74592px 3.54352272px, 434.27451px 2.75734611px, 434.66864px 1.89646464px, 434.91493px 0.97323153px, 435px 1.1036871416792E-15px, 494px 0px, 494px 0px, 494.973302px 0.07852983px, 495.896576px 0.30588384px, 496.757474px 0.66970881px, 497.543648px 1.15765152px, 498.24275px 1.75735875px, 498.842432px 2.45647728px, 499.330346px 3.24265389px, 499.694144px 4.10353536px, 499.921478px 5.02676847px, 500px 6px, 500px 102px, 500px 102px, 499.921478px 102.973302px, 499.694144px 103.896576px, 499.330346px 104.757474px, 498.842432px 105.543648px, 498.24275px 106.24275px, 497.543648px 106.842432px, 496.757474px 107.330346px, 495.896576px 107.694144px, 494.973302px 107.921478px, 494px 108px, 435px 108px, 435px 108px, 434.91493px 107.026698px, 434.66864px 106.103424px, 434.27451px 105.242526px, 433.74592px 104.456352px, 433.09625px 103.75725px, 432.33888px 103.157568px, 431.48719px 102.669654px, 430.55456px 102.305856px, 429.55437px 102.078522px, 428.5px 102px, 428.5px 102px, 427.44563px 102.078522px, 426.44544px 102.305856px, 425.51281px 102.669654px, 424.66112px 103.157568px, 423.90375px 103.75725px, 423.25408px 104.456352px, 422.72549px 105.242526px, 422.33136px 106.103424px, 422.08507px 107.026698px, 422px 108px, 6.00001px 108px, 6.00001px 108px, 5.02677819px 107.921478px, 4.10354432px 107.694144px, 3.24266173px 107.330346px, 2.45648376px 106.842432px, 1.75736375px 106.24275px, 1.15765504px 105.543648px, 0.66971097px 104.757474px, 0.30588488px 103.896576px, 0.07853011px 102.973302px, 9.9333611704463E-32px 102px, 0px 6px)',
        }}
      >
        <div className='shadow-xs pr-2 p-3 flex-col flex-1 flex'>
          <div className='font-semibold p-0 line-clamp-1'>
            {data?.subject ?? '[No subject]'}
          </div>
          <div className='p-0 max-w-[calc(476px-77px)]'>
            {data?.description ? (
              <MarkdownContentRenderer
                showAsInlineText
                markdownContent={data?.description}
                className='text-gray-500 line-clamp-3 text-sm'
              />
            ) : (
              <p className='text-gray-500 line-clamp-3 text-sm'>
                {'[No description]'}
              </p>
            )}
          </div>
        </div>
        <div className='footer p-0 relative h-[108px] flex flex-col justify-center min-w-[71px] border-l border-gray-200 border-dashed'>
          <div className='flex flex-col items-center justify-center overflow-hidden h-[103px] min-w-[66px] relative rounded-md'>
            <Tag
              size='md'
              variant='outline'
              colorScheme={statusColorScheme}
              className='z-[1] cursor-pointer'
            >
              <TagLabel>
                {['solved', 'closed'].includes(data.issueStatus)
                  ? 'Closed'
                  : 'Open'}
              </TagLabel>
            </Tag>
            <IssueBgPattern className='absolute w-[120%] h-full' />
          </div>
        </div>
      </div>
    </div>
  );
};

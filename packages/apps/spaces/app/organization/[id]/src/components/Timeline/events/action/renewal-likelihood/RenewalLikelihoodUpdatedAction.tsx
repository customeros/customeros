// import React from 'react';
//
// import { Flex } from '@ui/layout/Flex';
// import { Text } from '@ui/typography/Text';
// import { Icons, FeaturedIcon } from '@ui/media/Icon';
// import { Action, RenewalLikelihoodProbability } from '@graphql/types';
// import { getARRColor } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
// import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';
//
// import { getLikelihoodDisplayData } from '../utils';
//
// interface RenewalForecastUpdatedActionProps {
//   data: Action;
// }
//
// export const RenewalLikelihoodUpdatedAction: React.FC<
//   RenewalForecastUpdatedActionProps
// > = ({ data }) => {
//   const { openModal } = useTimelineEventPreviewMethodsContext();
//   if (!data.content) return null;
//   const { preText, likelihood, author } = getLikelihoodDisplayData(
//     data.content,
//   );
//
//   return (
//     <Flex
//       alignItems='center'
//       onClick={() => openModal(data.id)}
//       cursor='pointer'
//     >
//       <FeaturedIcon
//         size='md'
//         minW='10'
//         colorScheme={getARRColor(
//           likelihood.toUpperCase() as RenewalLikelihoodProbability,
//         )}
//       >
//         <Icons.HeartActivity />
//       </FeaturedIcon>
//
//       <Text
//         my={1}
//         maxW='500px'
//         noOfLines={2}
//         ml={2}
//         fontSize='sm'
//         color='gray.700'
//       >
//         {preText}
//         <Text as='span' fontWeight='semibold'>
//           {likelihood}
//         </Text>
//         <Text color='gray.500' as='span' ml={1}>
//           by {author}
//         </Text>
//       </Text>
//     </Flex>
//   );
// };

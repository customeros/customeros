// import React from 'react';
//
// import { Flex } from '@ui/layout/Flex';
// import { Action } from '@graphql/types';
// import { Text } from '@ui/typography/Text';
// import { Icons, FeaturedIcon } from '@ui/media/Icon';
// import { getARRColor } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
// import {
//   getMetadata,
//   getCurrencyString,
// } from '@organization/src/components/Timeline/events/action/utils';
// import { useTimelineEventPreviewMethodsContext } from '@organization/src/components/Timeline/preview/context/TimelineEventPreviewContext';
// const DEFAULT_COLOR_SCHEME = 'gray';
//
// interface RenewalForecastUpdatedActionProps {
//   data: Action;
// }
//
// export const RenewalForecastUpdatedAction: React.FC<
//   RenewalForecastUpdatedActionProps
// > = ({ data }) => {
//   const { openModal } = useTimelineEventPreviewMethodsContext();
//   const forecastedAmount = data.content && getCurrencyString(data.content);
//   const [preText, postText] = data.content?.split('by ') ?? [];
//   const isCreatedBySystem = data.content?.includes('default');
//   const metadata = getMetadata(data?.metadata);
//   const colorScheme =
//     forecastedAmount && isCreatedBySystem
//       ? getARRColor(metadata?.likelihood)
//       : DEFAULT_COLOR_SCHEME;
//
//   const authorText = isCreatedBySystem ? data.content : `${preText} by`;
//
//   return (
//     <Flex
//       alignItems='center'
//       onClick={() => openModal(data.id)}
//       cursor='pointer'
//     >
//       <FeaturedIcon size='md' minW='10' colorScheme={colorScheme}>
//         <Icons.Calculator />
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
//         {authorText}
//         {!isCreatedBySystem && (
//           <Text color='gray.500' as='span' ml={1}>
//             {postText}
//           </Text>
//         )}
//       </Text>
//     </Flex>
//   );
// };

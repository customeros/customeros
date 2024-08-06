export const isUserPlatformMac = () => {
  // @ts-expect-error navigator is a global variable
  return navigator.userAgentData.platform.indexOf('macOS') >= 0;
};

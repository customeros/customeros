export function genRandomNum(maxLimit = 100) {
  let rand = Math.random() * maxLimit;

  rand = Math.floor(rand);

  return rand;
}

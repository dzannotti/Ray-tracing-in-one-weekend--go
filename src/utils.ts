export const randomNum = () => {
  return Math.random();
};

export const randomBetween = (min: number, max: number) => {
  return min + (max - min) * randomNum();
};

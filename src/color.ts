import { Vec3 } from "./vec3";

const n = (m: number) => Math.floor(m * 255);

const linearToGamma = (linearComponent: number) => {
  return linearComponent > 0 ? Math.sqrt(linearComponent) : 0;
};

function fillRGB(
  ctx: CanvasRenderingContext2D,
  r: number,
  g: number,
  b: number,
) {
  r = linearToGamma(r);
  g = linearToGamma(g);
  b = linearToGamma(b);

  ctx.fillStyle = `rgb(${n(r)},${n(g)},${n(b)})`;
}

export function writeColor(
  ctx: CanvasRenderingContext2D,
  color: Vec3, // each component is [0,1]
  j: number,
  i: number,
) {
  fillRGB(ctx, color.x, color.y, color.z);
  ctx.fillRect(i, j, 1, 1);
}

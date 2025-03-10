import { writeColor } from "./color";
import { ray, Ray } from "./ray";
import { vec3, color, point3, Vec3 } from "./vec3";

const canvas = document.querySelector("canvas");
if (!canvas) throw new Error("No canvas");
const ctx = canvas?.getContext("2d");
if (!ctx) throw new Error("No ctx");

function hitSphere(center: Vec3, radius: number, r: Ray) {
  const oc = center.sub(r.origin);
  const a = r.direction.lengthSquared;
  const h = Vec3.dot(r.direction, oc);
  const c = oc.lengthSquared - radius * radius;
  const discriminant = h * h - a * c;

  if (discriminant < 0) {
    return -1;
  }

  const t = (h - Math.sqrt(discriminant)) / a;

  return t;
}

const rayColor = (r: Ray) => {
  const t = hitSphere(point3(0, 0, -1), 0.5, r);
  if (t > 0) {
    const N = r.at(t).sub(vec3(0, 0, -1)).unit;
    const c = N.add(color(1, 1, 1)).div(2);
    return c;
  }

  const unitDirection = r.direction.unit;
  const a = 0.5 * (unitDirection.y + 1.0);

  return color(1.0, 1.0, 1.0)
    .k(1 - a)
    .add(color(0.5, 0.7, 1.0).k(a));
};

// Image
const aspectRatio = 16 / 9;
const width = 800;
const height = Math.floor(width / aspectRatio);
canvas.width = width;
canvas.height = height;

// Camera
const focalLength = 1.0;
const viewportHeight = 2.0;
const viewportWidth = viewportHeight * (width / height);
const cameraCenter = point3(0, 0, 0);

const viewportU = vec3(viewportWidth, 0, 0);
const viewportV = vec3(0, -viewportHeight, 0);

const pixelDeltaU = viewportU.div(width);
const pixelDeltaV = viewportV.div(height);

const viewportUpperLeft = cameraCenter
  .sub(vec3(0, 0, focalLength))
  .sub(viewportU.div(2))
  .sub(viewportV.div(2));
const pixel00Loc = viewportUpperLeft.add(pixelDeltaU.add(pixelDeltaV).div(2));

for (let j = 0; j <= height - 1; j++) {
  console.log(`Rendering scanline ${j}`);
  for (let i = 0; i <= width - 1; i++) {
    const pixelCenter = pixel00Loc.add(pixelDeltaU.k(i)).add(pixelDeltaV.k(j));
    const rayDirection = pixelCenter.sub(cameraCenter);

    const r = ray(cameraCenter, rayDirection);

    const color = rayColor(r);

    writeColor(ctx, color, j, i);
  }
}
console.log("Done!");

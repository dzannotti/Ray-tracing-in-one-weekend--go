import { writeColor } from "./color";
import { HitRecord, Hittable } from "./hittable";
import { HittableList } from "./hittable-list";
import { ray, Ray } from "./ray";
import { Sphere } from "./sphere";
import { vec3, color, point3 } from "./vec3";

const canvas = document.querySelector("canvas");
if (!canvas) throw new Error("No canvas");
const ctx = canvas?.getContext("2d");
if (!ctx) throw new Error("No ctx");

const rayColor = (r: Ray, world: Hittable) => {
  const rec = new HitRecord();

  const [hasHit, resultRec] = world.hit(r, 0, Infinity, rec);
  if (hasHit) {
    return resultRec.normal!.add(color(1, 1, 1)).div(2);
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

// World

const world = new HittableList();
world.add(new Sphere(point3(0, 0, -1), 0.5));
world.add(new Sphere(point3(0, -100.5, -1), 100));

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

    const color = rayColor(r, world);

    writeColor(ctx, color, j, i);
  }
}
console.log("Done!");

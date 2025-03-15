import { writeColor } from "./color";
import { HitRecord, Hittable } from "./hittable";
import { interval } from "./interval";
import { ray, Ray } from "./ray";
import { ref, Ref } from "./ref";
import { degreesToRadians, randomNum } from "./utils";
import { color, point3, Vec3, vec3 } from "./vec3";

const sampleSquare = () => vec3(randomNum() - 0.5, randomNum() - 0.5, 0);

export class Camera {
  height: number;
  pixel00Loc!: Vec3;
  pixelDeltaU!: Vec3;
  pixelDeltaV!: Vec3;
  cameraCenter!: Vec3;
  pixelSamplesScale: number;
  canvas: HTMLCanvasElement;
  ctx: CanvasRenderingContext2D;
  width: number;
  aspectRatio: number;
  samplesPerPixel: number;
  maxDepth: number;
  vfov: number;
  lookFrom: Vec3;
  lookAt: Vec3;
  vUp: Vec3;
  defocusAngle: number;
  focusDist: number;
  defocusDiskU!: Vec3;
  defocusDiskV!: Vec3;

  constructor({
    width = 800,
    aspectRatio = 16 / 9,
    samplesPerPixel = 10,
    maxDepth = 10,
    vfov = 90,
    lookFrom = point3(0, 0, 0),
    lookAt = point3(0, 0, -1),
    vUp = vec3(0, 1, 0),
    defocusAngle = 0,
    focusDist = 10,
  }: {
    width?: number;
    aspectRatio?: number;
    samplesPerPixel?: number;
    maxDepth?: number;
    vfov?: number;
    lookFrom?: Vec3;
    lookAt?: Vec3;
    vUp?: Vec3;
    defocusAngle?: number;
    focusDist?: number;
  }) {
    const canvas = document.querySelector("canvas");
    if (!canvas) throw new Error("No canvas");
    this.canvas = canvas;
    const ctx = this.canvas?.getContext("2d");
    if (!ctx) throw new Error("No ctx");
    this.ctx = ctx;

    this.width = width;
    this.aspectRatio = aspectRatio;
    this.samplesPerPixel = samplesPerPixel;
    this.maxDepth = maxDepth;

    // Image
    this.height = Math.floor(width / aspectRatio);
    this.canvas.width = width;
    this.canvas.height = this.height;

    this.pixelSamplesScale = 1 / samplesPerPixel;

    this.vfov = vfov;
    this.lookFrom = lookFrom;
    this.lookAt = lookAt;
    this.vUp = vUp;
    this.defocusAngle = defocusAngle;
    this.focusDist = focusDist;

    this.initialize();
  }

  initialize() {
    this.cameraCenter = this.lookFrom.clone();

    // Camera
    const theta = degreesToRadians(this.vfov);
    const h = Math.tan(theta / 2);
    const viewportHeight = 2.0 * h * this.focusDist;
    const viewportWidth = viewportHeight * (this.width / this.height);

    // Calculate camera basis vectors for camera coordinates
    const w = this.lookFrom.sub(this.lookAt).unit;
    const u = Vec3.cross(this.vUp, w).unit;
    const v = Vec3.cross(w, u);

    const viewportU = u.k(viewportWidth);
    const viewportV = v.k(-viewportHeight);

    this.pixelDeltaU = viewportU.div(this.width);
    this.pixelDeltaV = viewportV.div(this.height);

    const viewportUpperLeft = this.cameraCenter
      .sub(w.k(this.focusDist))
      .sub(viewportU.div(2))
      .sub(viewportV.div(2));
    this.pixel00Loc = viewportUpperLeft.add(
      this.pixelDeltaU.add(this.pixelDeltaV).div(2),
    );

    // camera defocus disk basis vectors
    const defocusRadius =
      this.focusDist * Math.tan(degreesToRadians(this.defocusAngle / 2));
    this.defocusDiskU = u.k(defocusRadius);
    this.defocusDiskV = v.k(defocusRadius);
  }

  render(world: Hittable) {
    for (let j = 0; j <= this.height - 1; j++) {
      console.log(`Rendering scanline ${j}`);
      for (let i = 0; i <= this.width - 1; i++) {
        let pixelColor = vec3(0, 0, 0);

        for (let sample = 0; sample < this.samplesPerPixel; sample++) {
          const r = this.getRay(i, j);
          pixelColor = pixelColor.add(this.rayColor(r, this.maxDepth, world));
        }

        writeColor(this.ctx, pixelColor.k(this.pixelSamplesScale), j, i);
      }
    }
    console.log("Done!");
  }

  rayColor(r: Ray, depth: number, world: Hittable): Vec3 {
    if (depth <= 0) return color(0, 0, 0);

    const rec = new HitRecord();

    const [hasHit, resultRec] = world.hit(r, interval(0.0001, Infinity), rec);
    if (hasHit) {
      let scattered: Ref<Ray> = {};
      let attenuation: Ref<Vec3> = {};

      if (
        resultRec.material!.scatter(ref(r), resultRec, attenuation, scattered)
      ) {
        return this.rayColor(scattered.value!, depth - 1, world).vectorMultiply(
          attenuation.value!,
        );
      } else {
        return color(0, 0, 0);
      }
    }

    const unitDirection = r.direction.unit;
    const a = 0.5 * (unitDirection.y + 1.0);

    return color(1.0, 1.0, 1.0)
      .k(1 - a)
      .add(color(0.5, 0.7, 1.0).k(a));
  }

  getRay(i: number, j: number) {
    const offset = sampleSquare();

    const pixelSample = this.pixel00Loc
      .add(this.pixelDeltaU.k(i + offset.x))
      .add(this.pixelDeltaV.k(j + offset.y));

    const rayOrigin =
      this.defocusAngle <= 0 ? this.cameraCenter : this.defocusDiskSample();
    const rayDirection = pixelSample.sub(rayOrigin);

    return ray(rayOrigin, rayDirection);
  }

  defocusDiskSample() {
    const p = Vec3.randomInUnitDisk();
    return this.cameraCenter
      .add(this.defocusDiskU.k(p.x))
      .add(this.defocusDiskV.k(p.y));
  }
}

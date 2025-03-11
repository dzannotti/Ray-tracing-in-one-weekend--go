import { writeColor } from "./color";
import { HitRecord, Hittable } from "./hittable";
import { interval } from "./interval";
import { ray, Ray } from "./ray";
import { color, point3, Vec3, vec3 } from "./vec3";

export class Camera {
  height: number;
  pixel00Loc!: Vec3;
  pixelDeltaU!: Vec3;
  pixelDeltaV!: Vec3;
  cameraCenter!: Vec3;

  constructor(
    public width: number = 800,
    public aspectRatio: number = 16 / 9,
    public canvas: HTMLCanvasElement,
    public ctx: CanvasRenderingContext2D,
  ) {
    // Image
    this.height = Math.floor(width / aspectRatio);
    canvas.width = width;
    canvas.height = this.height;
    this.initialize();
  }

  render(world: Hittable) {
    for (let j = 0; j <= this.height - 1; j++) {
      console.log(`Rendering scanline ${j}`);
      for (let i = 0; i <= this.width - 1; i++) {
        const pixelCenter = this.pixel00Loc
          .add(this.pixelDeltaU.k(i))
          .add(this.pixelDeltaV.k(j));
        const rayDirection = pixelCenter.sub(this.cameraCenter);

        const r = ray(this.cameraCenter, rayDirection);

        const color = this.rayColor(r, world);

        writeColor(this.ctx, color, j, i);
      }
    }
    console.log("Done!");
  }

  initialize() {
    // Camera
    const focalLength = 1.0;
    const viewportHeight = 2.0;
    const viewportWidth = viewportHeight * (this.width / this.height);
    this.cameraCenter = point3(0, 0, 0);

    const viewportU = vec3(viewportWidth, 0, 0);
    const viewportV = vec3(0, -viewportHeight, 0);

    this.pixelDeltaU = viewportU.div(this.width);
    this.pixelDeltaV = viewportV.div(this.height);

    const viewportUpperLeft = this.cameraCenter
      .sub(vec3(0, 0, focalLength))
      .sub(viewportU.div(2))
      .sub(viewportV.div(2));
    this.pixel00Loc = viewportUpperLeft.add(
      this.pixelDeltaU.add(this.pixelDeltaV).div(2),
    );
  }

  rayColor(r: Ray, world: Hittable) {
    const rec = new HitRecord();

    const [hasHit, resultRec] = world.hit(r, interval(0, Infinity), rec);
    if (hasHit) {
      return resultRec.normal!.add(color(1, 1, 1)).div(2);
    }

    const unitDirection = r.direction.unit;
    const a = 0.5 * (unitDirection.y + 1.0);

    return color(1.0, 1.0, 1.0)
      .k(1 - a)
      .add(color(0.5, 0.7, 1.0).k(a));
  }
}

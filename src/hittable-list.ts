import { HitRecord, Hittable } from "./hittable";
import { interval, Interval } from "./interval";
import { Ray } from "./ray";

export class HittableList extends Hittable {
  objects: Hittable[] = [];

  clear() {
    this.objects = [];
  }

  add(obj: Hittable) {
    this.objects.push(obj);
  }

  hit(r: Ray, rayT: Interval, rec: HitRecord) {
    let tempRec = new HitRecord();
    let hitAnything = false;
    let closestSoFar = rayT.max;

    for (const obj of this.objects) {
      const [hasHit, resultRec] = obj.hit(
        r,
        interval(rayT.min, closestSoFar),
        tempRec,
      );
      if (hasHit) {
        hitAnything = true;
        closestSoFar = resultRec.t!;
        rec = resultRec;
      }
    }

    return [hitAnything, rec] as [boolean, HitRecord];
  }
}

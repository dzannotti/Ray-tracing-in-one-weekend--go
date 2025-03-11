import { HitRecord, Hittable } from "./hittable";
import { Ray } from "./ray";

export class HittableList {
  objects: Hittable[] = [];
  
  clear() {
    this.objects = [];
  }

  add(obj: Hittable) {
    this.objects.push(obj);
  }

  hit(r: Ray, rayTMin: number, rayTMax: number, rec: HitRecord) {
    let tempRec = new HitRecord();
    let hitAnything = false;
    let closestSoFar = rayTMax;

    for (const obj of this.objects) {
      const [hasHit, resultRec] = obj.hit(r, rayTMin, closestSoFar, tempRec);
      if (hasHit) {
        hitAnything = true;
        closestSoFar = resultRec.t!;
        rec = resultRec;
      }
    }

    return [hitAnything, rec];
  }
}

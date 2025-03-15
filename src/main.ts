import { Camera } from "./camera";
import { HittableList } from "./hittable-list";
import { Dialectric, Lambertian, Metal } from "./material";
import { Sphere } from "./sphere";
import { color, point3 } from "./vec3";

const ground = new Lambertian(color(0.8, 0.8, 0));
const center = new Lambertian(color(0.1, 0.2, 0.5));
const left = new Dialectric(1.5);
const bubble = new Dialectric(1 / 1.5);
const right = new Metal(color(0.8, 0.6, 0.2), 1.0);

// World
const world = new HittableList();
world.add(new Sphere(point3(0, -100.5, -1), 100, ground));
world.add(new Sphere(point3(0, 0, -1.2), 0.5, center));
world.add(new Sphere(point3(-1, 0, -1), 0.5, left));
world.add(new Sphere(point3(-1, 0, -1), 0.4, bubble));
world.add(new Sphere(point3(1, 0, -1), 0.5, right));

const cam = new Camera({
  width: 600,
  aspectRatio: 16 / 9,
  samplesPerPixel: 100,
  maxDepth: 50,
});

cam.render(world);
